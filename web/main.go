package main

import (
        "github.com/micro/go-log"
	        "github.com/micro/go-web"
                "github.com/julienschmidt/httprouter"
	"190222/utils"
	"190222/web/handler"
	)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.web"),
                web.Version("latest"),
                web.Address(":"+utils.APPPort),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	router:=httprouter.New()

	router.GET("/area/list",handler.GetArea)

	router.GET("/user/id",handler.GenerateId)
	router.GET("/user/captcha/:user_id",handler.Captcha)
	router.GET("/user/sms",handler.SmsCaptcha)
	router.POST("/user/register",handler.Register)
	router.POST("/user/info",handler.GetUserInfo)
	router.POST("/user/avatar",handler.Avatar)
	router.POST("/user/rename",handler.Rename)
	router.POST("/user/auth",handler.Auth)

	router.POST("/house/add",handler.AddHouse)
	router.POST("/house/list",handler.GetHouseList)
	router.GET("/house/desc/:house_id",handler.GetHouseDesc)
	router.POST("/house/image/:house_id",handler.UploadHouseImage)
	router.GET("/house/index",handler.GetIndexHouseList)

	router.POST("/order/add",handler.AddOrder)
	router.POST("/order/list",handler.GetOrderList)
	router.POST("/order/handle",handler.HandleOrder)
	router.POST("/order/comment",handler.Comment)

	// register call handler
	service.Handle("/", router)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
