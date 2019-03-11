package controllers

import (
	"github.com/astaxie/beego"
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
	this.Data["json"] = &ResponseJSON{Code: code, Msg: msg, Data: data}
	this.ServeJSON()
}

func (this *BaseController) handleSuccess(data interface{})  {
	this.handleResponse(200, "ok", data)
}

func (this *BaseController) handleError(code int, err error)  {
	this.handleResponse(code, err.Error(), nil)
}
