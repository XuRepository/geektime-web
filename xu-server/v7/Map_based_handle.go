package main

import "net/http"

type Routable interface {
	Route(method string, pattern string, handlerFunc func(ctx *Context))
}

type HandlerBasedMap struct {
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBasedMap) Route(method string, pattern string, handlerFunc func(ctx *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

func NewHandlerBasedMap() *HandlerBasedMap {
	return &HandlerBasedMap{handlers: make(map[string]func(ctx *Context))}
}

//实现Handler接口，这样就可以把HandlerBasedMap作为服务器的Handler，处理注册的路由！
//实现的功能是  处理http请求，判断请求的类型和路径，创建context封装resp和req，从map中取到相应的handlerFunc进行处理！
func (h *HandlerBasedMap) ServeHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	if handler, OK := h.handlers[key]; OK {
		handler(NewContext(ctx.W, ctx.R))
	} else {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("NOT FOUND"))
	}
}

//key 请求方法+请求路径
func (h *HandlerBasedMap) key(method string, pattern string) string {
	return method + "#" + pattern
}
