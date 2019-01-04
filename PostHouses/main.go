package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"sss/PostHouses/handler"
	"sss/PostHouses/subscriber"

	example "sss/PostHouses/proto/example"
	"github.com/micro/go-grpc"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostHouses"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostHouses", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostHouses", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
