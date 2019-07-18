package controllers

import (
	ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
	"190702/utils"
		"190702/models"
	"encoding/json"
	"github.com/zmx6999/FormValidation/FormValidation"
)

/*
4095268767
0xB068a830140F31153f5C6751B72349ca31Ef0A92
0xa895362f045e7d54c5dadf86d680a8c463e40d120506f5d23fa44e96a30a9406
3029617152
0xfF52e5677f38f2e31E3079Ce78e7d46131A261eC
0x1b0abc01b72131cab3e0c6c6623cb9b116762b7a181879bc24c5d79363b44a73
0x68fC96D98C79DC00f3498D2cb589a362bBCD046B
0x035053489f532c57652fd0f009342fa29151d0365ae15a2e395672845f470a39
 */

type UserController struct {
	BaseController
}

func (this *UserController) GenerateKey()  {
	_privateKey, err := ethereum_crypto.GenerateKey()
	if err != nil {
		this.error(1011, err.Error())
		return
	}

	privateKey := utils.EncodePrivateKey(_privateKey)
	address := utils.AddressFromPrivateKey(_privateKey)
	data := map[string]interface{}{
		"privateKey": privateKey,
		"address": address,
	}
	this.success(data)
}

func (this *UserController) getInfo(address string) (map[string]interface{}, error) {
	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		return nil, err
	}

	data, err := ccs.ChaincodeQuery("getUser", [][]byte{[]byte(address)})
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (this *UserController) Register()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"mobile cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"mobile",
			ValidMethodName:"ChineseMobile",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"invalid mobile",
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

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	mobile := request["mobile"].(string)
	tx, err := ccs.ChaincodeUpdate("addUser", [][]byte{[]byte(mobile), []byte(address)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *UserController) GetInfo()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	data, err := this.getInfo(address)
	if err != nil {
		this.error(1031, err.Error())
		return
	}

	this.success(data)
}

func (this *UserController) UpdateAvatar()  {
	privateKey := this.Ctx.Request.Form.Get("private_key")
	request := map[string]interface{}{
		"private_key": privateKey,
	}

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	avatarUrl, err := this.upload("avatar", []string{"jpg", "png", "jpeg"}, 1024*1024*2)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	tx, err := ccs.ChaincodeUpdate("updateUserAvatar", [][]byte{[]byte(address), []byte(avatarUrl)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *UserController) Rename()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"new_name",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"new_name cannot be empty",
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

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	newName := request["new_name"].(string)
	tx, err := ccs.ChaincodeUpdate("rename", [][]byte{[]byte(address), []byte(newName)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}

func (this *UserController) Auth()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"real_name",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"real_name cannot be empty",
			Trim:true,
			ValidEmpty:true,
		},
		&FormValidation.FieldValidation{
			FieldName:"id_card",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"id_card cannot be empty",
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

	address, err := this.validate(request)
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1021, err.Error())
		return
	}

	realName := request["real_name"].(string)
	idCard := request["id_card"].(string)
	tx, err := ccs.ChaincodeUpdate("auth", [][]byte{[]byte(address), []byte(realName), []byte(idCard)})
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	data := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(data)
}
