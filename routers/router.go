package routers

import (
	"190720/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/area/list", &controllers.AreaController{}, "GET:GetList")

    beego.Router("/user/key", &controllers.UserController{}, "GET:GenerateKey")
    beego.Router("/register", &controllers.UserController{}, "POST:Register")
	beego.Router("/user/avatar", &controllers.UserController{}, "POST:UpdateAvatar")
	beego.Router("/user/rename", &controllers.UserController{}, "POST:Rename")
	beego.Router("/user/identify", &controllers.UserController{}, "POST:Identify")
	beego.Router("/user/info", &controllers.UserController{}, "POST:GetInfo")

    beego.Router("/house/add", &controllers.HouseController{}, "POST:Add")
    beego.Router("/house/image", &controllers.HouseController{}, "POST:UpdateImage")
    beego.Router("/house/list", &controllers.HouseController{}, "GET:GetList")
    beego.Router("/house/info", &controllers.HouseController{}, "GET:GetInfo")
    beego.Router("/house/index", &controllers.HouseController{}, "GET:Search")

    beego.Router("/order/add", &controllers.OrderController{}, "POST:Add")
	beego.Router("/order/list", &controllers.OrderController{}, "POST:GetList")
    beego.Router("/order/handle", &controllers.OrderController{}, "POST:Handle")
    beego.Router("/order/comment", &controllers.OrderController{}, "POST:Comment")
}
