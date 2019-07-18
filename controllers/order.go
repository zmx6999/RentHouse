package controllers

import (
	"github.com/zmx6999/FormValidation/FormValidation"
	"190702/models"
	"encoding/json"
)

type OrderController struct {
	BaseController
}

func (this *OrderController) Add()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
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
		&FormValidation.FieldValidation{
			FieldName:"start",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"start cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"end cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"start",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"start invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end",
			ValidMethodName:"Date",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"end invalid",
			Trim:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"end",
			ValidMethodName:"EndDate",
			ValidMethodArgs:[]interface{}{request["start"]},
			ErrMsg:"end should be later than start",
			Trim:true,
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

	houseId := request["house_id"].(string)
	start := request["start"].(string)
	end := request["end"].(string)
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	tx, err := ccs.ChaincodeUpdate("addOrder", [][]byte{[]byte(userAddr), []byte(houseId), []byte(start), []byte(end)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *OrderController) GetList()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	userAddr, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	role := this.GetString("role")
	data, err := ccs.ChaincodeQuery("getOrderList", [][]byte{[]byte(userAddr), []byte(role)})
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

func (this *OrderController) Handle()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"order_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"order_id cannot be empty",
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

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	orderId := request["order_id"].(string)
	action := this.GetString("action")
	tx, err := ccs.ChaincodeUpdate("handleOrder", [][]byte{[]byte(userAddr), []byte(orderId), []byte(action)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *OrderController) Comment()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"order_id",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"order_id cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"comment",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"comment cannot be empty",
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

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	orderId := request["order_id"].(string)
	comment := request["comment"].(string)
	tx, err := ccs.ChaincodeUpdate("comment", [][]byte{[]byte(userAddr), []byte(orderId), []byte(comment)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}
