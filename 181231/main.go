package main

import (
        "github.com/micro/go-log"
        "github.com/micro/go-web"
        "github.com/julienschmidt/httprouter"
        "net/http"
        _ "sss/181231/models"
        "sss/181231/handler"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.181231"),
                web.Version("latest"),
                web.Address(":8099"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

        router:=httprouter.New()
        router.NotFound = http.FileServer(http.Dir("html"))

        router.GET("/api/v1.0/areas",handler.GetArea)
        router.GET("/api/v1.0/captcha/:uuid",handler.GetImageCd)
        router.GET("/api/v1.0/sms/:mobile",handler.GetSmsCd)
        router.POST("/api/v1.0/register",handler.Register)
        router.GET("/api/v1.0/session",handler.GetSession)
        router.POST("/api/v1.0/login",handler.Login)
        router.GET("/api/v1.0/logout",handler.Logout)
        router.GET("/api/v1.0/user",handler.GetUserInfo)
        router.POST("/api/v1.0/avatar",handler.PostAvatar)
        router.PUT("/api/v1.0/user/rename",handler.PutUserInfo)
        router.PUT("/api/v1.0/user/auth",handler.PostUserAuth)
        router.GET("/api/v1.0/houses",handler.GetUserHouses)
        router.POST("/api/v1.0/houses",handler.PostHouse)
        router.POST("/api/v1.0/houses/images/:id",handler.PostHouseImage)
        router.GET("/api/v1.0/houses/detail/:id",handler.GetHouseInfo)
        router.GET("/api/v1.0/banner",handler.GetIndexBanner)
        router.GET("/api/v1.0/houses/search",handler.GetHouses)
        router.POST("/api/v1.0/orders",handler.PostOrder)
        router.GET("/api/v1.0/orders",handler.GetUserOrder)
        router.PUT("/api/v1.0/orders/:id",handler.PutOrder)
        router.PUT("/api/v1.0/comment/:id",handler.PutComment)

        service.Handle("/",router)
        if err:=service.Run();err != nil {
                log.Fatal(err)
        }
}
