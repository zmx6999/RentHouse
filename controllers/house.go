package controllers

import (
	"github.com/zmx6999/FormValidation/FormValidation"
	"190720/models"
	"190720/utils"
	"encoding/json"
)

type HouseController struct {
	BaseController
}

func (this *HouseController) Add()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"title",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"title cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"area_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"area_id cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"address",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"address cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err = gv.Validate()
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	userAddr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	title := request["title"].(string)
	price := utils.StringValue(request, "price", "")
	areaId := request["area_id"].(string)
	address := request["address"].(string)
	roomCount := utils.StringValue(request, "room_count", "")
	acreage := utils.StringValue(request, "acreage", "")
	unit := utils.StringValue(request, "unit", "")
	capacity := utils.StringValue(request, "capacity", "")
	beds := utils.StringValue(request, "beds", "")
	deposit := utils.StringValue(request, "deposit", "")
	minDays := utils.StringValue(request, "min_days", "")
	maxDays := utils.StringValue(request, "max_days", "")
	facility := utils.StringValue(request, "facility", "")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	tx, err := ccs.ChaincodeUpdate("addHouse", [][]byte{
		[]byte(userAddr),
		[]byte(title),
		[]byte(price),
		[]byte(areaId),
		[]byte(address),
		[]byte(roomCount),
		[]byte(acreage),
		[]byte(unit),
		[]byte(capacity),
		[]byte(beds),
		[]byte(deposit),
		[]byte(minDays),
		[]byte(maxDays),
		[]byte(facility),
	})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *HouseController) UpdateImage()  {
	request := map[string]interface{}{
		"private_key": this.Ctx.Request.Form.Get("private_key"),
		"house_id": this.Ctx.Request.Form.Get("house_id"),
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"house_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"house_id cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	userAddr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	imageUrl, err := this.upload("image", []string{"jpg", "png", "jpeg"}, 1024*1024*2)
	if err != nil {
		this.error(1006, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	houseId := request["house_id"].(string)
	tx, err := ccs.ChaincodeUpdate("uploadHouseImage", [][]byte{[]byte(userAddr), []byte(houseId), []byte(imageUrl)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *HouseController) GetList()  {
	userId := this.GetString("user_id")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	payload, err := ccs.ChaincodeQuery("getHouseList", [][]byte{[]byte(userId)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal(payload, &data)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(data)
}

func (this *HouseController) GetInfo()  {
	houseId := this.GetString("house_id")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	payload, err := ccs.ChaincodeQuery("getHouseInfo", [][]byte{[]byte(houseId)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(payload, &data)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(data)
}

func (this *HouseController) Search()  {
	areaId := this.GetString("area_id")
	start := this.GetString("start")
	end := this.GetString("end")
	page := this.GetString("page")
	pageSize := this.GetString("page_size")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	payload, err := ccs.ChaincodeQuery("searchHouse", [][]byte{[]byte(areaId), []byte(start), []byte(end), []byte(page), []byte(pageSize)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(payload, &data)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(data)
}
