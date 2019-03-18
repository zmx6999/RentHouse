package main

import (
        "github.com/micro/go-log"
	        "github.com/micro/go-web"
                "190316/utils"
        "github.com/julienschmidt/httprouter"
	"190316/web/handler"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.web"),
                web.Version("latest"),
                web.Address(":" + utils.HttpPort),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	r := httprouter.New()

	r.GET("/area/list", handler.GetAreaList)

	r.GET("/user/id", handler.GenerateUserId)
	r.GET("/user/captcha/:user_id", handler.Captcha)
	r.GET("/user/sms", handler.SmsCaptcha)
	r.POST("/user/register", handler.Register)
	r.POST("/user/info", handler.GetUserInfo)
	r.POST("/user/avatar", handler.Avatar)
	r.POST("/user/rename", handler.Rename)
	r.POST("/user/auth", handler.Auth)

	r.POST("/house/add", handler.AddHouse)
	r.POST("/house/list", handler.GetLandlordHouseList)
	r.POST("/house/image", handler.UploadHouseImage)
	r.GET("/house/desc/:house_id", handler.GetHouseDesc)
	r.GET("/house/search", handler.SearchHouse)

	r.POST("/order/add", handler.AddOrder)
	r.POST("/order/list", handler.GetOrderList)
	r.POST("/order/handle", handler.HandleOrder)
	r.POST("/order/comment", handler.Comment)

	// register call handler
	service.Handle("/", r)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
