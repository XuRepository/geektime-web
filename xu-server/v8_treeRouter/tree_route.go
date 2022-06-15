package main

import (
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root *node
}

func NewHandlerBasedOnTree() *HandlerBasedOnTree {
	println("new tree router")
	root := &node{}
	return &HandlerBasedOnTree{root: root}
}

// ServerHTTP 方法就是处理接受到的http请求（请求被封装在ctx中），然后handlerBasedOnTree选择正确的handler进行处理（通过tree匹配）
func (h *HandlerBasedOnTree) ServerHTTP(ctx *Context) {
	handler, found := h.findRouter(ctx.R.URL.Path) //通过path获取到路径，然后进行tree匹配查找handleFunc
	if found {
		//找到目标handlerFunc  调用处理！
		handler(ctx)
	} else {
		//没找到对应的handleFunc，返回错误信息
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("handlerFunc Not Found 404"))
		return
	}
}

//Route 方法就是往树种插入节点，往Server中注册节点进去
func (h *HandlerBasedOnTree) Route(method string, pattern string, handlerFunc handlerFunc) {
	//1,把pattern字符串分割成为一个数组，方便进行逐个判断
	pattern = strings.Trim(pattern, "/") //去掉路径首尾的pattern
	paths := strings.Split(pattern, "/")

	cur := h.root //指向当前handler的根节点,root本身不存东西，第一个路径的匹配从root.child中间去找
	for index, path := range paths {
		//如果在子节点中找到匹配cur节点path的，就判断为当前path成功。index之前包括index本身 的是已经匹配成功的，index之后的是未匹配的
		matchChild, found := h.findMatchChild(cur, path)
		if found {
			//找到匹配的子节点，迭代继续遍历
			cur = matchChild
		} else {
			//没找到，就以当前节点为根据，向后创建子节点！
			h.createSubTree(cur, paths[index:], handlerFunc)
			return
		}
	}
}

//findRouter 通过path获取到路径，然后进行tree匹配查找handleFunc
func (h *HandlerBasedOnTree) findRouter(pattern string) (handlerFunc, bool) {
	//1,把pattern字符串分割成为一个数组，方便进行逐个判断
	pattern = strings.Trim(pattern, "/") //去掉路径首尾的pattern
	paths := strings.Split(pattern, "/")

	cur := h.root
	for _, path := range paths {
		matcheChind, found := h.findMatchChild(cur, path)
		if !found {
			return nil, false
		}
		cur = matcheChind
	}
	//寻找结束，没有return的话，就是匹配了整个路径，现在看看cur是否是叶子结点，不是叶子节点，就算查找完毕。
	if cur.handler == nil {
		return nil, false
	} else {
		return cur.handler, true
	}
}

func (h *HandlerBasedOnTree) findMatchChild(cur *node, path string) (*node, bool) {
	for _, child := range cur.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}
func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handler handlerFunc) {
	cur := root //需要迭代遍历来创建子节点，使用cur作为临时变量
	for _, path := range paths {
		nn := newNode(path) //新节点
		cur.children = append(cur.children, nn)

		cur = nn //继续迭代
	}
	//迭代结束，cur位置是叶子结点，对该节点设置handlerFunc
	cur.handler = handler
}

type node struct {
	path     string
	children []*node //切片

	//如果在叶子结点上面实现了匹配，可以调用handler来处理
	handler handlerFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}
