package routers

import (
	"190316/import_area_facility/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/area/add", &controllers.AreaController{}, "get:Add")
	beego.Router("/area/list", &controllers.AreaController{}, "get:List")

	beego.Router("/facility/add", &controllers.FacilityController{}, "get:Add")
	beego.Router("/facility/list", &controllers.FacilityController{}, "get:List")
}
