package routers

import (
	"190305/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/area/add", &controllers.AreaController{}, "get:Add")
    beego.Router("/area/list", &controllers.AreaController{}, "get:Get")

	beego.Router("/facility/add", &controllers.FacilityController{}, "get:Add")
	beego.Router("/facility/list", &controllers.FacilityController{}, "get:Get")
}
