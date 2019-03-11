package main

import (
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
		"math/rand"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"strconv"
	"strings"
	"math"
)

const (
	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"

	TIMEZONE="Asia/Saigon"

	INDEX_HOUSE_LIST_SIZE=5
	SEARCH_HOUSE_LIST_SIZE=2

	FDFS_HOST="45.63.94.102"
	FDFS_PORT="8888"

	USER_OBJECT_TYPE="User"
	HOUSE_OBJECT_TYPE="House"
	HOUSE_IMAGE_OBJECT_TYPE="HouseImage"
	AREA_OBJECT_TYPE="Area"
	FACILITY_OBJECT_TYPE="Facility"
	ORDER_OBJECT_TYPE="Order"
)

func checkArgsNum(args []string,n int) error {
	if len(args)!=n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}

	return nil
}

func displayTime(x time.Time) string {
	loc,_:=time.LoadLocation(TIMEZONE)
	return x.In(loc).Format("2006-01-02 15:04:05")
}

func displayDate(x time.Time) string {
	loc,_:=time.LoadLocation(TIMEZONE)
	return x.In(loc).Format("2006-01-02")
}

func addDomain2Url(url string) string {
	if url=="" {
		return ""
	}
	return "http://"+FDFS_HOST+":"+FDFS_PORT+"/"+url
}

func getInfo(stub shim.ChaincodeStubInterface,key string) ([]byte,error) {
	data,err:=stub.GetState(key)
	if err!=nil {
		return nil,err
	}
	if data==nil {
		return nil,errors.New("key "+key+" doesn't exist")
	}

	return data,nil
}

func getQueryResult(stub shim.ChaincodeStubInterface,query map[string]interface{}) (shim.StateQueryIteratorInterface, error) {
	_query,err:=json.Marshal(query)
	if err!=nil {
		return nil,err
	}

	queryStr:=string(_query)
	return stub.GetQueryResult(queryStr)
}

func keyExist(stub shim.ChaincodeStubInterface,key string) (bool,error) {
	data,err:=stub.GetState(key)
	if err!=nil {
		return false,err
	}
	if data==nil {
		return false,nil
	}

	return true,nil
}

func generateId(stub shim.ChaincodeStubInterface,prefix string,length int) (string,error) {
	rand.Seed(time.Now().UnixNano())
	x:=[]string{"1","2","3","4","5","6","7","8","9"}
	nx:=len(x)
	idStr:=""
	for i:=0; i<length; i++ {
		idStr+=x[rand.Intn(nx)]
	}

	key:=prefix+idStr
	exist,err:=keyExist(stub,key)
	if err!=nil {
		return "",err
	}
	if !exist {
		return idStr,nil
	}

	return generateId(stub,prefix,length)
}

type User struct {
	ObjectType string
	Id string
	Name string
	PublicKey string
	Mobile string
	RealName string
	IdCard string
	AvatarUrl string
}

func getUserById(stub shim.ChaincodeStubInterface,userId string) (*User,error) {
	uKey:=USER_OBJECT_TYPE+userId
	data,err:=getInfo(stub,uKey)
	if err!=nil {
		return nil,err
	}

	var user User
	err=json.Unmarshal(data,&user)
	if err!=nil {
		return nil,err
	}

	return &user,nil
}

func getUserByMobile(stub shim.ChaincodeStubInterface,mobile string) (*User,error) {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":USER_OBJECT_TYPE,
			"Mobile":mobile,
		},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return nil,err
	}
	defer iter.Close()

	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return nil,err
		}

		var user User
		err=json.Unmarshal(item.Value,&user)
		if err!=nil {
			return nil,err
		}

		return &user,nil
	}

	return nil,errors.New("user doesn't exist")
}

func userMobileExist(stub shim.ChaincodeStubInterface,mobile string) (bool,error) {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":USER_OBJECT_TYPE,
			"Mobile":mobile,
		},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return false,err
	}
	defer iter.Close()

	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return false,err
		}

		if item.Value!=nil {
			return true,nil
		}
	}

	return false,nil
}

func setUserInfo(stub shim.ChaincodeStubInterface,userMobile string,setUser func(user *User)) error {
	user,err:=getUserByMobile(stub,userMobile)
	if err!=nil {
		return err
	}

	setUser(user)

	data,err:=json.Marshal(user)
	if err!=nil {
		return err
	}

	uKey:=USER_OBJECT_TYPE+user.Id
	err=stub.PutState(uKey,data)
	if err!=nil {
		return err
	}

	return nil
}

type House struct {
	ObjectType string
	Id string
	Title string
	Price int
	Address string
	RoomCount int
	Acreage int
	Unit string
	Capacity int
	Beds string
	Deposit int
	MinDays int
	MaxDays int
	OrderCount int
	IndexImageUrl string
	CreateTime time.Time

	UserId string
	AreaId string

	FacilityIds []string
}

type HouseImage struct {
	ObjectType string
	Id string
	Url string

	HouseId string
}

func getHouse(stub shim.ChaincodeStubInterface,houseId string) (*House,error) {
	hKey:=HOUSE_OBJECT_TYPE+houseId
	_data,err:=getInfo(stub,hKey)
	if err!=nil {
		return nil,err
	}

	var house House
	err=json.Unmarshal(_data,&house)
	if err!=nil {
		return nil,err
	}

	return &house,nil
}

func getHouseIdListByLandlord(stub shim.ChaincodeStubInterface,userMobile string) ([]string,error) {
	user,err:=getUserByMobile(stub,userMobile)
	if err!=nil {
		return nil,err
	}
	userId:=user.Id

	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":HOUSE_OBJECT_TYPE,
			"UserId":userId,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"CreateTime":"desc"},
		},
		"use_index":[]string{"_design/houseDoc","house"},
		"fields":[]string{"Id"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return nil,err
	}
	defer iter.Close()

	var houseIdList []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return nil,err
		}

		var house House
		err=json.Unmarshal(item.Value,&house)
		if err!=nil {
			return nil,err
		}

		houseIdList=append(houseIdList,house.Id)
	}

	return houseIdList,nil
}

func getHouseImageUrlList(stub shim.ChaincodeStubInterface,houseId string) ([]string,error) {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":HOUSE_IMAGE_OBJECT_TYPE,
			"HouseId":houseId,
		},
		"use_index":[]string{"_design/houseImageDoc","houseImage"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return nil,err
	}
	defer iter.Close()

	var houseImageUrlList []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return nil,err
		}

		var houseImage HouseImage
		err=json.Unmarshal(item.Value,&houseImage)
		if err!=nil {
			return nil,err
		}

		houseImageUrlList=append(houseImageUrlList,addDomain2Url(houseImage.Url))
	}

	return houseImageUrlList,nil
}

func getHouseInfo(stub shim.ChaincodeStubInterface,houseId string) (map[string]interface{},error) {
	house,err:=getHouse(stub,houseId)
	if err!=nil {
		return nil,err
	}

	user,err:=getUserById(stub,house.UserId)
	if err!=nil {
		return nil,err
	}

	area,err:=getArea(stub,house.AreaId)
	if err!=nil {
		return nil,err
	}

	data:=make(map[string]interface{})
	data["house_id"]=houseId
	data["title"]=house.Title
	data["price"]=house.Price
	data["area_name"]=area.Name
	data["index_image_url"]=addDomain2Url(house.IndexImageUrl)
	data["room_count"]=house.RoomCount
	data["order_count"]=house.OrderCount
	data["address"]=house.Address
	data["user_avatar"]=addDomain2Url(user.AvatarUrl)
	data["create_time"]=displayTime(house.CreateTime)
	return data,nil
}

func getHouseDesc(stub shim.ChaincodeStubInterface,houseId string) (map[string]interface{},error) {
	house,err:=getHouse(stub,houseId)
	if err!=nil {
		return nil,err
	}

	user,err:=getUserById(stub,house.UserId)
	if err!=nil {
		return nil,err
	}

	var facilities []string
	for _,facilityId:=range house.FacilityIds{
		facility,err:=getFacility(stub,facilityId)
		if err!=nil {
			return nil,err
		}

		facilities=append(facilities,facility.Name)
	}

	orderIdList,err:=getOrderIdListByHouse(stub,houseId)
	if err!=nil {
		return nil,err
	}

	var comments []map[string]interface{}
	for _,orderId:=range orderIdList{
		order,err:=getOrder(stub,orderId)
		if err!=nil {
			return nil,err
		}

		if order.Comment=="" {
			continue
		}

		_user,err:=getUserById(stub,order.UserId)
		if err!=nil {
			return nil,err
		}

		comment:=make(map[string]interface{})
		comment["username"]=_user.Name
		comment["comment"]=order.Comment
		comment["create_time"]=displayTime(order.CreateTime)
		comments=append(comments,comment)
	}

	imgUrls,err:=getHouseImageUrlList(stub,houseId)
	if err!=nil {
		return nil,err
	}

	data:=make(map[string]interface{})
	data["house_id"]=houseId
	data["user_id"]=user.Id
	data["username"]=user.Name
	data["user_avatar"]=addDomain2Url(user.AvatarUrl)
	data["title"]=house.Title
	data["price"]=house.Price
	data["address"]=house.Address
	data["room_count"]=house.RoomCount
	data["acreage"]=house.Acreage
	data["unit"]=house.Unit
	data["capacity"]=house.Capacity
	data["beds"]=house.Beds
	data["deposit"]=house.Deposit
	data["min_days"]=house.MinDays
	data["max_days"]=house.MaxDays
	data["img_urls"]=imgUrls
	data["facilities"]=facilities
	data["comments"]=comments
	return data,nil
}

type Area struct {
	ObjectType string
	Id string
	Name string
}

func getArea(stub shim.ChaincodeStubInterface,areaId string) (*Area,error) {
	aKey:=AREA_OBJECT_TYPE+areaId
	_data,err:=getInfo(stub,aKey)
	if err!=nil {
		return nil,err
	}

	var area Area
	err=json.Unmarshal(_data,&area)
	if err!=nil {
		return nil,err
	}

	return &area,nil
}

type Facility struct {
	ObjectType string
	Id string
	Name string
}

func getFacility(stub shim.ChaincodeStubInterface,facilityId string) (*Facility,error) {
	fKey:=FACILITY_OBJECT_TYPE+facilityId
	_data,err:=getInfo(stub,fKey)
	if err!=nil {
		return nil,err
	}

	var facility Facility
	err=json.Unmarshal(_data,&facility)
	if err!=nil {
		return nil,err
	}

	return &facility,nil
}

type Order struct {
	ObjectType string
	Id string
	BeginDate time.Time
	EndDate time.Time
	Days int
	HousePrice int
	Amount int
	Status string
	Comment string
	CreateTime time.Time
	Credit bool

	UserId string
	HouseId string
}

func getOrder(stub shim.ChaincodeStubInterface,orderId string) (*Order,error) {
	oKey:=ORDER_OBJECT_TYPE+orderId
	_data,err:=getInfo(stub,oKey)
	if err!=nil {
		return nil,err
	}

	var order Order
	err=json.Unmarshal(_data,&order)
	if err!=nil {
		return nil,err
	}

	return &order,nil
}

func getOrderInfo(stub shim.ChaincodeStubInterface,orderId string) (map[string]interface{},error) {
	order,err:=getOrder(stub,orderId)
	if err!=nil {
		return nil,err
	}

	house,err:=getHouse(stub,order.HouseId)
	if err!=nil {
		return nil,err
	}

	data:=make(map[string]interface{})
	data["order_id"]=orderId
	data["title"]=house.Title
	data["image_url"]=addDomain2Url(house.IndexImageUrl)
	data["begin_date"]=displayDate(order.BeginDate)
	data["end_date"]=displayDate(order.EndDate)
	data["create_time"]=displayTime(order.CreateTime)
	data["days"]=order.Days
	data["amount"]=order.Amount
	data["status"]=order.Status
	data["comment"]=order.Comment
	data["credit"]=order.Credit
	return data,nil
}

func getOrderIdListByRenter(stub shim.ChaincodeStubInterface,renterId string) ([]string,error) {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":ORDER_OBJECT_TYPE,
			"UserId":renterId,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"CreateTime":"desc"},
		},
		"use_index":[]string{"_design/orderRenterDoc","orderRenter"},
		"fields":[]string{"Id"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return nil,err
	}
	defer iter.Close()

	var orderIdList []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return nil,err
		}

		var order Order
		err=json.Unmarshal(item.Value,&order)
		if err!=nil {
			return nil,err
		}

		orderIdList=append(orderIdList,order.Id)
	}

	return orderIdList,nil
}

func getOrderIdListByHouse(stub shim.ChaincodeStubInterface,houseId string) ([]string,error) {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":ORDER_OBJECT_TYPE,
			"HouseId":houseId,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"CreateTime":"desc"},
		},
		"use_index":[]string{"_design/orderDoc","order"},
		"fields":[]string{"Id"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return nil,err
	}
	defer iter.Close()

	var orderIdList []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return nil,err
		}

		var order Order
		err=json.Unmarshal(item.Value,&order)
		if err!=nil {
			return nil,err
		}

		orderIdList=append(orderIdList,order.Id)
	}

	return orderIdList,nil
}

type RentingChaincode struct {

}

func (this *RentingChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *RentingChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	switch fn {
	case "getAreaList":
		return this.getAreaList(stub,args)
	case "generateUserId":
		return this.generateUserId(stub,args)
	case "register":
		return this.register(stub,args)
	case "getUserId":
		return this.getUserId(stub,args)
	case "getUserPublicKey":
		return this.getUserPublicKey(stub,args)
	case "getUserInfo":
		return this.getUserInfo(stub,args)
	case "avatar":
		return this.avatar(stub,args)
	case "rename":
		return this.rename(stub,args)
	case "auth":
		return this.auth(stub,args)
	case "addHouse":
		return this.addHouse(stub,args)
	case "uploadHouseImage":
		return this.uploadHouseImage(stub,args)
	case "getLandlordHouseList":
		return this.getLandlordHouseList(stub,args)
	case "getHouseDesc":
		return this.getHouseDesc(stub,args)
	case "getIndexHouseList":
		return this.getIndexHouseList(stub,args)
	case "searchHouse":
		return this.searchHouse(stub,args)
	case "addOrder":
		return this.addOrder(stub,args)
	case "getOrderList":
		return this.getOrderList(stub,args)
	case "handleOrder":
		return this.handleOrder(stub,args)
	case "comment":
		return this.comment(stub,args)
	case "getOrderHouseId":
		return this.getOrderHouseId(stub,args)
	/*
	case "addArea":
		return this.addArea(stub,args)
	case "getFacilityList":
		return this.getFacilityList(stub,args)
	case "addFacility":
		return this.addFacility(stub,args)
	 */
	default:
		return shim.Error("invalid method")
	}
}

func (this *RentingChaincode) getAreaList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":AREA_OBJECT_TYPE,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"Id":"asc"},
		},
		"use_index":[]string{"_design/areaDoc","area"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var areaList []map[string]interface{}
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}

		var _area Area
		err=json.Unmarshal(item.Value,&_area)
		if err!=nil {
			return shim.Error(err.Error())
		}

		area:=make(map[string]interface{})
		area["area_id"]=_area.Id
		area["area_name"]=_area.Name

		areaList=append(areaList,area)
	}

	data,err:=json.Marshal(areaList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) generateUserId(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	id,err:=generateId(stub,USER_OBJECT_TYPE,11)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(id))
}

func (this *RentingChaincode) register(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var user User
	user.ObjectType=USER_OBJECT_TYPE

	user.Id=args[0]
	uKey:=USER_OBJECT_TYPE+user.Id
	exist,err:=keyExist(stub,uKey)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if exist {
		return shim.Error("user exists")
	}

	user.Mobile=args[1]
	exist,err=userMobileExist(stub,user.Mobile)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if exist {
		return shim.Error("mobile has been registered")
	}

	user.Name=user.Mobile

	user.PublicKey=args[2]

	data,err:=json.Marshal(user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(uKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getUserInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	userInfo:=make(map[string]interface{})
	userInfo["user_id"]=user.Id
	userInfo["name"]=user.Name
	userInfo["mobile"]=user.Mobile
	userInfo["real_name"]=user.RealName
	userInfo["id_card"]=user.IdCard
	userInfo["avatar_url"]=addDomain2Url(user.AvatarUrl)

	data,err:=json.Marshal(userInfo)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(data))
}

func (this *RentingChaincode) getUserId(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(user.Id))
}

func (this *RentingChaincode) getUserPublicKey(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(user.PublicKey))
}

func (this *RentingChaincode) avatar(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub, args[0], func(user *User) {
		user.AvatarUrl=args[1]
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) rename(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub, args[0], func(user *User) {
		user.Name=args[1]
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) auth(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub, args[0], func(user *User) {
		user.RealName=args[1]
		user.IdCard=args[2]
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) addHouse(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,14)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var house House
	house.ObjectType=HOUSE_OBJECT_TYPE

	house.Id,err=generateId(stub,HOUSE_OBJECT_TYPE,14)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.UserId=user.Id

	house.Title=args[1]

	house.Price,err=strconv.Atoi(args[2])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Price*=100

	house.AreaId=args[3]

	house.Address=args[4]

	house.RoomCount,err=strconv.Atoi(args[5])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.Acreage,err=strconv.Atoi(args[6])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.Unit=args[7]

	house.Capacity,err=strconv.Atoi(args[8])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.Beds=args[9]

	house.Deposit,err=strconv.Atoi(args[10])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Deposit*=100

	house.MinDays,err=strconv.Atoi(args[11])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.MaxDays,err=strconv.Atoi(args[12])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.FacilityIds=strings.Split(args[13],",")

	loc,_:=time.LoadLocation(TIMEZONE)
	house.CreateTime=time.Now().In(loc)

	data,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	hKey:=HOUSE_OBJECT_TYPE+house.Id
	err=stub.PutState(hKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) uploadHouseImage(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	userId:=user.Id

	house,err:=getHouse(stub,args[1])
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.UserId!=userId {
		return shim.Error("invalid user")
	}

	var houseImage HouseImage
	houseImage.ObjectType=HOUSE_IMAGE_OBJECT_TYPE

	houseImage.Id,err=generateId(stub,HOUSE_IMAGE_OBJECT_TYPE,16)
	if err!=nil {
		return shim.Error(err.Error())
	}

	houseImage.Url=args[2]

	houseImage.HouseId=house.Id

	houseImageData,err:=json.Marshal(houseImage)
	if err!=nil {
		return shim.Error(err.Error())
	}

	hiKey:=HOUSE_IMAGE_OBJECT_TYPE+houseImage.Id
	err=stub.PutState(hiKey,houseImageData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.IndexImageUrl=="" {
		house.IndexImageUrl=houseImage.Url

		data,err:=json.Marshal(house)
		if err!=nil {
			return shim.Error(err.Error())
		}

		hKey:=HOUSE_OBJECT_TYPE+house.Id
		err=stub.PutState(hKey,data)
		if err!=nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getLandlordHouseList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	houseIdList,err:=getHouseIdListByLandlord(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	var houseList []map[string]interface{}
	for _,houseId:=range houseIdList{
		house,err:=getHouseInfo(stub,houseId)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseList=append(houseList,house)
	}

	data,err:=json.Marshal(houseList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) getHouseDesc(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	house,err:=getHouseDesc(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	data,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) getIndexHouseList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":HOUSE_OBJECT_TYPE,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"CreateTime":"desc"},
		},
		"use_index":[]string{"_design/houseIndexDoc","houseIndex"},
		"fields":[]string{"Id"},
		"limit":INDEX_HOUSE_LIST_SIZE,
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var houseList []map[string]interface{}
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}

		var _house House
		err=json.Unmarshal(item.Value,&_house)
		if err!=nil {
			return shim.Error(err.Error())
		}

		house,err:=getHouseInfo(stub,_house.Id)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseList=append(houseList,house)
	}

	data,err:=json.Marshal(houseList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) searchHouse(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,4)
	if err!=nil {
		return shim.Error(err.Error())
	}

	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":HOUSE_OBJECT_TYPE,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"CreateTime":"desc"},
		},
		"use_index":[]string{"_design/houseIndexDoc","houseIndex"},
		"fields":[]string{"Id"},
	}

	if args[0]!="" {
		querySelector:=query["selector"].(map[string]interface{})
		querySelector["AreaId"]=args[0]
		query["selector"]=querySelector
		query["use_index"]=[]string{"_design/houseAreaDoc","houseArea"}
	}
	if args[1]!="" {
		start,_:=time.Parse("2006-01-02 15:04:05",args[1]+" 00:00:00")
		querySelector:=query["selector"].(map[string]interface{})
		querySelector["CreateTime"]=map[string]interface{}{"$gte":start}
		query["selector"]=querySelector
	}
	if args[2]!="" {
		end,_:=time.Parse("2006-01-02 15:04:05",args[2]+" 23:59:59")
		querySelector:=query["selector"].(map[string]interface{})
		querySelector["CreateTime"]=map[string]interface{}{"$lte":end}
		query["selector"]=querySelector
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	page:=1
	if args[3]!="" {
		page,err=strconv.Atoi(args[3])
		if err!=nil {
			page=1
		}
	}
	skip:=SEARCH_HOUSE_LIST_SIZE*(page-1)

	var houseList []map[string]interface{}
	totalRows:=0
	for iter.HasNext() {
		totalRows++

		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}

		if totalRows<=skip || totalRows>skip+SEARCH_HOUSE_LIST_SIZE {
			continue
		}

		var _house House
		err=json.Unmarshal(item.Value,&_house)
		if err!=nil {
			return shim.Error(err.Error())
		}

		house,err:=getHouseInfo(stub,_house.Id)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseList=append(houseList,house)
	}

	_data:=make(map[string]interface{})
	_data["houses"]=houseList
	_data["total_house"]=totalRows
	_data["total_page"]=int(math.Ceil(float64(totalRows)/float64(SEARCH_HOUSE_LIST_SIZE)))
	_data["current_page"]=page

	data,err:=json.Marshal(_data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) addOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,4)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var order Order
	order.ObjectType=ORDER_OBJECT_TYPE

	order.Id,err=generateId(stub,ORDER_OBJECT_TYPE,16)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	order.UserId=user.Id

	house,err:=getHouse(stub,args[1])
	if err!=nil {
		return shim.Error(err.Error())
	}
	order.HouseId=house.Id

	if house.UserId==user.Id {
		return shim.Error("invalid user")
	}

	order.BeginDate,err=time.Parse("2006-01-02",args[2])
	if err!=nil {
		return shim.Error(err.Error())
	}

	order.EndDate,err=time.Parse("2006-01-02",args[3])
	if err!=nil {
		return shim.Error(err.Error())
	}

	order.Days=int(order.EndDate.Sub(order.BeginDate).Seconds()/60/60/24)+1

	order.HousePrice=house.Price

	order.Amount=order.HousePrice*order.Days

	loc,_:=time.LoadLocation(TIMEZONE)
	order.CreateTime=time.Now().In(loc)

	order.Status=ORDER_STATUS_WAIT_ACCEPT

	data,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	oKey:=ORDER_OBJECT_TYPE+order.Id
	err=stub.PutState(oKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getOrderList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var orderList []map[string]interface{}
	if args[1]=="landlord" {
		houseIdList,err:=getHouseIdListByLandlord(stub,args[0])
		if err!=nil {
			return shim.Error(err.Error())
		}

		for _,houseId:=range houseIdList{
			orderIdList,err:=getOrderIdListByHouse(stub,houseId)
			if err!=nil {
				return shim.Error(err.Error())
			}

			if orderIdList!=nil && len(orderIdList)>0 {
				for _,orderId:=range orderIdList{
					order,err:=getOrderInfo(stub,orderId)
					if err!=nil {
						return shim.Error(err.Error())
					}
					orderList=append(orderList,order)
				}
			}
		}
	} else {
		user,err:=getUserByMobile(stub,args[0])
		if err!=nil {
			return shim.Error(err.Error())
		}
		userId:=user.Id

		orderIdList,err:=getOrderIdListByRenter(stub,userId)
		if err!=nil {
			return shim.Error(err.Error())
		}

		if orderIdList!=nil && len(orderIdList)>0 {
			for _,orderId:=range orderIdList{
				order,err:=getOrderInfo(stub,orderId)
				if err!=nil {
					return shim.Error(err.Error())
				}
				orderList=append(orderList,order)
			}
		}
	}

	data,err:=json.Marshal(orderList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) handleOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	userId:=user.Id

	order,err:=getOrder(stub,args[1])
	if err!=nil {
		return shim.Error(err.Error())
	}

	house,err:=getHouse(stub,order.HouseId)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.UserId!=userId {
		return shim.Error("invalid user")
	}
	if order.Status!=ORDER_STATUS_WAIT_ACCEPT {
		return shim.Error("invalid order status")
	}

	if args[2]=="reject" {
		order.Status=ORDER_STATUS_REJECTED
	} else {
		order.Status=ORDER_STATUS_WAIT_COMMENT
	}

	data,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	oKey:=ORDER_OBJECT_TYPE+order.Id
	err=stub.PutState(oKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if order.Status==ORDER_STATUS_WAIT_COMMENT {
		house,err:=getHouse(stub,order.HouseId)
		if err!=nil {
			return shim.Error(err.Error())
		}

		house.OrderCount++

		houseData,err:=json.Marshal(house)
		if err!=nil {
			return shim.Error(err.Error())
		}

		hKey:=HOUSE_OBJECT_TYPE+order.HouseId
		err=stub.PutState(hKey,houseData)
		if err!=nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) comment(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user,err:=getUserByMobile(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	userId:=user.Id

	order,err:=getOrder(stub,args[1])
	if err!=nil {
		return shim.Error(err.Error())
	}

	if order.UserId!=userId {
		return shim.Error("invalid user")
	}
	if order.Status!=ORDER_STATUS_WAIT_COMMENT {
		return shim.Error("invalid order status")
	}

	order.Comment=args[2]
	order.Status=ORDER_STATUS_COMPLETE

	data,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	oKey:=ORDER_OBJECT_TYPE+order.Id
	err=stub.PutState(oKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getOrderHouseId(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	order,err:=getOrder(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(order.HouseId))
}

/*
func (this *RentingChaincode) addArea(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var area Area
	area.ObjectType=AREA_OBJECT_TYPE
	area.Id=args[0]
	area.Name=args[1]

	data,err:=json.Marshal(area)
	if err!=nil {
		return shim.Error(err.Error())
	}

	aKey:=AREA_OBJECT_TYPE+area.Id
	err=stub.PutState(aKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getFacilityList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	query:=map[string]interface{}{
		"selector":map[string]interface{}{
			"ObjectType":FACILITY_OBJECT_TYPE,
		},
		"sort":[]map[string]interface{}{
			map[string]interface{}{"Id":"asc"},
		},
		"use_index":[]string{"_design/facilityDoc","facility"},
	}

	iter,err:=getQueryResult(stub,query)
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var facilityList []map[string]interface{}
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}

		var _facility Facility
		err=json.Unmarshal(item.Value,&_facility)
		if err!=nil {
			return shim.Error(err.Error())
		}

		facility:=make(map[string]interface{})
		facility["facility_id"]=_facility.Id
		facility["facility_name"]=_facility.Name

		facilityList=append(facilityList,facility)
	}

	data,err:=json.Marshal(facilityList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) addFacility(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var facility Facility
	facility.ObjectType=FACILITY_OBJECT_TYPE
	facility.Id=args[0]
	facility.Name=args[1]

	data,err:=json.Marshal(facility)
	if err!=nil {
		return shim.Error(err.Error())
	}

	fKey:=FACILITY_OBJECT_TYPE+facility.Id
	err=stub.PutState(fKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
 */

func main()  {
	shim.Start(new(RentingChaincode))
}
