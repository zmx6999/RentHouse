package controllers

/*
0x522aca1Cbc2489C9CC87d76bBF30976777859BAC
0xcb08966f133ab88951d247ebfe24c14ad1968299134d298b55395e30f4d468c2
4984157656
0xe6d11A65CcC26702C4786a0Feae374D28dE7DDD4
0x7e8d5ad666d84175b802f4b6160d8f653aa183c3ee0a78091cfc3e1903155244
3918506041
0x182Da35736703daDac8c7F98dEc968CeDF7b4b0a
0x9d4244ebb366cc93a7d929d743622555fe219f14624b7eb2919274bb0d1e1294
 */

import (
	ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
			"190720/utils"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190720/models"
	"encoding/json"
)

type UserController struct {
	BaseController
}

func (this *UserController) GenerateKey()  {
	privateKey, err := ethereum_crypto.GenerateKey()
	if err != nil {
		this.error(1011, err.Error())
		return
	}

	privateKeyHex := utils.EncodePrivateKey(privateKey)
	address := utils.AddressFromPrivateKey(privateKey)

	data := map[string]interface{}{
		"private_key": privateKeyHex,
		"address": address,
	}
	this.success(data)
}

func (this *UserController) Register()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1003, err.Error())
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

	addr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	mobile := request["mobile"].(string)
	tx, err := ccs.ChaincodeUpdate("addUser", [][]byte{[]byte(mobile), []byte(addr)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *UserController) UpdateAvatar()  {
	request := map[string]interface{}{
		"private_key": this.Ctx.Request.Form.Get("private_key"),
	}

	addr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	avatarUrl, err := this.upload("avatar", []string{"jpg", "png", "jpeg"}, 1024*1024*2)
	if err != nil {
		this.error(1006, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	tx, err := ccs.ChaincodeUpdate("updateUserAvatar", [][]byte{[]byte(addr), []byte(avatarUrl)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *UserController) Rename()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1003, err.Error())
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

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err = gv.Validate()
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	addr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	newName := request["new_name"].(string)
	tx, err := ccs.ChaincodeUpdate("rename", [][]byte{[]byte(addr), []byte(newName)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *UserController) Identify()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1003, err.Error())
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

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err = gv.Validate()
	if err != nil {
		this.error(1004, err.Error())
		return
	}

	addr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	realName := request["real_name"].(string)
	idCard := request["id_card"].(string)
	tx, err := ccs.ChaincodeUpdate("identify", [][]byte{[]byte(addr), []byte(realName), []byte(idCard)})
	if err != nil {
		this.error(1002, err.Error())
		return
	}

	r := map[string]interface{}{
		"transaction_id": string(tx),
	}
	this.success(r)
}

func (this *UserController) GetInfo()  {
	request, err := this.postParam()
	if err != nil {
		this.error(1003, err.Error())
		return
	}

	addr, err := validateUser(request)
	if err != nil {
		this.error(1005, err.Error())
		return
	}

	ccs, err := models.Initialize(models.ChannelId, models.UserName, models.OrgName, models.ChaincodeId, models.ConfigFile)
	if err != nil {
		this.error(1001, err.Error())
		return
	}

	payload, err := ccs.ChaincodeQuery("getUserInfo", [][]byte{[]byte(addr)})
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
