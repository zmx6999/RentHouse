package main

import (
        "github.com/micro/go-log"
	        "github.com/micro/go-web"
                "190120/utils"
        "github.com/julienschmidt/httprouter"
        _ "190120/models"
        "190120/web/handler"
        "net/http"
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

        r:=httprouter.New()
        r.NotFound=http.FileServer(http.Dir("html"))

        r.GET("/api/v3.0/areas",handler.GetArea)

        r.GET( "/api/v3.0/user/captcha/:uuid",handler.GetCaptcha)
        r.GET("/api/v3.0/user/sms",handler.GetSmsCaptcha)
        r.POST("/api/v3.0/user/register",handler.Register)
        r.POST("/api/v3.0/user/login",handler.Login)
        r.GET("/api/v3.0/user/logout",handler.Logout)
        r.GET("/api/v3.0/user/info",handler.Info)
        r.POST("/api/v3.0/user/avatar",handler.Avatar)
        r.PUT("/api/v3.0/user/name",handler.UpdateUserName)
        r.PUT("/api/v3.0/user/auth",handler.Auth)

        r.POST("/api/v3.0/house/add",handler.AddHouse)
        r.GET("/api/v3.0/house/list",handler.GetHouses)
        r.POST("/api/v3.0/house/image/:house_id",handler.UploadHouseImage)
        r.GET("/api/v3.0/house/detail/:house_id",handler.GetHouseDetail)
        r.GET("/api/v3.0/house/banner",handler.GetIndexBanner)
        r.GET("/api/v3.0/house/search",handler.Search)

        r.POST("/api/v3.0/order/add",handler.AddOrder)
        r.GET("/api/v3.0/order/list",handler.GetOrders)
        r.GET("/api/v3.0/order/handle/:order_id",handler.HandleOrder)
        r.POST("/api/v3.0/order/comment/:order_id",handler.Comment)

        service.Handle("/",r)

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/example/call", handler.ExampleCall)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
