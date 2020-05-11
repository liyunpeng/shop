// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: user/user.proto

/*
Package go_micro_srv_user is a generated protocol buffer package.

It is generated from these files:
	user/user.proto

It has these top-level messages:
	User
	RegisterRequest
	LoginRequest
	UpdatePasswordRequest
	Response
*/
package go_micro_srv_user

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

// Client API for UserService service

type UserService interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*Response, error)
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*Response, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...client.CallOption) (*Response, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.user"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "UserService.Register", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "UserService.Login", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "UserService.UpdatePassword", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	Register(context.Context, *RegisterRequest, *Response) error
	Login(context.Context, *LoginRequest, *Response) error
	UpdatePassword(context.Context, *UpdatePasswordRequest, *Response) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		Register(ctx context.Context, in *RegisterRequest, out *Response) error
		Login(ctx context.Context, in *LoginRequest, out *Response) error
		UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, out *Response) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) Register(ctx context.Context, in *RegisterRequest, out *Response) error {
	return h.UserServiceHandler.Register(ctx, in, out)
}

func (h *userServiceHandler) Login(ctx context.Context, in *LoginRequest, out *Response) error {
	return h.UserServiceHandler.Login(ctx, in, out)
}

func (h *userServiceHandler) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, out *Response) error {
	return h.UserServiceHandler.UpdatePassword(ctx, in, out)
}
