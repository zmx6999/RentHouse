package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
		"fmt"
	"190702/utils"
	"strings"
	"path"
	"github.com/weilaihui/fdfs_client"
	"github.com/zmx6999/FormValidation/FormValidation"
)

type BaseController struct {
	beego.Controller
}

type ResponseJSON struct {
	Code int
	Msg string
	Data interface{}
}

func (this *BaseController) handleResponse(code int, msg string, data interface{})  {
	this.Data["json"] = &ResponseJSON{Code:code, Msg:msg, Data:data}
	this.ServeJSON()
}

func (this *BaseController) success(data interface{})  {
	this.handleResponse(200, "OK", data)
}

func (this *BaseController) error(code int, msg string)  {
	this.handleResponse(code, msg, nil)
}

func (this *BaseController) postParam() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *BaseController) upload(key string, allowType []string, allowMaxSize int) (string, error) {
	file, head, err := this.GetFile(key)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.TrimPrefix(path.Ext(head.Filename), ".")
	if utils.Find(allowType, strings.ToLower(ext)) < 0 {
		return "", fmt.Errorf("type should be %s", strings.Join(allowType, ","))
	}
	if int(head.Size) > allowMaxSize {
		return "", fmt.Errorf("size exceed")
	}

	data := make([]byte, head.Size)
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}

	configFile := beego.AppConfig.String("fdfs_config_file")
	client, err := fdfs_client.NewFdfsClient(configFile)
	if err != nil {
		return "", err
	}

	r, err := client.UploadByBuffer(data, ext)
	if err != nil {
		return "", err
	}

	return r.RemoteFileId, nil
}

func (this *BaseController) validate(request map[string]interface{}) (string, error) {
	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:"private_key",
			ValidMethodName:"Require",
			ValidMethodArgs:[]interface{}{},
			ErrMsg:"private_key cannot be empty",
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
		return "", err
	}

	privateKey, err := utils.DecodePrivateKey(request["private_key"].(string))
	if err != nil {
		this.error(1003, err.Error())
		return "", err
	}

	address := utils.AddressFromPrivateKey(privateKey)
	return address, nil
}
