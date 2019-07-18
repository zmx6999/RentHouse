package controllers

import (
	"github.com/zmx6999/FormValidation/FormValidation"
	"190702/utils"
	"190702/models"
	"encoding/json"
)

type HouseController struct {
	BaseController
}

func (this *HouseController) Add()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
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
			FieldName:"address",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"address cannot be empty",
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
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err = gv.Validate()
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	userAddr, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	title := request["title"].(string)
	price := utils.GetStringValue(request, "price", "0")
	address := request["address"].(string)
	areaId := request["area_id"].(string)
	roomCount := utils.GetStringValue(request, "room_count", "1")
	acreage := utils.GetStringValue(request, "acreage", "0")
	unit := utils.GetStringValue(request, "unit", "")
	capacity := utils.GetStringValue(request, "capacity", "1")
	beds := utils.GetStringValue(request, "beds", "")
	deposit := utils.GetStringValue(request, "deposit", "0")
	minDays := utils.GetStringValue(request, "min_days", "1")
	maxDays := utils.GetStringValue(request, "max_days", "0")
	facility := utils.GetStringValue(request, "facilities", "")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
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
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *HouseController) UploadImage()  {
	privateKey := this.Ctx.Request.Form.Get("private_key")
	houseId := this.Ctx.Request.Form.Get("house_id")
	request := map[string]interface{}{
		"private_key": privateKey,
		"house_id": houseId,
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

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	url, err := this.upload("image", []string{"jpg", "png", "jpeg"}, 1024*1024*2)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	tx, err := ccs.ChaincodeUpdate("updateHouseImage", [][]byte{[]byte(address), []byte(houseId), []byte(url)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *HouseController) GetList()  {
	userId := this.GetString("user_id")
	request := map[string]interface{}{
		"user_id": userId,
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"user_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"user_id cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
	}

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	data, err := ccs.ChaincodeQuery("getUserHouseList", [][]byte{[]byte(userId)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	r := []map[string]interface{}{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(r)
}

func (this *HouseController) GetDetail()  {
	houseId := this.GetString("house_id")
	request := map[string]interface{}{
		"house_id": houseId,
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

	gv:=&FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	data, err := ccs.ChaincodeQuery("getHouseDetail", [][]byte{[]byte(houseId)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(data, &r)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(r)
}

func (this *HouseController) Search()  {
	areaId := this.GetString("area_id")
	start := this.GetString("start")
	end := this.GetString("end")
	page := this.GetString("page")
	pageSize := this.GetString("page_size")

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	data, err := ccs.ChaincodeQuery("searchHouse", [][]byte{[]byte(areaId), []byte(start), []byte(end), []byte(page), []byte(pageSize)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(data, &r)
	if err != nil {
		this.error(1007, err.Error())
		return
	}

	this.success(r)
}
