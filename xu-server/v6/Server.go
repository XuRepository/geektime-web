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
	root    Filter // root过滤器，被层层Filter增强的handler，在这里是Handler.ServerHTTP,当请求过来的时候，执行被多个过滤器增强的func(ctx)方法
}

// Route 把目标的handleFunc func(ctx *Context) 注册到 s.handler.handlers这个map中去！
func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	s.handler.Route(method, pattern, handleFunc)
}

func (s *sdkHttpServer) Start(address string) error {

	//定义一个handler处理逻辑，为w，r封装ctx并且调用filter实现AOP！通过root层层封装，最终调用我们的handler的ServeHTTP方法！
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		s.root(ctx)
	}

	//这样一来从 “/”路径进来的http请求都会被上面定义的handler处理
	http.HandleFunc("/", handler)
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 第2个参数是当前服务器要绑定的Filter，用于对原有处理path的handler方法进行增强
func NewHttpServer(name string, builders ...FilterBuilder) Server {

	handler := NewHandlerBasedMap()

	//最基层的root filter，用于调用handler的ServeHTTP，处理最基础的http请求！
	var root = func(ctx *Context) {
		handler.ServeHTTP(ctx.W, ctx.R)
	}

	//使用builders中的filterBuilder对server的 root Filter 进行增强建造。
	//倒序遍历，从后往前进行层层包裹封装，先进入的filter被封装在内层
	for i := len(builders) - 1; i >= 0; i-- {
		builder := builders[i]
		root = builder(root)
	}

	return &sdkHttpServer{ //当返回实际类型所实现的接口的时候，需要返回指针
		Name:    name,
		handler: handler,
		root:    root,
	}
}

// SignUp 登录的handler
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
