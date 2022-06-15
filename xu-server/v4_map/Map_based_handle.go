package main

import "net/http"

type HandlerBasedMap struct {
	handlers map[string]func(ctx *Context)
}

func NewHandlerBasedMap() *HandlerBasedMap {
	return &HandlerBasedMap{handlers: make(map[string]func(ctx *Context))}
}

//实现Handler接口，这样就可以把HandlerBasedMap作为服务器的Handler，处理注册的路由！
func (h *HandlerBasedMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.key(request.Method, request.URL.Path)
	if handler, OK := h.handlers[key]; OK {
		handler(NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("NOT FOUND"))
	}
}

//key 请求方法+请求路径
func (h *HandlerBasedMap) key(method string, pattern string) string {
	return method + "#" + pattern
}
