package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

type Response struct {
	Code int
	Msg string
	Data interface{}
}

func (this *BaseController) handleResponse(code int, msg string, data interface{}) {
	this.Data["json"] = &Response{Code: code, Msg: msg, Data: data}
	this.ServeJSON()
}
