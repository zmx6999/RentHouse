package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"math/rand"
	"time"
	"strings"
	"math"
)

const (
	AreaObjectType = "Area"
	FacilityObjectType = "Facility"
	UserObjectType = "User"
	HouseObjectType = "House"
	HouseImageObjectType = "HouseImage"
	OrderObjectType = "Order"

	TimeZone = "Asia/Saigon"
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"

	FastDFSHost = "149.28.210.102"
	FastDFSPort = 8888

	OrderStatusWaitAccept = "WAIT_ACCEPT"
	OrderStatusWaitComment = "WAIT_COMMENT"
	OrderStatusComplete = "COMPLETE"
	OrderStatusRejected = "REJECTED"
)

func put(stub shim.ChaincodeStubInterface, objectType string, id string, obj interface{}) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
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

func find(stub shim.ChaincodeStubInterface, objectType string, id string) ([]byte, error) {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return nil, err
	}

	data, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func add(stub shim.ChaincodeStubInterface, objectType string, id string, obj interface{}) error {
	data, err := find(stub, objectType, id)
	if err != nil {
		return err
	}
	if data != nil {
		return fmt.Errorf("key exist")
	}

	return put(stub, objectType, id, obj)
}

func update(stub shim.ChaincodeStubInterface, objectType string, id string, obj interface{}) error {
	data, err := find(stub, objectType, id)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("key does not exist")
	}

	return put(stub, objectType, id, obj)
}

func del(stub shim.ChaincodeStubInterface, objectType string, id string) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return err
	}

	err = stub.DelState(key)
	if err != nil {
		return err
	}

	return nil
}

func getQueryResult(stub shim.ChaincodeStubInterface, query map[string]interface{}) (shim.StateQueryIteratorInterface, error) {
	queryData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	return stub.GetQueryResult(string(queryData))
}

func generateId(stub shim.ChaincodeStubInterface, objectType string, length int) (string, error) {
	id := ""
	for i := 0; i < length; i++ {
		if i > 0 {
			id += strconv.Itoa(rand.Intn(10))
		} else {
			id += strconv.Itoa(rand.Intn(9)+1)
		}
	}

	data, err := find(stub, objectType, id)
	if err != nil {
		return "", err
	}
	if data != nil {
		return generateId(stub, objectType, length)
	}

	return id, nil
}

func checkArgs(args []string, n int) error {
	if len(args) < n {
		return fmt.Errorf("%d parameter(s) required", n)
	}

	return nil
}

func toTimeZone(t time.Time) time.Time {
	loc, _ := time.LoadLocation(TimeZone)
	return t.In(loc)
}

func today() time.Time {
	now := toTimeZone(time.Now()).Format(DateLayout)
	loc, _ := time.LoadLocation(TimeZone)
	r, _ := time.ParseInLocation(TimeLayout, now+" 00:00:00", loc)
	return r
}

func addDomain(url string) string {
	if url == "" {
		return ""
	}
	return "http://"+FastDFSHost+":"+strconv.Itoa(FastDFSPort)+"/"+url
}

type RentChaincode struct {

}

type Area struct {
	ObjectType string
	AreaId string
	AreaName string
}

func findArea(stub shim.ChaincodeStubInterface, id string) (*Area, error) {
	data, err := find(stub, AreaObjectType, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var area Area
	err = json.Unmarshal(data, &area)
	if err != nil {
		return nil, err
	}

	return &area, nil
}

func (this *RentChaincode) addArea(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	var area Area
	area.ObjectType = AreaObjectType
	area.AreaId = args[0]
	area.AreaName = args[1]

	err = add(stub, AreaObjectType, area.AreaId, area)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getAreaList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": AreaObjectType,
		},
		"use_index": []string{"_design/areaDoc", "area"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"AreaId": "asc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var areaList []map[string]interface{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var _area Area
		err = json.Unmarshal(item.Value, &_area)
		if err != nil {
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
	ObjectType string
	FacilityId string
	FacilityName string
}

func findFacilityMap(stub shim.ChaincodeStubInterface) (map[string]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": FacilityObjectType,
		},
		"use_index": []string{"_design/facilityDoc", "facility"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"FacilityId": "asc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	facilityMap := make(map[string]string)
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var facility Facility
		err = json.Unmarshal(item.Value, &facility)
		if err == nil {
			facilityMap[facility.FacilityId] = facility.FacilityName
		}
	}

	return facilityMap, nil
}

func (this *RentChaincode) addFacility(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	var facility Facility
	facility.ObjectType = FacilityObjectType
	facility.FacilityId = args[0]
	facility.FacilityName = args[1]

	err = add(stub, FacilityObjectType, facility.FacilityId, facility)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getFacilityList(stub shim.ChaincodeStubInterface) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": FacilityObjectType,
		},
		"use_index": []string{"_design/facilityDoc", "facility"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"FacilityId": "asc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var facilityList []map[string]interface{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var _facility Facility
		err = json.Unmarshal(item.Value, &_facility)
		if err != nil {
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

type User struct {
	ObjectType string
	UserId string
	UserName string
	Address string
	Mobile string
	RealName string
	IdCard string
	AvatarUrl string
	CreateTime time.Time
}

func findUserById(stub shim.ChaincodeStubInterface, id string) (*User, error) {
	data, err := find(stub, UserObjectType, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func findUserByMobile(stub shim.ChaincodeStubInterface, mobile string) (*User, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/user_mobileDoc", "user_mobile"},
		"selector": map[string]interface{}{
			"ObjectType": UserObjectType,
			"Mobile": mobile,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var user User
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		err = json.Unmarshal(item.Value, &user)
		if err == nil {
			return &user, nil
		}
	}

	return nil, nil
}

func findUserByAddress(stub shim.ChaincodeStubInterface, address string) (*User, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/user_addressDoc", "user_address"},
		"selector": map[string]interface{}{
			"ObjectType": UserObjectType,
			"Address": address,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var user User
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		err = json.Unmarshal(item.Value, &user)
		if err == nil {
			return &user, nil
		}
	}

	return nil, nil
}

func (this *RentChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	mobile := args[0]
	_user, err := findUserByMobile(stub, mobile)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _user != nil {
		return shim.Error("mobile exist")
	}

	address := args[1]
	_user, err = findUserByAddress(stub, address)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _user != nil {
		return shim.Error("address exist")
	}

	id, err := generateId(stub, UserObjectType, 10)
	if err != nil {
		return shim.Error(err.Error())
	}

	var user User
	user.ObjectType = UserObjectType
	user.UserId = id
	user.UserName = mobile
	user.Address = address
	user.Mobile = mobile
	user.CreateTime = toTimeZone(time.Now())

	err = add(stub, UserObjectType, id, user)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func updateUser(stub shim.ChaincodeStubInterface, address string, setUser func(user *User)) error {
	user, err := findUserByAddress(stub, address)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user does not exist")
	}

	setUser(user)

	err = update(stub, UserObjectType, user.UserId, user)
	if err != nil {
		return err
	}

	return nil
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

func (this *RentChaincode) identify(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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

func (this *RentChaincode) getUserInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	addr := args[0]

	user, err := findUserByAddress(stub, addr)
	if err != nil {
		return shim.Error(err.Error())
	}
	if user == nil {
		return shim.Error("user does not exist")
	}

	info := map[string]interface{}{
		"user_id": user.UserId,
		"username": user.UserName,
		"address": user.Address,
		"mobile": user.Mobile,
		"real_name": user.RealName,
		"id_card": user.IdCard,
		"avatar_url": addDomain(user.AvatarUrl),
		"create_time": user.CreateTime.Format(TimeLayout),
	}

	data, err := json.Marshal(info)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

type House struct {
	ObjectType string
	HouseId string
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
	FacilityIdList []string
}

func findHouse(stub shim.ChaincodeStubInterface, id string) (*House, error) {
	data, err := find(stub, HouseObjectType, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var house House
	err = json.Unmarshal(data, &house)
	if err != nil {
		return nil, err
	}

	return &house, nil
}

func getHouseInfo(stub shim.ChaincodeStubInterface, id string) (map[string]interface{}, error) {
	house, err := findHouse(stub, id)
	if err != nil {
		return nil, err
	}
	if house == nil {
		return nil, fmt.Errorf("house does not exist")
	}

	area, err := findArea(stub, house.AreaId)
	if err != nil || area == nil {
		area = &Area{}
	}

	user, err := findUserById(stub, house.UserId)
	if err != nil || user == nil {
		user = &User{}
	}

	return map[string]interface{}{
		"house_id": house.HouseId,
		"title": house.Title,
		"price": house.Price,
		"area_name": area.AreaName,
		"index_image_url": addDomain(house.IndexImageUrl),
		"room_count": house.RoomCount,
		"order_count": house.OrderCount,
		"address": house.Address,
		"user_avatar": addDomain(user.AvatarUrl),
		"create_time": house.CreateTime.Format(TimeLayout),
	}, nil
}

func getHouseDetail(stub shim.ChaincodeStubInterface, id string) (map[string]interface{}, error) {
	house, err := findHouse(stub, id)
	if err != nil {
		return nil, err
	}
	if house == nil {
		return nil, fmt.Errorf("house does not exist")
	}

	user, err := findUserById(stub, house.UserId)
	if err != nil || user == nil {
		user = &User{}
	}

	var imageUrlList []string
	_imageUrlList, err := findHouseImageUrlList(stub, id)
	if err == nil {
		for _, imageUrl := range _imageUrlList {
			imageUrlList = append(imageUrlList, addDomain(imageUrl))
		}
	}

	var facilities []string
	facilityMap, err := findFacilityMap(stub)
	if err == nil {
		for _, facilityId := range house.FacilityIdList{
			facilities = append(facilities, facilityMap[facilityId])
		}
	}

	var comments []map[string]interface{}
	orderIdList, err := findHouseOrderIdList(stub, id)
	if err == nil {
		for _, orderId := range orderIdList{
			order, err := findOrder(stub, orderId)
			if err != nil || order == nil {
				continue
			}
			if order.Status != OrderStatusComplete || order.Comment == "" {
				continue
			}

			user, err := findUserById(stub, order.UserId)
			if err != nil || user == nil {
				user = &User{}
			}

			comment := map[string]interface{}{
				"username": user.UserName,
				"comment": order.Comment,
				"create_time": order.CreateTime.Format(TimeLayout),
			}
			comments = append(comments, comment)
		}
	}

	return map[string]interface{}{
		"house_id": house.HouseId,
		"user_id": house.UserId,
		"username": user.UserName,
		"user_avatar": addDomain(user.AvatarUrl),
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
		"facilities": facilities,
		"comments": comments,
	}, nil
}

func findUserHouseIdList(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	user, err := findUserById(stub, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user does not exist")
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HouseObjectType,
			"UserId": user.UserId,
		},
		"use_index": []string{"_design/house_userDoc", "house_user"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var houseIdList []string
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var house House
		err = json.Unmarshal(item.Value, &house)
		if err == nil {
			houseIdList = append(houseIdList, house.HouseId)
		}
	}

	return houseIdList, nil
}

func (this *RentChaincode) addHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := findUserByAddress(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if user == nil {
		return shim.Error("user does not exist")
	}

	title := args[1]

	price, err := strconv.Atoi(args[2])
	if err != nil || price < 0 {
		price = 0
	}

	areaId := args[3]
	address := args[4]

	roomCount, err := strconv.Atoi(args[5])
	if err != nil || roomCount < 1 {
		roomCount = 1
	}

	acreage, err := strconv.Atoi(args[6])
	if err != nil || acreage < 0 {
		acreage = 0
	}

	unit := args[7]

	capacity, err := strconv.Atoi(args[8])
	if err != nil || capacity < 1 {
		capacity = 1
	}

	beds := args[9]

	deposit, err := strconv.Atoi(args[10])
	if err != nil || deposit < 0 {
		deposit = 0
	}

	minDays, err := strconv.Atoi(args[11])
	if err != nil || minDays < 1 {
		minDays = 1
	}

	maxDays, err := strconv.Atoi(args[12])
	if err != nil || maxDays < 0 {
		maxDays = 0
	}

	facility := strings.Split(args[13], ",")

	id, err := generateId(stub, HouseObjectType, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	var house House
	house.ObjectType = HouseObjectType
	house.HouseId = id
	house.UserId = user.UserId
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
	house.FacilityIdList = facility
	house.CreateTime = toTimeZone(time.Now())

	err = add(stub, HouseObjectType, id, house)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentChaincode) uploadHouseImage(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := findUserByAddress(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if user == nil {
		return shim.Error("user does not exist")
	}

	house, err := findHouse(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if house == nil {
		return shim.Error("house does not exist")
	}

	if house.UserId != user.UserId {
		return shim.Error("forbidden")
	}

	imageUrl := args[2]

	houseImageId, err := generateId(stub, HouseImageObjectType, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	var houseImage HouseImage
	houseImage.ObjectType = HouseImageObjectType
	houseImage.HouseImageId = houseImageId
	houseImage.HouseId = house.HouseId
	houseImage.Url = imageUrl

	err = add(stub, HouseImageObjectType, houseImage.HouseImageId, houseImage)
	if err != nil {
		return shim.Error(err.Error())
	}

	if house.IndexImageUrl == "" {
		house.IndexImageUrl = imageUrl
		err = update(stub, HouseObjectType, house.HouseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentChaincode) getHouseList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	userId := args[0]
	houseIdList, err := findUserHouseIdList(stub, userId)
	if err != nil {
		return shim.Error(err.Error())
	}

	var houseList []map[string]interface{}
	for _, houseId := range houseIdList{
		house, err := getHouseInfo(stub, houseId)
		if err == nil {
			houseList = append(houseList, house)
		}
	}

	data, err := json.Marshal(houseList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentChaincode) getHouseInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgs(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	houseId := args[0]
	_house, err := findHouse(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _house == nil {
		return shim.Error("house does not exist")
	}

	house, err := getHouseDetail(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := json.Marshal(house)
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

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HouseObjectType,
		},
		"use_index": []string{"_design/houseDoc", "house"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	selector := query["selector"].(map[string]interface{})

	areaId := args[0]
	if areaId != "" {
		query["use_index"] = []string{"_design/house_areaDoc", "house_area"}
		selector["AreaId"] = areaId
	}

	start := args[1]
	end := args[2]
	loc, _ := time.LoadLocation(TimeZone)
	if start != "" && end != "" {
		_start, _ := time.ParseInLocation(TimeLayout, start+" 00:00:00", loc)
		_end, _ := time.ParseInLocation(TimeLayout, end+" 23:59:59", loc)
		selector["$and"] = []map[string]interface{}{
			map[string]interface{}{
				"CreateTime": map[string]interface{}{"$gte": _start},
			},
			map[string]interface{}{
				"CreateTime": map[string]interface{}{"$lte": _end},
			},
		}
	}
	if start != "" && end == "" {
		_start, _ := time.ParseInLocation(TimeLayout, start+" 00:00:00", loc)
		selector["CreateTime"] = map[string]interface{}{"$gte": _start}
	}
	if start == "" && end != "" {
		_end, _ := time.ParseInLocation(TimeLayout, end+" 23:59:59", loc)
		selector["CreateTime"] = map[string]interface{}{"$lte": _end}
	}

	query["selector"] = selector

	page, err := strconv.Atoi(args[3])
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(args[4])
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	total := 0
	skip := pageSize*(page-1)
	var houseList []map[string]interface{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		total++
		if total <= skip || total > skip+pageSize {
			continue
		}

		var _house House
		err = json.Unmarshal(item.Value, &_house)
		if err == nil {
			houseId := _house.HouseId
			house, err := getHouseInfo(stub, houseId)
			if err == nil {
				houseList = append(houseList, house)
			}
		}
	}

	totalPage := int(math.Ceil(float64(total)/float64(pageSize)))

	r := map[string]interface{}{
		"houses": houseList,
		"total": total,
		"total_page": totalPage,
		"current_page": page,
	}

	data, err := json.Marshal(r)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

type HouseImage struct {
	ObjectType string
	HouseImageId string
	HouseId string
	Url string
}

func findHouseImageUrlList(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"use_index": []string{"_design/user_imageDoc", "user_image"},
		"selector": map[string]interface{}{
			"ObjectType": HouseImageObjectType,
			"HouseId": houseId,
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var urlList []string
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var houseImage HouseImage
		err = json.Unmarshal(item.Value, &houseImage)
		if err == nil {
			urlList = append(urlList, houseImage.Url)
		}
	}

	return urlList, nil
}

type Order struct {
	ObjectType string
	OrderId string
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

func findOrder(stub shim.ChaincodeStubInterface, id string) (*Order, error) {
	data, err := find(stub, OrderObjectType, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func findHouseOrderIdList(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": OrderObjectType,
			"HouseId": houseId,
		},
		"use_index": []string{"_design/order_houseDoc", "order_house"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var orderIdList []string
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var order Order
		err = json.Unmarshal(item.Value, &order)
		if err == nil {
			orderIdList = append(orderIdList, order.OrderId)
		}
	}

	return orderIdList, nil
}

func findUserOrderIdList(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": OrderObjectType,
			"UserId": userId,
		},
		"use_index": []string{"_design/order_userDoc", "order_user"},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var orderIdList []string
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var order Order
		err = json.Unmarshal(item.Value, &order)
		if err == nil {
			orderIdList = append(orderIdList, order.OrderId)
		}
	}

	return orderIdList, nil
}

func getOrderInfo(stub shim.ChaincodeStubInterface, id string) (map[string]interface{}, error) {
	order, err := findOrder(stub, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, nil
	}

	house, err := findHouse(stub, order.HouseId)
	if err != nil || house == nil {
		house = &House{}
	}

	return map[string]interface{}{
		"order_id": order.OrderId,
		"title": house.Title,
		"image_url": addDomain(house.IndexImageUrl),
		"begin_date": order.BeginDate.Format(DateLayout),
		"end_date": order.EndDate.Format(DateLayout),
		"create_time": order.CreateTime.Format(TimeLayout),
		"days": order.Days,
		"amount": order.Amount,
		"status": order.Status,
		"comment": order.Comment,
	}, nil
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
	if user == nil {
		return shim.Error("user does not exist")
	}

	houseId := args[1]
	house, err := findHouse(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if house == nil {
		return shim.Error("house does not exist")
	}

	if house.UserId == user.UserId {
		return shim.Error("forbidden")
	}

	start := args[2]
	end := args[3]

	loc, _ := time.LoadLocation(TimeZone)
	_start, _ := time.ParseInLocation(TimeLayout, start+" 00:00:00", loc)
	_end, _ := time.ParseInLocation(TimeLayout, end+" 00:00:00", loc)

	if _start.Sub(_end).Seconds() > 0 {
		return shim.Error("start_date should be no later than end_date")
	}
	if today().Sub(_start).Seconds() > 0 {
		return shim.Error("today should be no later than start_date")
	}

	days := int(math.Round(_end.Sub(_start).Seconds()/86400))+1
	if days < house.MinDays {
		return shim.Error(fmt.Sprintf("at least %d day(s)", house.MinDays))
	}
	if house.MaxDays > 0 && days > house.MaxDays {
		return shim.Error(fmt.Sprintf("at most %d day(s)", house.MaxDays))
	}

	orderId, err := generateId(stub, OrderObjectType, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	var order Order
	order.ObjectType = OrderObjectType
	order.OrderId = orderId
	order.UserId = user.UserId
	order.HouseId = house.HouseId
	order.BeginDate = _start
	order.EndDate = _end
	order.HousePrice = house.Price
	order.Days = days
	order.Amount = house.Price*days
	order.Status = OrderStatusWaitAccept
	order.CreateTime = toTimeZone(time.Now())

	err = add(stub, OrderObjectType, orderId, order)
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
	if user == nil {
		return shim.Error("user does not exist")
	}

	var orderList []map[string]interface{}

	role := args[1]
	if role == "landlord" {
		houseIdList, err := findUserHouseIdList(stub, user.UserId)
		if err == nil {
			for _, houseId := range houseIdList{
				orderIdList, err := findHouseOrderIdList(stub, houseId)
				if err == nil {
					for _, orderId := range orderIdList{
						order, err := getOrderInfo(stub, orderId)
						if err == nil {
							orderList = append(orderList, order)
						}
					}
				}
			}
		}
	} else {
		orderIdList, err := findUserOrderIdList(stub, user.UserId)
		if err == nil {
			for _, orderId := range orderIdList{
				order, err := getOrderInfo(stub, orderId)
				if err == nil {
					orderList = append(orderList, order)
				}
			}
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
	if user == nil {
		return shim.Error("user does not exist")
	}

	orderId := args[1]
	order, err := findOrder(stub, orderId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if order == nil {
		return shim.Error("order does not exist")
	}

	houseId := order.HouseId
	house, err := findHouse(stub, houseId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if house == nil {
		return shim.Error("house does not exist")
	}

	if house.UserId != user.UserId {
		return shim.Error("forbidden")
	}
	if order.Status != OrderStatusWaitAccept {
		return shim.Error("forbidden")
	}

	action := args[2]
	if action == "reject" {
		order.Status = OrderStatusRejected
	} else {
		order.Status = OrderStatusWaitComment

		house.OrderCount += 1
		err = update(stub, HouseObjectType, houseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	err = update(stub, OrderObjectType, orderId, order)
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
	if user == nil {
		return shim.Error("user does not exist")
	}

	orderId := args[1]
	order, err := findOrder(stub, orderId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if order == nil {
		return shim.Error("order does not exist")
	}

	if order.UserId != user.UserId {
		return shim.Error("forbidden")
	}
	if order.Status != OrderStatusWaitComment {
		return shim.Error("forbidden")
	}

	comment := args[2]
	order.Status = OrderStatusComplete
	order.Comment = comment
	err = update(stub, OrderObjectType, orderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
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
	case "identify":
		return this.identify(stub, args)
	case "getUserInfo":
		return this.getUserInfo(stub, args)
	case "addHouse":
		return this.addHouse(stub, args)
	case "uploadHouseImage":
		return this.uploadHouseImage(stub, args)
	case "getHouseList":
		return this.getHouseList(stub, args)
	case "getHouseInfo":
		return this.getHouseInfo(stub, args)
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
	default:
		return shim.Error("invalid")
	}
}

func main()  {
	shim.Start(new(RentChaincode))
}
