package main

type Handler interface {
	// Routable Route用于向Handler中注册路由！
	Routable

	//ServerHTTP 用于处理接收到的封装在ctx中的http request，
	// 解析路径并且选择正确的handlerFunc(例如SignUp方法)，进行处理后，在把响应封装好并且发送给客户端
	//实际上，ServerHTTP方法就是最内层的root filter！它被层层filterBuilder进行增强，实现AOP
	ServerHTTP(ctx *Context)
}

type Routable interface {
	Route(method string, pattern string, handlerFunc handlerFunc)
}

type handlerFunc func(ctx *Context)

var _ Handler = &HandlerBasedOnTree{}
var _ Handler = &HandlerBasedMap{}
