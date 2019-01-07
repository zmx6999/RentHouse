package main

import (
        "github.com/micro/go-log"
        "github.com/micro/go-web"
        "190105/utils"
        "github.com/julienschmidt/httprouter"
        _ "190105/models"
        "net/http"
        "190105/web/handler"
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
        router.NotFound=http.FileServer(http.Dir("html"))

        router.GET("/api/v2.0/areas",handler.GetArea)

        router.GET("/api/v2.0/user/captcha/:id",handler.GetImageCpt)
        router.GET("/api/v2.0/user/sms/:id",handler.GetSmsCpt)
        router.POST("/api/v2.0/user/register",handler.Register)
        router.POST("/api/v2.0/user/login",handler.Login)
        router.GET("/api/v2.0/user/logout",handler.Logout)
        router.GET("/api/v2.0/user/info",handler.GetUserInfo)
        router.PUT("/api/v2.0/user/rename",handler.Rename)
        router.PUT("/api/v2.0/user/auth",handler.Auth)
        router.POST("/api/v2.0/user/avatar",handler.UploadAvatar)

        router.POST("/api/v2.0/houses/add",handler.AddHouse)
        router.GET("/api/v2.0/houses/list",handler.GetUserHouses)
        router.POST("/api/v2.0/houses/image/:id",handler.UploadHouseImage)
        router.GET("/api/v2.0/houses/detail/:id",handler.GetHouseDetail)
        router.GET("/api/v2.0/houses/banner",handler.GetIndexBanner)
        router.GET("/api/v2.0/houses/search",handler.SearchHouse)

        router.POST("/api/v2.0/order/add",handler.AddOrder)
        router.GET("/api/v2.0/order/list",handler.GetOrders)
        router.PUT("/api/v2.0/order/handle/:id",handler.HandleOrder)
        router.PUT("/api/v2.0/order/comment/:id",handler.CommentOrder)

	// register html handler
	    service.Handle("/", router)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
