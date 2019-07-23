package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/zmx6999/FormValidation/FormValidation"
	"190720/utils"
	"path"
	"strings"
	"fmt"
	"github.com/weilaihui/fdfs_client"
)

type ResponseJSON struct {
	Code int
	Msg string
	Data interface{}
}

type BaseController struct {
	beego.Controller
}

func (this *BaseController) handleResponse(code int, msg string, data interface{})  {
	this.Data["json"] = &ResponseJSON{code, msg, data}
	this.ServeJSON()
}

func (this *BaseController) success(data interface{})  {
	this.handleResponse(200, "ok", data)
}

func (this *BaseController) error(code int, msg string)  {
	this.handleResponse(code, msg, nil)
}

func (this *BaseController) postParam() (map[string]interface{}, error) {
	request := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func validateUser(request map[string]interface{}) (string, error) {
	fvs:=[]*FormValidation.FieldValidation{
		&FormValidation.FieldValidation{
			FieldName:       "private_key",
			ValidMethodName: "Require",
			ValidMethodArgs: []interface{}{},
			ErrMsg:          "private_key cannot be empty",
			Trim:            true,
			ValidEmpty:      true,
		},
	}

	gv := &FormValidation.GroupValidation{
		request,
		fvs,
	}
	_, err := gv.Validate()
	if err != nil {
		return "", err
	}

	privateKeyHex := request["private_key"].(string)
	privateKey, err := utils.DecodePrivateKey(privateKeyHex)
	if err != nil {
		return "", err
	}

	address := utils.AddressFromPrivateKey(privateKey)
	return address, nil
}

func (this *BaseController) upload(key string, allowTypeList []string, allowMaxSize int) (string, error) {
	file, head, err := this.GetFile(key)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := path.Ext(head.Filename)
	ext = strings.ToLower(ext)
	ext = strings.TrimPrefix(ext, ".")
	if utils.Find(allowTypeList, ext) < 0 {
		return "", fmt.Errorf("type should be %s", strings.Join(allowTypeList, ","))
	}
	if int(head.Size) > allowMaxSize {
		return "", fmt.Errorf("size exceed")
	}

	data := make([]byte, head.Size)
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}

	client, err := fdfs_client.NewFdfsClient(beego.AppConfig.String("fdfs_client_config"))
	if err != nil {
		return "", err
	}

	r, err := client.UploadByBuffer(data, ext)
	if err != nil {
		return "", err
	}

	return r.RemoteFileId, nil
}
