package main

import (
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
	"math/rand"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"fmt"
	"strings"
	"math"
)

const (
	USER_OBJECT_TYPE = "User"
	HOUSE_OBJECT_TYPE = "House"
	AREA_OBJECT_TYPE = "Area"
	FACILITY_OBJECT_TYPE = "Facility"
	HOUSE_IMAGE_OBJECT_TYPE = "HouseImage"
	ORDER_OBJECT_TYPE = "Order"

	TOTAL_KEY = "Total"

	TIMEZONE = "Asia/Saigon"
	TIME_FORMAT = "2006-01-02 15:04:05"
	DATE_FORMAT = "2006-01-02"

	FAST_DFS_HOST = "45.63.94.102"
	FAST_DFS_PORT = "8888"

	HOUSE_LIST_PAGE_SIZE = 2

	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"
)

func addDomain2Url(url string) string {
	if url == "" {
		return ""
	}

	return "http://" + FAST_DFS_HOST + ":" + FAST_DFS_PORT + "/" + url
}

func timeWithTimeZone(t time.Time) time.Time {
	lo, _ := time.LoadLocation(TIMEZONE)
	return t.In(lo)
}

func showTime(t time.Time) string {
	lo, _ := time.LoadLocation(TIMEZONE)
	return t.In(lo).Format(TIME_FORMAT)
}

func showDate(t time.Time) string {
	lo, _ := time.LoadLocation(TIMEZONE)
	return t.In(lo).Format(DATE_FORMAT)
}

func generateId(stub shim.ChaincodeStubInterface, objectType string, length int) (string, error) {
	x := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	id := ""
	for i := 0; i < length; i++ {
		if i > 0 {
			id += x[rand.Intn(len(x))]
		} else {
			id += x[rand.Intn(len(x) - 1)]
		}
	}

	exist, err := recordExist(stub, objectType, id)
	if err != nil {
		return "", err
	}

	if exist {
		return generateId(stub, objectType, length)
	}

	return id, nil
}

func recordExist(stub shim.ChaincodeStubInterface, objectType string, id string) (bool, error) {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return false, err
	}

	data, err := stub.GetState(key)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil
	}

	return true, nil
}

func get(stub shim.ChaincodeStubInterface, objectType string, id string) ([]byte, error) {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return nil, err
	}

	data, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("RECORD NOT FOUND")
	}

	return data, nil
}

func add(stub shim.ChaincodeStubInterface, objectType string, id string, obj interface{}, addTotal bool) error {
	exist, err := recordExist(stub, objectType, id)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(objectType + "Id already exists")
	}

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

	if !addTotal {
		return nil
	}

	totalKey, err := stub.CreateCompositeKey(objectType, []string{TOTAL_KEY})
	if err != nil {
		return err
	}

	totalData, err := stub.GetState(totalKey)
	if err != nil {
		return err
	}

	var total int
	if totalData != nil {
		total, err = strconv.Atoi(string(totalData))
		if err != nil {
			return err
		}
	}
	total++

	err = stub.PutState(totalKey, []byte(strconv.Itoa(total)))
	if err != nil {
		return err
	}

	return nil
}

func set(stub shim.ChaincodeStubInterface, objectType string, id string, obj interface{}) error {
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

func checkArgsNum(args []string, n int) error {
	if len(args) < n {
		return errors.New(fmt.Sprintf("%d argument(s) required", n))
	}

	return nil
}

func getQueryResult(stub shim.ChaincodeStubInterface, query map[string]interface{}) (shim.StateQueryIteratorInterface, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	return stub.GetQueryResult(string(queryJson))
}

type User struct {
	UserId string
	ObjectType string
	UserName string
	PublicKey string
	Mobile string
	RealName string
	IdCard string
	AvatarUrl string
}

func getUserById(stub shim.ChaincodeStubInterface, userId string) (*User, error) {
	data, err := get(stub, USER_OBJECT_TYPE, userId)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUserByMobile(stub shim.ChaincodeStubInterface, mobile string) (*User, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": USER_OBJECT_TYPE,
			"Mobile": mobile,
		},
		"use_index": []string{"_design/userMobileDoc", "userMobile"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var user User
		json.Unmarshal(item.Value, &user)
		if user.UserId != "" {
			return &user, nil
		}
	}

	return nil, errors.New("USER NOT FOUND")
}

func mobileExist(stub shim.ChaincodeStubInterface, mobile string) (bool, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": USER_OBJECT_TYPE,
			"Mobile": mobile,
		},
		"use_index": []string{"_design/userMobileDoc", "userMobile"},
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return false, err
	}
	defer iter.Close()

	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		var user User
		json.Unmarshal(item.Value, &user)
		if user.UserId != "" {
			return true, nil
		}
	}

	return false, nil
}

func setUserInfo(stub shim.ChaincodeStubInterface, mobile string, setUser func(*User)) error {
	user, err := getUserByMobile(stub, mobile)
	if err != nil {
		return err
	}

	setUser(user)

	err = set(stub, USER_OBJECT_TYPE, user.UserId, user)
	if err != nil {
		return err
	}

	return nil
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

	UserId string
	AreaId string

	FacilityIds []string
}

func getHouse(stub shim.ChaincodeStubInterface, houseId string) (*House, error) {
	data, err := get(stub, HOUSE_OBJECT_TYPE, houseId)
	if err != nil {
		return nil, err
	}

	var house House
	err = json.Unmarshal(data, &house)
	if err != nil {
		return nil, err
	}

	return &house, nil
}

func getHouseInfo(stub shim.ChaincodeStubInterface, houseId string) (map[string]interface{}, error) {
	house, err := getHouse(stub, houseId)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["house_id"] = house.HouseId
	data["title"] = house.Title
	data["price"] = house.Price
	data["index_image_url"] = addDomain2Url(house.IndexImageUrl)
	data["room_count"] = house.RoomCount
	data["order_count"] = house.OrderCount
	data["address"] = house.Address
	data["create_time"] = showTime(house.CreateTime)

	data["area_name"] = ""
	area, _ := getArea(stub, house.AreaId)
	if area != nil {
		data["area_name"] = area.AreaName
	}

	data["user_avatar"] = ""
	user, _ := getUserById(stub, house.UserId)
	if user != nil {
		data["user_avatar"] = addDomain2Url(user.AvatarUrl)
	}

	return data,nil
}

func getHouseDesc(stub shim.ChaincodeStubInterface, houseId string) (map[string]interface{}, error) {
	house, err := getHouse(stub, houseId)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["house_id"] = house.HouseId
	data["title"] = house.Title
	data["price"] = house.Price
	data["address"] = house.Address
	data["room_count"] = house.RoomCount
	data["acreage"] = house.Acreage
	data["unit"] = house.Unit
	data["capacity"] = house.Capacity
	data["beds"] = house.Beds
	data["deposit"] = house.Deposit
	data["min_days"] = house.MinDays
	data["max_days"] = house.MaxDays

	data["img_urls"] = []string{}
	urlList, _ := getHouseImageUrlList(stub, houseId)
	if len(urlList) > 0 {
		data["img_urls"] = urlList
	}

	data["user_id"] = ""
	data["username"] = ""
	data["user_avatar"] = ""
	user, _ := getUserById(stub, house.UserId)
	if user != nil {
		data["user_id"] = user.UserId
		data["username"] = user.RealName
		data["user_avatar"] = addDomain2Url(user.AvatarUrl)
	}

	var facilities []string
	for _, facilityId := range house.FacilityIds{
		facility, _ := getFacility(stub, facilityId)
		if facility != nil {
			facilities = append(facilities, facility.FacilityName)
		}
	}
	data["facilities"] = facilities

	var comments []map[string]interface{}
	orderIdList, _ := getOrderIdListByHouse(stub, houseId)
	for _, orderId := range orderIdList{
		order, _ := getOrder(stub, orderId)
		if order != nil && order.Comment != "" {
			comment := make(map[string]interface{})
			comment["comment"] = order.Comment
			comment["create_time"] = showTime(order.CreateTime)

			comment["username"] = ""
			orderUser, _ := getUserById(stub, order.UserId)
			if orderUser != nil {
				comment["username"] = orderUser.UserName
			}

			comments = append(comments, comment)
		}
	}
	data["comments"] = comments

	return data, nil
}

func getHouseIdListByUser(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HOUSE_OBJECT_TYPE,
			"UserId": userId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
		"use_index": []string{"_design/houseUserDoc", "houseUser"},
		"fields": []string{"HouseId"},
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
		json.Unmarshal(item.Value, &house)
		if house.HouseId != "" {
			houseIdList = append(houseIdList, house.HouseId)
		}
	}

	return houseIdList, nil
}

type Area struct {
	AreaId string
	ObjectType string
	AreaName string
}

func getArea(stub shim.ChaincodeStubInterface, areaId string) (*Area, error) {
	data, err := get(stub, AREA_OBJECT_TYPE, areaId)
	if err != nil {
		return nil, err
	}

	var area Area
	err = json.Unmarshal(data, &area)
	if err != nil {
		return nil, err
	}

	return &area, nil
}

type Facility struct {
	FacilityId string
	ObjectType string
	FacilityName string
}

func getFacility(stub shim.ChaincodeStubInterface, facilityId string) (*Facility, error) {
	data, err := get(stub, FACILITY_OBJECT_TYPE, facilityId)
	if err != nil {
		return nil, err
	}

	var facility Facility
	err = json.Unmarshal(data, &facility)
	if err != nil {
		return nil, err
	}

	return &facility, nil
}

type HouseImage struct {
	HouseImageId string
	ObjectType string
	HouseId string
	Url string
}

func getHouseImageUrlList(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HOUSE_IMAGE_OBJECT_TYPE,
			"HouseId": houseId,
		},
		"fields": []string{"Url"},
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
		json.Unmarshal(item.Value, &houseImage)
		if houseImage.Url != "" {
			urlList = append(urlList, addDomain2Url(houseImage.Url))
		}
	}

	return urlList, nil
}

type Order struct {
	OrderId    string
	ObjectType string
	BeginDate  time.Time
	EndDate    time.Time
	Days       int
	HousePrice int
	Amount     int
	Status     string
	Comment    string
	CreateTime time.Time
	Credit bool

	UserId string
	HouseId string
}

func getOrder(stub shim.ChaincodeStubInterface, orderId string) (*Order, error) {
	data, err := get(stub, ORDER_OBJECT_TYPE, orderId)
	if err != nil {
		return nil, err
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func getOrderIdListByHouse(stub shim.ChaincodeStubInterface, houseId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": ORDER_OBJECT_TYPE,
			"HouseId": houseId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
		"use_index": []string{"_design/orderHouseDoc", "orderHouse"},
		"fields": []string{"OrderId"},
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
		json.Unmarshal(item.Value, &order)
		if order.OrderId != "" {
			orderIdList = append(orderIdList, order.OrderId)
		}
	}

	return orderIdList, nil
}

func getOrderIdListByUser(stub shim.ChaincodeStubInterface, userId string) ([]string, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": ORDER_OBJECT_TYPE,
			"UserId": userId,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
		"use_index": []string{"_design/orderUserDoc", "orderUser"},
		"fields": []string{"OrderId"},
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
		json.Unmarshal(item.Value, &order)
		if order.OrderId != "" {
			orderIdList = append(orderIdList, order.OrderId)
		}
	}

	return orderIdList, nil
}

func getOrderInfo(stub shim.ChaincodeStubInterface, orderId string) (map[string]interface{}, error) {
	order, err := getOrder(stub, orderId)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["order_id"] = order.OrderId
	data["begin_date"] = showDate(order.BeginDate)
	data["end_date"] = showDate(order.EndDate)
	data["create_time"] = showTime(order.CreateTime)
	data["days"] = order.Days
	data["amount"] = order.Amount
	data["status"] = order.Status
	data["comment"] = order.Comment
	data["credit"] = order.Credit

	data["title"] = ""
	data["image_url"] = ""
	house, _ := getHouse(stub, order.HouseId)
	if house != nil {
		data["title"] = house.Title
		data["image_url"] = addDomain2Url(house.IndexImageUrl)
	}

	return data, nil
}

type RentingChaincode struct {

}

func (this *RentingChaincode) getAreaList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": AREA_OBJECT_TYPE,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"AreaId": "asc"},
		},
		"use_index": []string{"_design/areaDoc", "area"},
		"fields": []string{"AreaId", "AreaName"},
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
		json.Unmarshal(item.Value, &_area)

		if _area.AreaId != "" {
			area := make(map[string]interface{})
			area["area_id"] = _area.AreaId
			area["area_name"] = _area.AreaName
			areaList = append(areaList, area)
		}
	}

	data, err := json.Marshal(areaList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) generateUserId(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	id, err := generateId(stub, USER_OBJECT_TYPE, 11)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(id))
}

func (this *RentingChaincode) register(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	exist, err := mobileExist(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if exist {
		return shim.Error("Mobile has been registered")
	}

	var user User
	user.ObjectType = USER_OBJECT_TYPE
	user.UserId = args[0]
	user.Mobile = args[1]
	user.UserName = user.Mobile
	user.PublicKey = args[2]

	err = add(stub, USER_OBJECT_TYPE, user.UserId, user, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getUserInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	_data := make(map[string]interface{})
	_data["user_id"] = user.UserId
	_data["name"] = user.UserName
	_data["mobile"] = user.Mobile
	_data["real_name"] = user.RealName
	_data["id_card"] = user.IdCard
	_data["avatar_url"] = addDomain2Url(user.AvatarUrl)

	data, err := json.Marshal(_data)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) getUserPublicKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(user.PublicKey))
}

func (this *RentingChaincode) getUserId(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(user.UserId))
}

func (this *RentingChaincode) avatar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = setUserInfo(stub, args[0], func(user *User) {
		user.AvatarUrl = args[1]
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) rename(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = setUserInfo(stub, args[0], func(user *User) {
		user.UserName = args[1]
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) auth(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = setUserInfo(stub, args[0], func(user *User) {
		user.RealName = args[1]
		user.IdCard = args[2]
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) addHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	var house House
	house.ObjectType = HOUSE_OBJECT_TYPE

	house.HouseId, err = generateId(stub, HOUSE_OBJECT_TYPE, 14)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	house.UserId = user.UserId

	house.Title = args[1]

	house.Price, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	house.Price *= 100

	house.AreaId = args[3]

	house.Address = args[4]

	house.RoomCount, err = strconv.Atoi(args[5])
	if err != nil {
		return shim.Error(err.Error())
	}

	house.Acreage, err = strconv.Atoi(args[6])
	if err != nil {
		return shim.Error(err.Error())
	}

	house.Unit = args[7]

	house.Capacity, err = strconv.Atoi(args[8])
	if err != nil {
		return shim.Error(err.Error())
	}

	house.Beds = args[9]

	house.Deposit, err = strconv.Atoi(args[10])
	if err != nil {
		return shim.Error(err.Error())
	}
	house.Deposit *= 100

	house.MinDays, err = strconv.Atoi(args[11])
	if err != nil {
		return shim.Error(err.Error())
	}

	house.MaxDays, err = strconv.Atoi(args[12])
	if err != nil {
		return shim.Error(err.Error())
	}

	house.FacilityIds = strings.Split(args[13], ",")

	house.CreateTime = timeWithTimeZone(time.Now())

	err = add(stub, HOUSE_OBJECT_TYPE, house.HouseId, &house, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getHouseList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	houseIdList, err := getHouseIdListByUser(stub, user.UserId)
	if err != nil {
		return shim.Error(err.Error())
	}

	var houseList []map[string]interface{}
	for _, houseId := range houseIdList{
		house, _ := getHouseInfo(stub, houseId)
		if house != nil {
			houseList = append(houseList, house)
		}
	}

	data, err := json.Marshal(houseList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) uploadImage(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	house, err := getHouse(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	if house.UserId != user.UserId {
		return shim.Error("invalid user")
	}

	var houseImage HouseImage
	houseImage.ObjectType = HOUSE_IMAGE_OBJECT_TYPE

	houseImage.HouseImageId, err = generateId(stub, HOUSE_IMAGE_OBJECT_TYPE, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	houseImage.HouseId = args[1]

	houseImage.Url = args[2]

	err = add(stub, HOUSE_IMAGE_OBJECT_TYPE, houseImage.HouseImageId, &houseImage, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	if house.IndexImageUrl=="" {
		house.IndexImageUrl = args[2]
		err = set(stub, HOUSE_OBJECT_TYPE, house.HouseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getHouseDesc(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	house, err := getHouseDesc(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := json.Marshal(house)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) searchHouse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 4)
	if err != nil {
		return shim.Error(err.Error())
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": HOUSE_OBJECT_TYPE,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"CreateTime": "desc"},
		},
		"use_index": []string{"_design/houseDoc", "house"},
		"fields": []string{"HouseId"},
		// "limit": HOUSE_LIST_PAGE_SIZE,
	}

	if args[0] != "" {
		selector := query["selector"].(map[string]interface{})
		selector["AreaId"] = args[0]
		query["selector"] = selector
		query["use_index"] = []string{"_design/houseAreaDoc", "houseArea"}
	}

	if args[1] != "" && args[2] != "" {
		start, _ := time.Parse(TIME_FORMAT, args[1] + " 00:00:00")
		end, _ := time.Parse(TIME_FORMAT, args[2] + " 23:59:59")
		selector := query["selector"].(map[string]interface{})
		selector["$and"] = []map[string]interface{}{
			map[string]interface{}{
				"CreateTime": map[string]interface{}{"$gte": start},
			},
			map[string]interface{}{
				"CreateTime": map[string]interface{}{"$lte": end},
			},
		}
		query["selector"] = selector
	}
	if args[1] != "" && args[2] == "" {
		start, _ := time.Parse(TIME_FORMAT, args[1] + " 00:00:00")
		selector := query["selector"].(map[string]interface{})
		selector["CreateTime"] = map[string]interface{}{"$gte": start}
		query["selector"] = selector
	}
	if args[1] == "" && args[2] != "" {
		end, _ := time.Parse(TIME_FORMAT, args[2] + " 23:59:59")
		selector := query["selector"].(map[string]interface{})
		selector["CreateTime"] = map[string]interface{}{"$lte": end}
		query["selector"] = selector
	}

	iter, err := getQueryResult(stub, query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	page := 1
	if args[3] != "" {
		page, err = strconv.Atoi(args[3])
		if err != nil {
			page = 1
		}
	}
	skip := HOUSE_LIST_PAGE_SIZE * ( page - 1 )

	total := 0
	var houseList []map[string]interface{}
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}

		total++
		if total <= skip || total > skip + HOUSE_LIST_PAGE_SIZE {
			continue
		}

		var house House
		json.Unmarshal(item.Value, &house)
		if house.HouseId != "" {
			house, _ := getHouseInfo(stub, house.HouseId)
			if house != nil {
				houseList = append(houseList, house)
			}
		}
	}

	totalPage := int(math.Ceil(float64(total) / float64(HOUSE_LIST_PAGE_SIZE)))

	_data := make(map[string]interface{})
	_data["houses"] = houseList
	_data["total"] = total
	_data["total_page"] = totalPage
	_data["current_page"] = page

	data, err := json.Marshal(_data)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) addOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 4)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	house, err := getHouse(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	if house.UserId == user.UserId {
		return shim.Error("invalid user")
	}

	var order Order
	order.ObjectType = ORDER_OBJECT_TYPE

	order.OrderId, err = generateId(stub, ORDER_OBJECT_TYPE, 16)
	if err != nil {
		return shim.Error(err.Error())
	}

	order.UserId = user.UserId

	order.HouseId = args[1]

	order.CreateTime = timeWithTimeZone(time.Now())

	order.Status = ORDER_STATUS_WAIT_ACCEPT

	order.BeginDate, err = time.Parse(DATE_FORMAT, args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	order.EndDate, err = time.Parse(DATE_FORMAT, args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	order.Days = int(order.EndDate.Sub(order.BeginDate).Seconds() / 60 / 60 / 24) + 1

	order.HousePrice = house.Price

	order.Amount = order.HousePrice * order.Days

	err = add(stub, ORDER_OBJECT_TYPE, order.OrderId, &order, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getOrderList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	var orderList []map[string]interface{}
	if args[1] == "landlord" {
		houseIdList, err := getHouseIdListByUser(stub, user.UserId)
		if err != nil {
			return shim.Error(err.Error())
		}

		for _, houseId := range houseIdList{
			orderIdList, err := getOrderIdListByHouse(stub, houseId)
			if err != nil {
				continue
			}

			for _, orderId := range orderIdList{
				order , _ := getOrderInfo(stub, orderId)
				if order != nil {
					orderList = append(orderList, order)
				}
			}
		}
	} else {
		orderIdList, err := getOrderIdListByUser(stub, user.UserId)
		if err != nil {
			return shim.Error(err.Error())
		}

		for _, orderId := range orderIdList{
			order , _ := getOrderInfo(stub, orderId)
			if order != nil {
				orderList = append(orderList, order)
			}
		}
	}

	data, err := json.Marshal(orderList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) handleOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	order, err := getOrder(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	if order.Status != ORDER_STATUS_WAIT_ACCEPT {
		return shim.Error("invalid order status")
	}

	house, err := getHouse(stub, order.HouseId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if user.UserId != house.UserId {
		return shim.Error("invalid user")
	}

	if args[2] == "reject" {
		order.Status = ORDER_STATUS_REJECTED
	} else {
		order.Status = ORDER_STATUS_WAIT_COMMENT
	}

	err = set(stub, ORDER_OBJECT_TYPE, order.OrderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	if order.Status == ORDER_STATUS_WAIT_COMMENT {
		house.OrderCount++
		err = set(stub, HOUSE_OBJECT_TYPE, house.HouseId, house)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getOrderHouseId(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	order, err := getOrder(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(order.HouseId))
}

func (this *RentingChaincode) comment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 3)
	if err != nil {
		return shim.Error(err.Error())
	}

	user, err := getUserByMobile(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	order, err := getOrder(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	if user.UserId != order.UserId {
		return shim.Error("invalid user")
	}

	if order.Status != ORDER_STATUS_WAIT_COMMENT {
		return shim.Error("invalid order status")
	}

	order.Comment = args[2]
	order.Status = ORDER_STATUS_COMPLETE

	err = set(stub, ORDER_OBJECT_TYPE, order.OrderId, order)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) addArea(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	var area Area
	area.ObjectType = AREA_OBJECT_TYPE
	area.AreaId = args[0]
	area.AreaName = args[1]

	err = add(stub, AREA_OBJECT_TYPE, area.AreaId, &area, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) getFacilityList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": FACILITY_OBJECT_TYPE,
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"FacilityId": "asc"},
		},
		"use_index": []string{"_design/facilityDoc", "facility"},
		"fields": []string{"FacilityId", "FacilityName"},
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
		json.Unmarshal(item.Value, &_facility)

		if _facility.FacilityId != "" {
			facility := make(map[string]interface{})
			facility["facility_id"] = _facility.FacilityId
			facility["facility_name"] = _facility.FacilityName
			facilityList = append(facilityList, facility)
		}
	}

	data, err := json.Marshal(facilityList)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) addFacility(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	var facility Facility
	facility.ObjectType = FACILITY_OBJECT_TYPE
	facility.FacilityId = args[0]
	facility.FacilityName = args[1]

	err = add(stub, FACILITY_OBJECT_TYPE, facility.FacilityId, &facility, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) clear(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"ObjectType": "",
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

		err = stub.DelState(item.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	data ,err := get(stub, args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (this *RentingChaincode) setMobile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := checkArgsNum(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = setUserInfo(stub, args[0], func(user *User) {
		user.Mobile = args[1]
	})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (this *RentingChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *RentingChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "getAreaList":
		return this.getAreaList(stub, args)
	case "generateUserId":
		return this.generateUserId(stub, args)
	case "register":
		return this.register(stub, args)
	case "getUserInfo":
		return this.getUserInfo(stub, args)
	case "getUserPublicKey":
		return this.getUserPublicKey(stub, args)
	case "getUserId":
		return this.getUserId(stub, args)
	case "avatar":
		return this.avatar(stub, args)
	case "rename":
		return this.rename(stub, args)
	case "auth":
		return this.auth(stub, args)
	case "addHouse":
		return this.addHouse(stub, args)
	case "getHouseList":
		return this.getHouseList(stub, args)
	case "uploadImage":
		return this.uploadImage(stub, args)
	case "getHouseDesc":
		return this.getHouseDesc(stub, args)
	case "searchHouse":
		return this.searchHouse(stub, args)
	case "addOrder":
		return this.addOrder(stub, args)
	case "getOrderList":
		return this.getOrderList(stub, args)
	case "handleOrder":
		return this.handleOrder(stub, args)
	case "getOrderHouseId":
		return this.getOrderHouseId(stub, args)
	case "comment":
		return this.comment(stub, args)
	case "addArea":
		return this.addArea(stub, args)
	case "getFacilityList":
		return this.getFacilityList(stub, args)
	case "addFacility":
		return this.addFacility(stub, args)
	case "get":
		return this.get(stub, args)
	case "clear":
		return this.clear(stub, args)
	case "setMobile":
		return this.setMobile(stub, args)
	default:
		return shim.Error("invalid method")
	}
}

func main()  {
	shim.Start(new(RentingChaincode))
}
