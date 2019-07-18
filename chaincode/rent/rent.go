package main

import (
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
	"math"
)

type RentChaincode struct {

}

const (
	UserPrefix = "User"
	HousePrefix = "House"
	AreaPrefix = "Area"
	FacilityPrefix = "Facility"
	HouseImagePrefix = "HouseImage"
	OrderPrefix = "Order"

	OrderWaitAccept= "WAIT_ACCEPT"
	OrderWaitComment= "WAIT_COMMENT"
	OrderComplete= "COMPLETE"
	OrderRejected= "REJECTED"

	TimeZone = "Asia/Saigon"
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"

	FDFSHost = "149.28.210.102"
	FDFSPort = 8888

	HouseListPageSize = 10
)

func find(stub shim.ChaincodeStubInterface, prefix string, id string) ([]byte, error) {
	key, err := stub.CreateCompositeKey(prefix, []string{id})
	if err != nil {
		return nil, err
	}

	data, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getQueryResult(stub shim.ChaincodeStubInterface, query map[string]interface{}) (shim.StateQueryIteratorInterface, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	iter, err := stub.GetQueryResult(string(queryJson))
	if err != nil {
		return nil, err
	}

	return iter, nil
}

func generateId(stub shim.ChaincodeStubInterface, prefix string, length int) (string, error) {
	id := ""
	a := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	n := len(a)
	for {
		id = ""
		for i := 0; i < length; i++ {
			if i > 0 {
				id += a[rand.Intn(n)]
			} else {
				id += a[rand.Intn(n-1)]
			}
		}

		data, err := find(stub, prefix, id)
		if err != nil {
			return "", err
		}
		if data == nil {
			break
		}
	}

	return id, nil
}

func put(stub shim.ChaincodeStubInterface, prefix string, id string, obj interface{}) error {
	key, err := stub.CreateCompositeKey(prefix, []string{id})
	if err != nil {
		return err
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = stub.PutState(key, data)
	if err != nil {
		return err
	}

	return nil
}

func add(stub shim.ChaincodeStubInterface, prefix string, id string, obj interface{}) error {
	data, err := find(stub, prefix, id)
	if err != nil {
		return err
	}
	if data != nil {
		err = fmt.Errorf("key already exists")
		return err
	}

	err = put(stub, prefix, id, obj)
	if err != nil {
		return err
	}

	return nil
}

func update(stub shim.ChaincodeStubInterface, prefix string, id string, obj interface{}) error {
	data, err := find(stub, prefix, id)
	if err != nil {
		return err
	}
	if data == nil {
		err = fmt.Errorf("key doesn't exist")
		return err
	}

	err = put(stub, prefix, id, obj)
	if err != nil {
		return err
	}

	return nil
}

func del(stub shim.ChaincodeStubInterface, prefix string, id string) (err error) {
	key, err := stub.CreateCompositeKey(prefix, []string{id})
	if err != nil {
		return
	}

	err = stub.DelState(key)
	if err != nil {
		return
	}

	return
}

func checkArgs(args []string, n int) error {
	if len(args) < n {
		err := fmt.Errorf("%d parameter(s) required", n)
		return err
	}

	return nil
}

func toTimeZone(t time.Time) time.Time {
	loc, _ := time.LoadLocation(TimeZone)
	zt := t.In(loc)
	return zt
}

func today() time.Time {
	now := toTimeZone(time.Now()).Format(DateLayout)
	loc, _ := time.LoadLocation(TimeZone)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", now+" 00:00:00", loc)
	return t
}

func toFDFSUrl(url string) string {
	if url == "" {
		return ""
	}

	fdfsUrl := "http://"+FDFSHost+":"+strconv.Itoa(FDFSPort)+"/"+url
	return fdfsUrl
}

type Area struct {
	AreaId string
	ObjectType string
	AreaName string
}

func getAreaName(stub shim.ChaincodeStubInterface, areaId string) (string, error) {
	data, err := find(stub, AreaPrefix, areaId)
	if err != nil {
		return "", err
	}
	if data == nil {
		err = fmt.Errorf("area not found")
		return "", err
	}

	var area Area
	err = json.Unmarshal(data, &area)
	if err != nil {
		return "", err
	}

	areaName := area.AreaName
	return areaName, nil
}

func (this *RentChaincode) addArea(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	areaId := args[0]
	areaName := args[1]

	var area Area
	area.AreaId = areaId
	area.ObjectType = AreaPrefix
	area.AreaName = areaName

	err = add(stub, AreaPrefix, areaId, area)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getAreaList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"use_index": []string{"_design/areaDoc", "area"},
		"selector": map[string]interface{}{
			"ObjectType": AreaPrefix,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"AreaId": "asc"},
		},
		"fields": []string{"AreaId", "AreaName"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	areaList := []map[string]interface{}{}
	for iter.HasNext() {
		item, _err := iter.Next()
		if _err != nil {
			continue
		}

		var _area Area
		_err = json.Unmarshal(item.Value, &_area)
		if _err != nil {
			continue
		}

		area := map[string]interface{}{
			"area_id": _area.AreaId,
			"area_name": _area.AreaName,
		}
		areaList = append(areaList, area)
	}

	data, err := json.Marshal(areaList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

type Facility struct {
	FacilityId string
	ObjectType string
	FacilityName string
}

func (this *RentChaincode) addFacility(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	facilityId := args[0]
	facilityName := args[1]

	var facility Facility
	facility.FacilityId = facilityId
	facility.ObjectType = FacilityPrefix
	facility.FacilityName = facilityName

	err = add(stub, FacilityPrefix, facilityId, facility)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getFacilityList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"use_index": []string{"_design/facilityDoc", "facility"},
		"selector": map[string]interface{}{
			"ObjectType": FacilityPrefix,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"FacilityId": "asc"},
		},
		"fields": []string{"FacilityId", "FacilityName"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	facilityList := []map[string]interface{}{}
	for iter.HasNext() {
		item, _err := iter.Next()
		if _err != nil {
			continue
		}

		var _facility Facility
		_err = json.Unmarshal(item.Value, &_facility)
		if _err != nil {
			continue
		}

		facility := map[string]interface{}{
			"facility_id": _facility.FacilityId,
			"facility_name": _facility.FacilityName,
		}
		facilityList = append(facilityList, facility)
	}

	data, err := json.Marshal(facilityList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func getFacilityList(stub shim.ChaincodeStubInterface) (map[string]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/facilityDoc", "facility"},
		"selector": map[string]interface{}{
			"ObjectType": FacilityPrefix,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"FacilityId": "asc"},
		},
		"fields": []string{"FacilityId", "FacilityName"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	facilityList := make(map[string]string)
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var facility Facility
		err = json.Unmarshal(item.Value, &facility)
		if err != nil {
			continue
		}

		facilityId := facility.FacilityId
		facilityName := facility.FacilityName
		facilityList[facilityId] = facilityName
	}

	return facilityList, nil
}

type User struct {
	UserId string
	ObjectType string
	UserName string
	Address string
	Mobile string
	RealName string
	IdCard string
	AvatarUrl string
	CreateTime time.Time
}

func findUserById(stub shim.ChaincodeStubInterface, id string) (User, error) {
	data, err := find(stub, UserPrefix, id)
	if err != nil {
		return User{}, err
	}
	if data == nil {
		return User{}, nil
	}

	user := User{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func findUserByMobile(stub shim.ChaincodeStubInterface, mobile string) (User, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/userDoc", "user"},
		"selector": map[string]interface{}{
			"ObjectType": UserPrefix,
			"Mobile": mobile,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return User{}, err
	}
	defer iter.Close()

	user := User{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		err = json.Unmarshal(item.Value, &user)
		if err != nil {
			continue
		}
	}

	return user, nil
}

func findUserByAddress(stub shim.ChaincodeStubInterface, address string) (User, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/userAddressDoc", "userAddress"},
		"selector": map[string]interface{}{
			"ObjectType": UserPrefix,
			"Address": address,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return User{}, err
	}
	defer iter.Close()

	user := User{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		err = json.Unmarshal(item.Value, &user)
		if err != nil {
			continue
		}
	}

	return user, nil
}

func updateUser(stub shim.ChaincodeStubInterface, address string, setUser func(*User)) error {
	user, err := findUserByAddress(stub, address)
	if err != nil {
		return err
	}
	if user.UserId == "" {
		err = fmt.Errorf("user does not exist")
		return err
	}

	setUser(&user)
	err = update(stub, UserPrefix, user.UserId, user)
	if err != nil {
		return err
	}

	return nil
}

func (this *RentChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	mobile := args[0]
	address := args[1]

	_user, err := findUserByMobile(stub, mobile)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _user.UserId != "" {
		return shim.Error("mobile already registered")
	}

	_user, err = findUserByAddress(stub, address)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _user.UserId != "" {
		return shim.Error("address already registered")
	}

	userId, err := generateId(stub, UserPrefix, 10)
	if err != nil {
		return shim.Error(err.Error())
	}

	var user User
	user.UserId = userId
	user.ObjectType = UserPrefix
	user.UserName = mobile
	user.Address = address
	user.Mobile = mobile
	user.CreateTime = toTimeZone(time.Now())

	err = add(stub, UserPrefix, userId, user)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) updateUserAvatar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	address := args[0]
	avatarUrl := args[1]

	err = updateUser(stub, address, func(user *User) {
		user.AvatarUrl = avatarUrl
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) rename(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	address := args[0]
	newName := args[1]

	err = updateUser(stub, address, func(user *User) {
		user.UserName = newName
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) auth(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	address := args[0]
	realName := args[1]
	idCard := args[2]

	err = updateUser(stub, address, func(user *User) {
		user.RealName = realName
		user.IdCard = idCard
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	address := args[0]
	user, err := findUserByAddress(stub, address)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	_user := map[string]interface{}{
		"user_id": user.UserId,
		"username": user.UserName,
		"address": user.Address,
		"mobile": user.Mobile,
		"real_name": user.RealName,
		"id_card": user.IdCard,
		"avatar_url": toFDFSUrl(user.AvatarUrl),
		"create_time": user.CreateTime.Format(TimeLayout),
	}

	data, err := json.Marshal(_user)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

type House struct {
	HouseId string
	ObjectType string
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
	FacilityId []string
}

func getHouse(stub shim.ChaincodeStubInterface, houseId string) (House, error) {
	data, err := find(stub, HousePrefix, houseId)
	if err != nil {
		return House{}, err
	}
	if data == nil {
		err = fmt.Errorf("house not found")
		return House{}, err
	}

	house := House{}
	err = json.Unmarshal(data, &house)
	if err != nil {
		return House{}, err
	}

	if house.HouseId == "" {
		err = fmt.Errorf("house not found")
		return House{}, err
	}

	return house, nil
}

func getHouseInfo(stub shim.ChaincodeStubInterface, houseId string) (map[string]interface{}, error) {
	house, err := getHouse(stub, houseId)
	if err != nil {
		return nil, err
	}

	areaName, err := getAreaName(stub, house.AreaId)
	if err != nil {
		return nil, err
	}

	user, err := findUserById(stub, house.UserId)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"house_id": house.HouseId,
		"title": house.Title,
		"price": house.Price,
		"area_name": areaName,
		"index_image_url": toFDFSUrl(house.IndexImageUrl),
		"room_count": house.RoomCount,
		"order_count": house.OrderCount,
		"address": house.Address,
		"user_avatar": toFDFSUrl(user.AvatarUrl),
		"create_time": house.CreateTime.Format(TimeLayout),
	}

	return info, nil
}

func getHouseDesc(stub shim.ChaincodeStubInterface, houseId string) (map[string]interface{}, error) {
	house, err := getHouse(stub, houseId)
	if err != nil {
		return nil, err
	}

	user, err := findUserById(stub, house.UserId)
	if err != nil {
		return nil, err
	}

	imageUrlList, err := getHouseImageUrlList(stub, houseId)
	if err != nil {
		return nil, err
	}

	allFacilityList, err := getFacilityList(stub)
	if err != nil {
		return nil, err
	}

	var facilityList  []string
	for _, facilityId := range house.FacilityId {
		facilityList = append(facilityList, allFacilityList[facilityId])
	}

	orderIdList, err := getHouseOrderIdList(stub, houseId)
	if err != nil {
		return nil, err
	}

	var commentList []map[string]interface{}
	for _, orderId := range orderIdList{
		order, err := getOrder(stub, orderId)
		if err != nil {
			continue
		}
		if order.Status != OrderComplete || order.Comment == "" {
			continue
		}

		user, err := findUserById(stub, order.UserId)
		if err != nil {
			continue
		}

		comment := map[string]interface{}{
			"username": user.UserName,
			"comment": order.Comment,
			"create_time": order.CreateTime.Format(TimeLayout),
		}
		commentList = append(commentList, comment)
	}

	info := map[string]interface{}{
		"house_id": house.HouseId,
		"user_id": house.UserId,
		"username": user.UserName,
		"user_avatar": toFDFSUrl(user.AvatarUrl),
		"title": house.Title,
		"price": house.Price,
		"address": house.Address,
		"room_count": house.RoomCount,
		"acreage": house.Acreage,
		"unit": house.Unit,
		"capacity": house.Capacity,
		"beds": house.Beds,
		"deposit": house.Deposit,
		"min_days": house.MinDays,
		"max_days": house.MaxDays,
		"img_urls": imageUrlList,
		"facilities": facilityList,
		"comments": commentList,
	}

	return info, nil
}

func getUserHouseIdList(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/houseUserDoc", "houseUser"},
		"selector": map[string]interface{}{
			"ObjectType": HousePrefix,
			"UserId": userId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	idList := []string{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var house House
		err = json.Unmarshal(item.Value, &house)
		if err != nil {
			continue
		}

		idList = append(idList, house.HouseId)
	}

	return idList, nil
}

func (this *RentChaincode) addHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	title := args[1]

	price, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	areaId := args[3]
	address := args[4]

	roomCount, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error(err.Error())
	}

	acreage, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error(err.Error())
	}

	unit := args[7]

	capacity, err := strconv.Atoi(args[8])
	if err != nil {
		return shim.Error(err.Error())
	}

	beds := args[9]

	deposit, err := strconv.Atoi(args[10])
	if err != nil {
		return shim.Error(err.Error())
	}

	minDays, err := strconv.Atoi(args[11])
	if err != nil {
		return shim.Error(err.Error())
	}

	maxDays, err := strconv.Atoi(args[12])
	if err != nil {
		return shim.Error(err.Error())
	}

	facilityId := strings.Split(args[13], ",")

	houseId, err := generateId(stub, HousePrefix, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	var house House
	house.ObjectType = HousePrefix
	house.HouseId = houseId
	house.Title = title
	house.Price = price
	house.AreaId = areaId
	house.Address = address
	house.RoomCount = roomCount
	house.Acreage = acreage
	house.Unit = unit
	house.Capacity = capacity
	house.Beds = beds
	house.Deposit = deposit
	house.MinDays = minDays
	house.MaxDays = maxDays
	house.UserId = user.UserId
	house.FacilityId = facilityId
	house.CreateTime = toTimeZone(time.Now())

	err = add(stub, HousePrefix, houseId, house)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) updateHouseImage(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	houseId := args[1]
	imageUrl := args[2]

	house, err := getHouse(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if user.UserId != house.UserId {
		return shim.Error("forbidden")
	}

	houseImageId, err := generateId(stub, HouseImagePrefix, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	var houseImage HouseImage
	houseImage.ObjectType = HouseImagePrefix
	houseImage.HouseImageId = houseImageId
	houseImage.HouseId = houseId
	houseImage.Url = imageUrl

	err = add(stub, HouseImagePrefix, houseImageId, houseImage)
	if err != nil {
		return shim.Error(err.Error())
	}

	if house.IndexImageUrl == "" {
		house.IndexImageUrl = imageUrl
		err = update(stub, HousePrefix, houseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getUserHouseList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	userId := args[0]
	houseIdList, err := getUserHouseIdList(stub, userId)
	if err != nil {
		return shim.Error(err.Error())
	}

	var houseList []map[string]interface{}
	for _, houseId := range houseIdList{
		house, err := getHouseInfo(stub, houseId)
		if err != nil {
			continue
		}
		houseList = append(houseList, house)
	}

	data, err := json.Marshal(houseList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) getHouseDetail(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	houseId := args[0]
	info, err := getHouseDesc(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := json.Marshal(info)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) searchHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 5)
	if err != nil {
		return shim.Error(err.Error())
	}

	areaId := args[0]
	start := args[1]
	end := args[2]

	page, err := strconv.Atoi(args[3])
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(args[4])
	if err != nil || pageSize < 1 {
		pageSize = HouseListPageSize
	}

	query := map[string]interface{}{
		"use_index": []string{"_design/houseDoc", "house"},
		"selector": map[string]interface{}{
			"ObjectType": HousePrefix,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{
				"CreateTime": "desc",
			},
		},
	}

	selector := query["selector"].(map[string]interface{})

	if areaId != "" {
		selector["AreaId"] = areaId
		query["use_index"] = []string{"_design/houseAreaDoc", "houseArea"}
	}

	if start != "" && end != "" {
		_start, err := time.Parse(TimeLayout, start+" 00:00:00")
		if err != nil {
			return shim.Error(err.Error())
		}

		_end, err := time.Parse(TimeLayout, end+" 23:59:59")
		if err != nil {
			return shim.Error(err.Error())
		}

		selector["$and"] = []map[string]interface{}{
			map[string]interface{}{
				"CreateTime": map[string]interface{}{
					"$gte": _start,
				},
			},
			map[string]interface{}{
				"CreateTime": map[string]interface{}{
					"$lte": _end,
				},
			},
		}
	}

	if start != "" && end == "" {
		_start, err := time.Parse(TimeLayout, start+" 00:00:00")
		if err != nil {
			return shim.Error(err.Error())
		}

		selector["CreateTime"] = map[string]interface{}{
			"$gte": _start,
		}
	}

	if start == "" && end != "" {
		_end, err := time.Parse(TimeLayout, end+" 23:59:59")
		if err != nil {
			return shim.Error(err.Error())
		}

		selector["CreateTime"] = map[string]interface{}{
			"$lte": _end,
		}
	}

	query["selector"] = selector
	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	info := make(map[string]interface{})

	total := 0
	var houseList []map[string]interface{}
	offset := pageSize*(page-1)
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		total += 1
		if total <= offset || total > offset+pageSize {
			continue
		}

		var _house House
		err = json.Unmarshal(item.Value, &_house)
		if err != nil {
			continue
		}

		house, err := getHouseInfo(stub, _house.HouseId)
		if err != nil {
			continue
		}

		houseList = append(houseList, house)
	}

	info["houses"] = houseList
	info["total"] = total
	info["current_page"] = page
	info["total_page"] = int(math.Ceil(float64(total)/float64(pageSize)))

	data, err := json.Marshal(info)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

type HouseImage struct {
	HouseImageId string
	ObjectType string
	HouseId string
	Url string
}

func getHouseImageUrlList(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/houseImageDoc", "houseImage"},
		"selector": map[string]interface{}{
			"ObjectType": HouseImagePrefix,
			"HouseId": houseId,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	urlList := []string{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var houseImage HouseImage
		err = json.Unmarshal(item.Value, &houseImage)
		if err != nil {
			continue
		}

		urlList = append(urlList, toFDFSUrl(houseImage.Url))
	}

	return urlList, nil
}

type Order struct {
	OrderId string
	ObjectType string
	BeginDate time.Time
	EndDate time.Time
	Days int
	HousePrice int
	Amount int
	Status string
	Comment string
	CreateTime time.Time

	HouseId string
	UserId string
}

func getOrder(stub shim.ChaincodeStubInterface, orderId string) (Order, error) {
	data, err := find(stub, OrderPrefix, orderId)
	if err != nil {
		return Order{}, err
	}
	if data == nil {
		err = fmt.Errorf("order not found")
		return Order{}, err
	}

	order := Order{}
	err = json.Unmarshal(data, &order)
	if err != nil {
		return Order{}, err
	}

	if order.OrderId == "" {
		err = fmt.Errorf("order not found")
		return Order{}, err
	}

	return order, nil
}

func getOrderInfo(stub shim.ChaincodeStubInterface, orderId string) (map[string]interface{}, error) {
	order, err := getOrder(stub, orderId)
	if err != nil {
		return nil, err
	}

	house, err := getHouse(stub, order.HouseId)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"order_id": order.OrderId,
		"title": house.Title,
		"image_url": toFDFSUrl(house.IndexImageUrl),
		"begin_date": order.BeginDate.Format(DateLayout),
		"end_date": order.EndDate.Format(DateLayout),
		"create_time": order.CreateTime.Format(TimeLayout),
		"days": order.Days,
		"amount": order.Amount,
		"status": order.Status,
		"comment": order.Comment,
	}

	return info, nil
}

func getHouseOrderIdList(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/orderHouseDoc", "orderHouse"},
		"selector": map[string]interface{}{
			"ObjectType": OrderPrefix,
			"HouseId": houseId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	idList := []string{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var order Order
		err = json.Unmarshal(item.Value, &order)
		if err != nil {
			continue
		}

		idList = append(idList, order.OrderId)
	}

	return idList, nil
}

func getUserOrderIdList(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/orderUserDoc", "orderUser"},
		"selector": map[string]interface{}{
			"ObjectType": OrderPrefix,
			"UserId": userId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	idList := []string{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var order Order
		err = json.Unmarshal(item.Value, &order)
		if err != nil {
			continue
		}

		idList = append(idList, order.OrderId)
	}

	return idList, nil
}

func (this *RentChaincode) addOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 4)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	houseId := args[1]
	if houseId == "" {
		return shim.Error("please input house_id")
	}

	loc, _ := time.LoadLocation(TimeZone)
	start, _ := time.ParseInLocation(TimeLayout, args[2]+" 00:00:00", loc)
	end, _ := time.ParseInLocation(TimeLayout, args[3]+" 00:00:00", loc)

	if start.Sub(end) > 0 {
		return shim.Error("end date no earlier than start date")
	}

	if today().Sub(start) > 0 {
		return shim.Error("start date no earlier than today")
	}

	house, err := getHouse(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if user.UserId == house.UserId {
		return shim.Error("forbidden")
	}

	days := int(end.Sub(start).Seconds()/86400)+1
	if days < house.MinDays {
		return shim.Error("no less than "+strconv.Itoa(house.MinDays)+" days")
	}

	if house.MaxDays > 0 && days > house.MaxDays {
		return shim.Error("no more than "+strconv.Itoa(house.MaxDays)+" days")
	}

	housePrice := house.Price
	amount := housePrice*days+house.Deposit

	orderId, err := generateId(stub, OrderPrefix, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	var order Order
	order.ObjectType = OrderPrefix
	order.OrderId = orderId
	order.UserId = user.UserId
	order.HouseId = houseId
	order.CreateTime = toTimeZone(time.Now())
	order.BeginDate = start
	order.EndDate = end
	order.HousePrice = housePrice
	order.Days = days
	order.Amount = amount
	order.Comment = ""
	order.Status = OrderWaitAccept

	err = add(stub, OrderPrefix, orderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getOrderList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	role := args[1]

	var orderList []map[string]interface{}
	if role == "landlord" {
		houseIdList, err := getUserHouseIdList(stub, user.UserId)
		if err != nil {
			return shim.Error(err.Error())
		}

		for _, houseId := range houseIdList{
			orderIdList, err := getHouseOrderIdList(stub, houseId)
			if err != nil {
				continue
			}

			for _, orderId := range orderIdList{
				order, err := getOrderInfo(stub, orderId)
				if err != nil {
					continue
				}

				orderList = append(orderList, order)
			}
		}
	} else {
		orderIdList, err := getUserOrderIdList(stub, user.UserId)
		if err != nil {
			return shim.Error(err.Error())
		}

		for _, orderId := range orderIdList{
			order, err := getOrderInfo(stub, orderId)
			if err != nil {
				continue
			}

			orderList = append(orderList, order)
		}
	}

	data, err := json.Marshal(orderList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) handleOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	orderId := args[1]
	action := args[2]

	order, err := getOrder(stub, orderId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if order.Status != OrderWaitAccept {
		return shim.Error("forbidden")
	}

	house, err := getHouse(stub, order.HouseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if user.UserId != house.UserId {
		return shim.Error("forbidden")
	}

	if action == "reject" {
		order.Status = OrderRejected
	} else {
		house.OrderCount += 1
		err = update(stub, HousePrefix, order.HouseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}

		order.Status = OrderWaitComment
	}

	err = update(stub, OrderPrefix, orderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) comment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAddr := args[0]
	user, err := findUserByAddress(stub, userAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user.UserId == "" {
		return shim.Error("user does not exist")
	}

	orderId := args[1]
	comment := args[2]

	order, err := getOrder(stub, orderId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if order.Status != OrderWaitComment {
		return shim.Error("forbidden")
	}

	if user.UserId != order.UserId {
		return shim.Error("forbidden")
	}

	order.Comment = comment
	order.Status = OrderComplete

	err = update(stub, OrderPrefix, orderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) delArea(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"use_index": []string{"_design/areaDoc", "area"},
		"selector": map[string]interface{}{
			"ObjectType": AreaPrefix,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"AreaId": "asc"},
		},
		"fields": []string{"AreaId"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var area Area
		err = json.Unmarshal(item.Value, &area)
		if err != nil {
			continue
		}

		err = del(stub, AreaPrefix, area.AreaId)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getUserList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": UserPrefix,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var userList []User
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var user User
		err = json.Unmarshal(item.Value, &user)
		if err != nil {
			continue
		}

		userList = append(userList, user)
	}

	data, err := json.Marshal(userList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) delUser(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": UserPrefix,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var user User
		err = json.Unmarshal(item.Value, &user)
		if err != nil {
			continue
		}

		err = del(stub, UserPrefix, user.UserId)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentChaincode) delHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	houseId := args[0]
	err := del(stub, HousePrefix, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getHouseImageList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HouseImagePrefix,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var houseImageList []HouseImage
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var houseImage HouseImage
		err = json.Unmarshal(item.Value, &houseImage)
		if err != nil {
			continue
		}

		houseImageList = append(houseImageList, houseImage)
	}

	data, err := json.Marshal(houseImageList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *RentChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "addArea":
		return this.addArea(stub, args)
	case "getAreaList":
		return this.getAreaList(stub)
	case "addFacility":
		return this.addFacility(stub, args)
	case "getFacilityList":
		return this.getFacilityList(stub)
	case "addUser":
		return this.addUser(stub, args)
	case "updateUserAvatar":
		return this.updateUserAvatar(stub, args)
	case "rename":
		return this.rename(stub, args)
	case "auth":
		return this.auth(stub, args)
	case "getUser":
		return this.getUser(stub, args)
	case "addHouse":
		return this.addHouse(stub, args)
	case "updateHouseImage":
		return this.updateHouseImage(stub, args)
	case "getHouseDetail":
		return this.getHouseDetail(stub, args)
	case "getUserHouseList":
		return this.getUserHouseList(stub, args)
	case "searchHouse":
		return this.searchHouse(stub, args)
	case "addOrder":
		return this.addOrder(stub, args)
	case "getOrderList":
		return this.getOrderList(stub, args)
	case "handleOrder":
		return this.handleOrder(stub, args)
	case "comment":
		return this.comment(stub, args)
	case "delArea":
		return this.delArea(stub)
	case "getUserList":
		return this.getUserList(stub)
	case "delUser":
		return this.delUser(stub)
	case "delHouse":
		return this.delHouse(stub, args)
	case "getHouseImageList":
		return this.getHouseImageList(stub)
	default:
		return shim.Error("forbidden")
	}
}

func main()  {
	shim.Start(new(RentChaincode))
}
