package main

import (
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
	"crypto/sha256"
	"strconv"
	"golang.org/x/crypto/ripemd160"
	"strings"
	"encoding/hex"
)

const (
	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"

	TIME_ZONE="Asia/Saigon"

	INDEX_HOUSE_LIST_SIZE=5

	FDFSHost="45.77.250.9"
	FDFSPort="8888"
)

type User struct {
	Id string
	Name string
	PublicKey string
	Mobile string
	RealName string
	IdCard string
	AvatarUrl string

	HouseIds []string
	OrderIds []string
}

type House struct {
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

	AreaId string
	UserId string

	FacilityIds []string
	HouseImageIds []string
	OrderIds []string
}

type Area struct {
	Id string
	Name string
}

type Facility struct {
	Id string
	Name string
}

type HouseImage struct {
	Id string
	Url string

	HouseId string
}

type Order struct {
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

type RentingChaincode struct {

}

func addDomain2Url(url string) string {
	if url=="" {
		return ""
	}
	return "http://"+FDFSHost+":"+FDFSPort+"/"+url
}

func displayTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func displayDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func getInfo(stub shim.ChaincodeStubInterface,key string) ([]byte,error) {
	data,err:=stub.GetState(key)
	if err!=nil {
		return nil,err
	}
	if data==nil {
		return nil,errors.New(key+" doesn't exist")
	}
	return data,nil
}

func getUser(stub shim.ChaincodeStubInterface,key string) (*User,error) {
	data,err:=getInfo(stub,key)
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

func getHouseInfo(stub shim.ChaincodeStubInterface,houseId string) (map[string]interface{},error) {
	hKey:="house_"+houseId
	houseData,err:=getInfo(stub,hKey)
	if err!=nil {
		return nil,err
	}

	var house House
	err=json.Unmarshal(houseData,&house)
	if err!=nil {
		return nil,err
	}

	aKey:="area_"+house.AreaId
	areaData,err:=getInfo(stub,aKey)
	if err!=nil {
		return nil,err
	}

	var area Area
	err=json.Unmarshal(areaData,&area)
	if err!=nil {
		return nil,err
	}

	uKey:="user_"+house.UserId
	userData,err:=getInfo(stub,uKey)
	if err!=nil {
		return nil,err
	}

	var user User
	err=json.Unmarshal(userData,&user)
	if err!=nil {
		return nil,err
	}

	info:=make(map[string]interface{})
	info["house_id"]=house.Id
	info["title"]=house.Title
	info["price"]=house.Price
	info["area_name"]=area.Name
	info["index_image_url"]=addDomain2Url(house.IndexImageUrl)
	info["room_count"]=house.RoomCount
	info["order_count"]=house.OrderCount
	info["address"]=house.Address
	info["user_avatar"]=addDomain2Url(user.AvatarUrl)
	info["create_time"]=displayTime(house.CreateTime)
	return info,nil
}

func getHouseDesc(stub shim.ChaincodeStubInterface,houseId string) (map[string]interface{},error) {
	hKey:="house_"+houseId
	houseData,err:=getInfo(stub,hKey)
	if err!=nil {
		return nil,err
	}

	var house House
	err=json.Unmarshal(houseData,&house)
	if err!=nil {
		return nil,err
	}

	uKey:="user_"+house.UserId
	userData,err:=getInfo(stub,uKey)
	if err!=nil {
		return nil,err
	}

	var user User
	err=json.Unmarshal(userData,&user)
	if err!=nil {
		return nil,err
	}

	var imgUrls []string
	var houseImage HouseImage
	for _,v:=range house.HouseImageIds{
		hiKey:="house_image_"+v
		houseImageData,err:=getInfo(stub,hiKey)
		if err!=nil {
			return nil,err
		}

		err=json.Unmarshal(houseImageData,&houseImage)
		if err!=nil {
			return nil,err
		}

		imgUrls=append(imgUrls,addDomain2Url(houseImage.Url))
	}

	var facilities []string
	var facility Facility
	for _,v:=range house.FacilityIds{
		fKey:="facility_"+v
		facilityData,err:=getInfo(stub,fKey)
		if err!=nil {
			return nil,err
		}

		err=json.Unmarshal(facilityData,&facility)
		if err!=nil {
			return nil,err
		}

		facilities=append(facilities,facility.Name)
	}

	var comments []map[string]interface{}
	var order Order
	for _,v:=range house.OrderIds{
		oKey:="order_"+v
		orderData,err:=getInfo(stub,oKey)
		if err!=nil {
			return nil,err
		}

		err=json.Unmarshal(orderData,&order)
		if err!=nil {
			return nil,err
		}

		if order.Comment=="" {
			continue
		}

		_uKey:="user_"+order.UserId
		_userData,err:=getInfo(stub,_uKey)
		if err!=nil {
			return nil,err
		}

		var _user User
		err=json.Unmarshal(_userData,&_user)
		if err!=nil {
			return nil,err
		}

		comment:=make(map[string]interface{})
		comment["username"]=_user.Name
		comment["comment"]=order.Comment
		comment["create_time"]=displayTime(order.CreateTime)
		comments=append(comments,comment)
	}

	info:=make(map[string]interface{})
	info["house_id"]=house.Id
	info["user_id"]=user.Id
	info["username"]=user.Name
	info["user_avatar"]=addDomain2Url(user.AvatarUrl)
	info["title"]=house.Title
	info["price"]=house.Price
	info["address"]=house.Address
	info["room_count"]=house.RoomCount
	info["acreage"]=house.Acreage
	info["uint"]=house.Unit
	info["capacity"]=house.Capacity
	info["beds"]=house.Beds
	info["deposit"]=house.Deposit
	info["min_days"]=house.MinDays
	info["max_days"]=house.MaxDays
	info["img_urls"]=imgUrls
	info["facilities"]=facilities
	info["comments"]=comments
	return info,nil
}

func getOrderInfo(stub shim.ChaincodeStubInterface,orderId string) (map[string]interface{},error) {
	oKey:="order_"+orderId
	orderData,err:=getInfo(stub,oKey)
	if err!=nil {
		return nil,err
	}

	var order Order
	err=json.Unmarshal(orderData,&order)
	if err!=nil {
		return nil,err
	}

	hKey:="house_"+order.HouseId
	houseData,err:=getInfo(stub,hKey)
	if err!=nil {
		return nil,err
	}

	var house House
	err=json.Unmarshal(houseData,&house)
	if err!=nil {
		return nil,err
	}

	info:=make(map[string]interface{})
	info["order_id"]=order.Id
	info["title"]=house.Title
	info["image_url"]=addDomain2Url(house.IndexImageUrl)
	info["begin_date"]=displayDate(order.BeginDate)
	info["end_date"]=displayDate(order.EndDate)
	info["create_time"]=displayTime(order.CreateTime)
	info["days"]=order.Days
	info["amount"]=order.Amount
	info["status"]=order.Status
	info["comment"]=order.Comment
	info["credit"]=order.Credit
	return info,nil
}

func checkArgsNum(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d argument(s) required",n))
	}
	return nil
}

func recordExists(stub shim.ChaincodeStubInterface,key string) (bool,error) {
	data,err:=stub.GetState(key)
	if err!=nil {
		return false,err
	}

	if data==nil {
		return false, nil
	} else {
		return true,nil
	}
}

func generateId(stub shim.ChaincodeStubInterface,prefix string) (string,error) {
	x:=time.Now().UnixNano()
	h256:=sha256.Sum256([]byte(strconv.Itoa(int(x))))
	h160:=ripemd160.New()
	h160.Write(h256[:])
	h:=h160.Sum(nil)
	hStr:=hex.EncodeToString(h)

	key:=prefix+hStr
	for  {
		exist,err:=recordExists(stub,key)
		if err!=nil {
			return "",err
		}
		if exist {
			return generateId(stub,prefix)
		} else {
			break
		}
	}

	return hStr,nil
}

/*
func (this *RentingChaincode) getInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	data,err:=getInfo(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}
 */

/*
func (this *RentingChaincode) addArea(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var area Area
	area.Id=args[0]
	area.Name=args[1]

	data,err:=json.Marshal(area)
	if err!=nil {
		return shim.Error(err.Error())
	}

	aKey:="area_"+area.Id
	err=stub.PutState(aKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalKey:="total_area"
	totalAreaData,err:=stub.GetState(totalKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var totalArea int
	if totalAreaData==nil {
		totalArea=1
	} else {
		totalArea,err=strconv.Atoi(string(totalAreaData))
		if err!=nil {
			return shim.Error(err.Error())
		}

		totalArea++
	}

	err=stub.PutState(totalKey,[]byte(strconv.Itoa(totalArea)))
	if err!=nil {
		return shim.Error(err.Error())
	}

	areaIdKey:="area_id_"+strconv.Itoa(totalArea)
	err=stub.PutState(areaIdKey,[]byte(area.Id))
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) delArea(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalAreaData,err:=getInfo(stub,"total_area")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalArea,err:=strconv.Atoi(string(totalAreaData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	for i:=1; i<=totalArea; i++ {
		aiKey:="area_id_"+strconv.Itoa(i)
		areaIdData,err:=getInfo(stub,aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		areaId:=string(areaIdData)

		aKey:="area_"+areaId
		err=stub.DelState(aKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		err=stub.DelState(aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}
	}

	err=stub.DelState("total_area")
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) addFacility(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var facility Facility
	facility.Id=args[0]
	facility.Name=args[1]

	data,err:=json.Marshal(facility)
	if err!=nil {
		return shim.Error(err.Error())
	}

	fKey:="facility_"+facility.Id
	err=stub.PutState(fKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalKey:="total_facility"
	totalFacilityData,err:=stub.GetState(totalKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var totalFacility int
	if totalFacilityData==nil {
		totalFacility=1
	} else {
		totalFacility,err=strconv.Atoi(string(totalFacilityData))
		if err!=nil {
			return shim.Error(err.Error())
		}

		totalFacility++
	}

	err=stub.PutState(totalKey,[]byte(strconv.Itoa(totalFacility)))
	if err!=nil {
		return shim.Error(err.Error())
	}

	facilityIdKey:="facility_id_"+strconv.Itoa(totalFacility)
	err=stub.PutState(facilityIdKey,[]byte(facility.Id))
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) delFacility(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalFacilityData,err:=getInfo(stub,"total_facility")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalFacility,err:=strconv.Atoi(string(totalFacilityData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	for i:=1; i<=totalFacility; i++ {
		aiKey:="facility_id_"+strconv.Itoa(i)
		facilityIdData,err:=getInfo(stub,aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		facilityId:=string(facilityIdData)

		aKey:="facility_"+facilityId
		err=stub.DelState(aKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		err=stub.DelState(aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}
	}

	err=stub.DelState("total_facility")
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
 */

func (this *RentingChaincode) getAreaList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalAreaData,err:=getInfo(stub,"total_area")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalArea,err:=strconv.Atoi(string(totalAreaData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	var areaList []Area
	for i:=1; i<=totalArea; i++ {
		areaIdData,err:=getInfo(stub,"area_id_"+strconv.Itoa(i))
		if err!=nil {
			return shim.Error(err.Error())
		}

		areaId:=string(areaIdData)

		aKey:="area_"+areaId
		areaData,err:=getInfo(stub,aKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		var area Area
		err=json.Unmarshal(areaData,&area)
		if err!=nil {
			return shim.Error(err.Error())
		}

		areaList=append(areaList,area)
	}

	data,err:=json.Marshal(areaList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

/*
func (this *RentingChaincode) getFacilityList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalFacilityData,err:=getInfo(stub,"total_facility")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalFacility,err:=strconv.Atoi(string(totalFacilityData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	var facilityList []Facility
	for i:=1; i<=totalFacility; i++ {
		facilityIdData,err:=getInfo(stub,"facility_id_"+strconv.Itoa(i))
		if err!=nil {
			return shim.Error(err.Error())
		}

		facilityId:=string(facilityIdData)

		aKey:="facility_"+facilityId
		facilityData,err:=getInfo(stub,aKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		var facility Facility
		err=json.Unmarshal(facilityData,&facility)
		if err!=nil {
			return shim.Error(err.Error())
		}

		facilityList=append(facilityList,facility)
	}

	data,err:=json.Marshal(facilityList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}
 */

/*
func (this *RentingChaincode) delHouse(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalHouseData,err:=getInfo(stub,"total_house")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalHouse,err:=strconv.Atoi(string(totalHouseData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	for i:=1; i<=totalHouse; i++ {
		aiKey:="house_id_"+strconv.Itoa(i)
		houseIdData,err:=getInfo(stub,aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseId:=string(houseIdData)

		aKey:="house_"+houseId
		err=stub.DelState(aKey)
		if err!=nil {
			return shim.Error(err.Error())
		}

		err=stub.DelState(aiKey)
		if err!=nil {
			return shim.Error(err.Error())
		}
	}

	err=stub.DelState("total_house")
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub,"17380416834", func(user *User) {
		user.HouseIds=nil
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub,"15928788003", func(user *User) {
		user.HouseIds=nil
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) delOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=stub.DelState("order_21ed384369802907a1f95e58bac244f21728cfa9")
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub,"15928788003", func(user *User) {
		user.OrderIds=nil
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	_houseData,err:=getInfo(stub,"house_851817007e043f154e95f268355052dd98ab7dd7")
	if err!=nil {
		return shim.Error(err.Error())
	}

	var house House
	err=json.Unmarshal(_houseData,&house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.OrderIds=nil

	houseData,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState("house_851817007e043f154e95f268355052dd98ab7dd7",houseData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
 */

func (this *RentingChaincode) generateUserId(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	userId,err:=generateId(stub,"user_")
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(userId))
}

func (this *RentingChaincode) register(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var user User

	user.Id=args[0]
	key:="user_"+user.Id
	exist,err:=recordExists(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if exist {
		return shim.Error("User exists")
	}

	user.Mobile=args[1]
	user.Name=user.Mobile
	mKey:="user_"+user.Mobile
	exist,err=recordExists(stub,mKey)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if exist {
		return shim.Error("Mobile has been registered")
	}

	user.PublicKey=args[2]

	data,err:=json.Marshal(user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(key,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(mKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

/*
func (this *RentingChaincode) delUser(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="user_"+args[0]
	user,err:=getUser(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	iKey:="user_"+user.Id
	err=stub.DelState(iKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	mKey:="user_"+user.Mobile
	err=stub.DelState(mKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) setUserMobile(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="user_"+args[0]
	user,err:=getUser(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	mKey:="user_"+user.Mobile
	err=stub.DelState(mKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user.Mobile=args[1]
	data,err:=json.Marshal(user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	iKey:="user_"+user.Id
	err=stub.PutState(iKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	mKey="user_"+user.Mobile
	err=stub.PutState(mKey,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
 */

func (this *RentingChaincode) getUserInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="user_"+args[0]
	user,err:=getUser(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	info:=make(map[string]interface{})
	info["user_id"]=user.Id
	info["name"]=user.Name
	info["mobile"]=user.Mobile
	info["real_name"]=user.RealName
	info["id_card"]=user.IdCard
	info["avatar_url"]=addDomain2Url(user.AvatarUrl)

	data,err:=json.Marshal(info)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) getUserPublicKey(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="user_"+args[0]
	user,err:=getUser(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(user.PublicKey))
}

func setUserInfo(stub shim.ChaincodeStubInterface,key string,setUser func(*User)) error {
	key="user_"+key
	_data,err:=getInfo(stub,key)
	if err!=nil {
		return err
	}

	var user User
	err=json.Unmarshal(_data,&user)
	if err!=nil {
		return err
	}

	setUser(&user)
	data,err:=json.Marshal(user)
	if err!=nil {
		return err
	}

	idKey:="user_"+user.Id
	err=stub.PutState(idKey,data)
	if err!=nil {
		return err
	}

	mKey:="user_"+user.Mobile
	err=stub.PutState(mKey,data)
	if err!=nil {
		return err
	}

	return nil
}

func (this *RentingChaincode) setUserAvatar(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub,args[0], func(user *User) {
		user.AvatarUrl=args[1]
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) setUserName(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=setUserInfo(stub,args[0], func(user *User) {
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

	err=setUserInfo(stub,args[0], func(user *User) {
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

	houseId,err:=generateId(stub,"house_")
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Id=houseId

	house.Title=args[0]

	price,err:=strconv.Atoi(args[1])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Price=price*100

	house.AreaId=args[2]
	house.Address=args[3]

	roomCount,err:=strconv.Atoi(args[4])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.RoomCount=roomCount

	acreage,err:=strconv.Atoi(args[5])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Acreage=acreage

	house.Unit=args[6]

	capacity,err:=strconv.Atoi(args[7])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Capacity=capacity

	house.Beds=args[8]

	deposit,err:=strconv.Atoi(args[9])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.Deposit=deposit*100

	minDays,err:=strconv.Atoi(args[10])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.MinDays=minDays

	maxDays,err:=strconv.Atoi(args[11])
	if err!=nil {
		return shim.Error(err.Error())
	}
	house.MaxDays=maxDays

	facilities:=strings.Split(args[12],",")
	house.FacilityIds=facilities

	house.UserId=args[13]
	err=setUserInfo(stub,house.UserId, func(user *User) {
		user.HouseIds=append(user.HouseIds,house.Id)
	})
	if err!=nil {
		return shim.Error(err.Error())
	}

	loc,_:=time.LoadLocation(TIME_ZONE)
	house.CreateTime=time.Now().In(loc)

	data,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="house_"+houseId
	err=stub.PutState(key,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalKey:="total_house"
	totalHouseData,err:=stub.GetState(totalKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var totalHouse int
	if totalHouseData==nil {
		totalHouse=1
	} else {
		totalHouse,err=strconv.Atoi(string(totalHouseData))
		if err!=nil {
			return shim.Error(err.Error())
		}

		totalHouse++
	}

	err=stub.PutState(totalKey,[]byte(strconv.Itoa(totalHouse)))
	if err!=nil {
		return shim.Error(err.Error())
	}

	houseIdKey:="house_id_"+strconv.Itoa(totalHouse)
	err=stub.PutState(houseIdKey,[]byte(houseId))
	if err!=nil {
		return shim.Error(err.Error())
	}

	/*houseAreaIdListKey:="house_area_id_list_"+house.AreaId
	_houseAreaIdListData,err:=stub.GetState(houseAreaIdListKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var houseAreaIdList []string
	if _houseAreaIdListData==nil {
		houseAreaIdList=[]string{houseId}
	} else {
		err=json.Unmarshal(_houseAreaIdListData,&houseAreaIdList)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseAreaIdList=append([]string{houseId},houseAreaIdList...)
	}

	houseAreaIdListData,err:=json.Marshal(houseAreaIdList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(houseAreaIdListKey,houseAreaIdListData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	houseDateIdListKey:="house_date_id_list_"+house.CreateTime.In(loc).Format("2006-01-02")
	_houseDateIdListData,err:=stub.GetState(houseDateIdListKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var houseDateIdList []string
	if _houseDateIdListData==nil {
		houseDateIdList=[]string{houseId}
	} else {
		err=json.Unmarshal(_houseDateIdListData,&houseDateIdList)
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseDateIdList=append([]string{houseId},houseDateIdList...)
	}

	houseDateIdListData,err:=json.Marshal(houseDateIdList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(houseDateIdListKey,houseDateIdListData)
	if err!=nil {
		return shim.Error(err.Error())
	}*/

	return shim.Success(nil)
}

func (this *RentingChaincode) getHouseList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	key:="user_"+args[0]
	_data,err:=getInfo(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var user User
	err=json.Unmarshal(_data,&user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var houseInfoList []map[string]interface{}
	for _,v:=range user.HouseIds{
		houseInfo,err:=getHouseInfo(stub,v)
		if err!=nil {
			return shim.Error(err.Error())
		}
		houseInfoList=append(houseInfoList,houseInfo)
	}

	data,err:=json.Marshal(houseInfoList)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) addHouseImage(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,3)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="house_"+args[0]
	_data,err:=getInfo(stub,key)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var house House
	err=json.Unmarshal(_data,&house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.UserId!=args[2] {
		return shim.Error("invalid user")
	}

	var houseImage HouseImage

	houseImageId,err:=generateId(stub,"house_image_")
	if err!=nil {
		return shim.Error(err.Error())
	}
	houseImage.Id=houseImageId

	houseImage.Url=args[1]

	houseImageData,err:=json.Marshal(houseImage)
	if err!=nil {
		return shim.Error(err.Error())
	}

	hiKey:="house_image_"+houseImageId
	err=stub.PutState(hiKey,houseImageData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.IndexImageUrl=="" {
		house.IndexImageUrl=args[1]
	}
	house.HouseImageIds=append(house.HouseImageIds,houseImageId)

	data,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(key,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getHouseInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}

	info,err:=getHouseDesc(stub,args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}

	data,err:=json.Marshal(info)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) getIndexHouseList(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	totalHouseData,err:=getInfo(stub,"total_house")
	if err!=nil {
		return shim.Error(err.Error())
	}

	totalHouse,err:=strconv.Atoi(string(totalHouseData))
	if err!=nil {
		return shim.Error(err.Error())
	}

	var houseList []map[string]interface{}
	for i:=totalHouse; i>totalHouse-INDEX_HOUSE_LIST_SIZE && i>0; i-- {
		houseIdData,err:=getInfo(stub,"house_id_"+strconv.Itoa(i))
		if err!=nil {
			return shim.Error(err.Error())
		}

		houseId:=string(houseIdData)

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

func (this *RentingChaincode) addOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,4)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var order Order
	order.UserId=args[0]

	orderId,err:=generateId(stub,"order_")
	if err!=nil {
		return shim.Error(err.Error())
	}
	order.Id=orderId

	order.HouseId=args[1]

	startDate,err:=time.Parse("2006-01-02",args[2])
	if err!=nil {
		return shim.Error(err.Error())
	}
	order.BeginDate=startDate

	endDate,err:=time.Parse("2006-01-02",args[3])
	if err!=nil {
		return shim.Error(err.Error())
	}
	order.EndDate=endDate

	days:=int(endDate.Sub(startDate).Seconds()/60/60/24)+1
	order.Days=days

	hKey:="house_"+order.HouseId
	_houseData,err:=getInfo(stub,hKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var house House
	err=json.Unmarshal(_houseData,&house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if order.UserId==house.UserId {
		return shim.Error("invalid user")
	}

	order.HousePrice=house.Price
	order.Amount=order.HousePrice*order.Days

	order.Status=ORDER_STATUS_WAIT_ACCEPT
	order.Comment=""

	loc,_:=time.LoadLocation(TIME_ZONE)
	order.CreateTime=time.Now().In(loc)

	order.Credit=false

	data,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	key:="order_"+orderId
	err=stub.PutState(key,data)
	if err!=nil {
		return shim.Error(err.Error())
	}

	house.OrderIds=append(house.OrderIds,orderId)
	houseData,err:=json.Marshal(house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(hKey,houseData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	uKey:="user_"+order.UserId
	_userData,err:=getInfo(stub,uKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var user User
	err=json.Unmarshal(_userData,&user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	user.OrderIds=append(user.OrderIds,orderId)
	userData,err:=json.Marshal(user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(uKey,userData)
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

	uKey:="user_"+args[0]
	_data,err:=getInfo(stub,uKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var user User
	err=json.Unmarshal(_data,&user)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var orderInfoList []map[string]interface{}
	if args[1]=="1" {
		for _,v:=range user.HouseIds{
			hKey:="house_"+v
			houseData,err:=getInfo(stub,hKey)
			if err!=nil {
				return shim.Error(err.Error())
			}

			var house House
			err=json.Unmarshal(houseData,&house)
			if err!=nil {
				return shim.Error(err.Error())
			}

			for _,vv:=range house.OrderIds{
				orderInfo,err:=getOrderInfo(stub,vv)
				if err!=nil {
					return shim.Error(err.Error())
				}
				orderInfoList=append([]map[string]interface{}{orderInfo},orderInfoList...)
			}
		}
	} else {
		for _,v:=range user.OrderIds{
			orderInfo,err:=getOrderInfo(stub,v)
			if err!=nil {
				return shim.Error(err.Error())
			}
			orderInfoList=append([]map[string]interface{}{orderInfo},orderInfoList...)
		}
	}

	data,err:=json.Marshal(orderInfoList)
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

	oKey:="order_"+args[0]
	_orderData,err:=getInfo(stub,oKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var order Order
	err=json.Unmarshal(_orderData,&order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	hKey:="house_"+order.HouseId
	_houseData,err:=getInfo(stub,hKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var house House
	err=json.Unmarshal(_houseData,&house)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if house.UserId!=args[1] {
		return shim.Error("NO DATA")
	}
	if order.Status!=ORDER_STATUS_WAIT_ACCEPT {
		return shim.Error("NO DATA")
	}

	if args[2]=="1" {
		order.Status=ORDER_STATUS_REJECTED
	} else {
		order.Status=ORDER_STATUS_WAIT_COMMENT
		house.OrderCount++
	}

	orderData,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(oKey,orderData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if order.Status==ORDER_STATUS_WAIT_COMMENT {
		houseData,err:=json.Marshal(house)
		if err!=nil {
			return shim.Error(err.Error())
		}

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

	oKey:="order_"+args[0]
	_orderData,err:=getInfo(stub,oKey)
	if err!=nil {
		return shim.Error(err.Error())
	}

	var order Order
	err=json.Unmarshal(_orderData,&order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	if order.UserId!=args[1] {
		return shim.Error("NO DATA")
	}
	if order.Status!=ORDER_STATUS_WAIT_COMMENT {
		return shim.Error("NO DATA")
	}

	order.Comment=args[2]
	order.Status=ORDER_STATUS_COMPLETE
	orderData,err:=json.Marshal(order)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err=stub.PutState(oKey,orderData)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *RentingChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	switch fn {
	/*
	case "getInfo":
		return this.getInfo(stub,args)
	 */
	/*
	case "addArea":
		return this.addArea(stub,args)
	case "delArea":
		return this.delArea(stub,args)
	case "addFacility":
		return this.addFacility(stub,args)
	case "delFacility":
		return this.delFacility(stub,args)
	case "getFacilityList":
		return this.getFacilityList(stub,args)
	 */
	case "getAreaList":
		return this.getAreaList(stub,args)
	case "generateUserId":
		return this.generateUserId(stub,args)
	case "getUserInfo":
		return this.getUserInfo(stub,args)
	case "getUserPublicKey":
		return this.getUserPublicKey(stub,args)
	case "register":
		return this.register(stub,args)
	/*
	case "delUser":
		return this.delUser(stub,args)
	case "setUserMobile":
		return this.setUserMobile(stub,args)
	 */
	/*
	case "delHouse":
		return this.delHouse(stub,args)
	case "delOrder":
		return this.delOrder(stub,args)
	 */
	case "setUserAvatar":
		return this.setUserAvatar(stub,args)
	case "setUserName":
		return this.setUserName(stub,args)
	case "auth":
		return this.auth(stub,args)
	case "addHouse":
		return this.addHouse(stub,args)
	case "getHouseList":
		return this.getHouseList(stub,args)
	case "addHouseImage":
		return this.addHouseImage(stub,args)
	case "getHouseInfo":
		return this.getHouseInfo(stub,args)
	case "getIndexHouseList":
		return this.getIndexHouseList(stub,args)
	case "addOrder":
		return this.addOrder(stub,args)
	case "getOrderList":
		return this.getOrderList(stub,args)
	case "handleOrder":
		return this.handleOrder(stub,args)
	case "comment":
		return this.comment(stub,args)
	default:
		return shim.Error("invalid parameter")
	}
}

func main()  {
	shim.Start(new(RentingChaincode))
}
