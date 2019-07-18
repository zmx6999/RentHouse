package routers

import (
		"github.com/astaxie/beego"
	"190702/controllers"
)

func init() {
	beego.Router("/area/list", &controllers.AreaController{}, "GET:GetList")

    beego.Router("/user/key", &controllers.UserController{}, "GET:GenerateKey")
	beego.Router("/register", &controllers.UserController{}, "POST:Register")
	beego.Router("/user/info", &controllers.UserController{}, "POST:GetInfo")
	beego.Router("/user/avatar", &controllers.UserController{}, "POST:UpdateAvatar")
	beego.Router("/user/rename", &controllers.UserController{}, "POST:Rename")
	beego.Router("/user/auth", &controllers.UserController{}, "POST:Auth")

	beego.Router("/house/add", &controllers.HouseController{}, "POST:Add")
	beego.Router("/house/image", &controllers.HouseController{}, "POST:UploadImage")
	beego.Router("/house/list", &controllers.HouseController{}, "GET:GetList")
	beego.Router("/house/detail", &controllers.HouseController{}, "GET:GetDetail")
	beego.Router("/house/index", &controllers.HouseController{}, "GET:Search")

	beego.Router("/order/add", &controllers.OrderController{}, "POST:Add")
	beego.Router("/order/list", &controllers.OrderController{}, "POST:GetList")
	beego.Router("/order/handle", &controllers.OrderController{}, "POST:Handle")
	beego.Router("/order/comment", &controllers.OrderController{}, "POST:Comment")
}
