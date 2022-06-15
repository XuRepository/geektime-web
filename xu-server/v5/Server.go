package main

import (
	"fmt"
	"net/http"
)

type Server interface {
	Routable
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler *HandlerBasedMap
}

// Route 把目标的handleFunc func(ctx *Context) 注册到 s.handler.handlers这个map中去！
func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	s.handler.Route(method, pattern, handleFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	//把我们的HandlerBasedMap注册到服务器中作为Handler使用！
	http.Handle("/", s.handler)
	return http.ListenAndServe(address, nil)
}

func NewHttpServer(name string) Server {
	return &sdkHttpServer{ //当返回实际类型所实现的接口的时候，需要返回指针
		Name:    name,
		handler: NewHandlerBasedMap(),
	}
}

func SignUp(ctx *Context) {
	req := &signUpReq{}

	////////////////// 么有context时，使用原生的方法读json文件等处理 ////////////////

	err := ctx.ReadJson(req)
	if err != nil {

		fmt.Fprintf(ctx.W, "error %v", err)
	}

	resp := commonResponse{
		Data: 123,
	}

	err = ctx.OKJson(resp)
	if err != nil {
		fmt.Fprintf(ctx.W, "error %v ", err)
	}
	//////////////////// 没有Context时response的处理 //////////////////
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
