// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/example/example.proto

/*
Package go_micro_srv_Order is a generated protocol buffer package.

It is generated from these files:
	proto/example/example.proto

It has these top-level messages:
	Message
	AddOrderRequest
	AddOrderResponse
	GetOrdersRequest
	GetOrdersResponse
	HandleOrderRequest
	HandleOrderResponse
	CommentOrderRequest
	CommentOrderResponse
*/
package go_micro_srv_Order

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Example service

type ExampleService interface {
	AddOrder(ctx context.Context, in *AddOrderRequest, opts ...client.CallOption) (*AddOrderResponse, error)
	GetOrders(ctx context.Context, in *GetOrdersRequest, opts ...client.CallOption) (*GetOrdersResponse, error)
	HandleOrder(ctx context.Context, in *HandleOrderRequest, opts ...client.CallOption) (*HandleOrderResponse, error)
	CommentOrder(ctx context.Context, in *CommentOrderRequest, opts ...client.CallOption) (*CommentOrderResponse, error)
}

type exampleService struct {
	c    client.Client
	name string
}

func NewExampleService(name string, c client.Client) ExampleService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.Order"
	}
	return &exampleService{
		c:    c,
		name: name,
	}
}

func (c *exampleService) AddOrder(ctx context.Context, in *AddOrderRequest, opts ...client.CallOption) (*AddOrderResponse, error) {
	req := c.c.NewRequest(c.name, "Example.AddOrder", in)
	out := new(AddOrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) GetOrders(ctx context.Context, in *GetOrdersRequest, opts ...client.CallOption) (*GetOrdersResponse, error) {
	req := c.c.NewRequest(c.name, "Example.GetOrders", in)
	out := new(GetOrdersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) HandleOrder(ctx context.Context, in *HandleOrderRequest, opts ...client.CallOption) (*HandleOrderResponse, error) {
	req := c.c.NewRequest(c.name, "Example.HandleOrder", in)
	out := new(HandleOrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleService) CommentOrder(ctx context.Context, in *CommentOrderRequest, opts ...client.CallOption) (*CommentOrderResponse, error) {
	req := c.c.NewRequest(c.name, "Example.CommentOrder", in)
	out := new(CommentOrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Example service

type ExampleHandler interface {
	AddOrder(context.Context, *AddOrderRequest, *AddOrderResponse) error
	GetOrders(context.Context, *GetOrdersRequest, *GetOrdersResponse) error
	HandleOrder(context.Context, *HandleOrderRequest, *HandleOrderResponse) error
	CommentOrder(context.Context, *CommentOrderRequest, *CommentOrderResponse) error
}

func RegisterExampleHandler(s server.Server, hdlr ExampleHandler, opts ...server.HandlerOption) error {
	type example interface {
		AddOrder(ctx context.Context, in *AddOrderRequest, out *AddOrderResponse) error
		GetOrders(ctx context.Context, in *GetOrdersRequest, out *GetOrdersResponse) error
		HandleOrder(ctx context.Context, in *HandleOrderRequest, out *HandleOrderResponse) error
		CommentOrder(ctx context.Context, in *CommentOrderRequest, out *CommentOrderResponse) error
	}
	type Example struct {
		example
	}
	h := &exampleHandler{hdlr}
	return s.Handle(s.NewHandler(&Example{h}, opts...))
}

type exampleHandler struct {
	ExampleHandler
}

func (h *exampleHandler) AddOrder(ctx context.Context, in *AddOrderRequest, out *AddOrderResponse) error {
	return h.ExampleHandler.AddOrder(ctx, in, out)
}

func (h *exampleHandler) GetOrders(ctx context.Context, in *GetOrdersRequest, out *GetOrdersResponse) error {
	return h.ExampleHandler.GetOrders(ctx, in, out)
}

func (h *exampleHandler) HandleOrder(ctx context.Context, in *HandleOrderRequest, out *HandleOrderResponse) error {
	return h.ExampleHandler.HandleOrder(ctx, in, out)
}

func (h *exampleHandler) CommentOrder(ctx context.Context, in *CommentOrderRequest, out *CommentOrderResponse) error {
	return h.ExampleHandler.CommentOrder(ctx, in, out)
}
