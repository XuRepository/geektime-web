package main

import (
	"fmt"
	"time"
)

type Filter func(ctx *Context) // 就是过滤器，一层一层封装Filter，实现AOP。原始的Filter就是本来的Handler(类似于SignUp)，一层一层封装

// FilterBuilder 用于对参数中的next这个Filter进行增强！
type FilterBuilder func(next Filter) Filter // filter建造者，返回的是一个Filter，参数next是被封装的Filter，用于对next进行增强！责任链模式

// MetricsFilterBuilder1 是一个FilterBuilder，用于对参数中的next这个Filter进行增强！
func MetricsFilterBuilder1(next Filter) Filter {

	return func(ctx *Context) {
		println("[Filter1][[MetricsFilterBuilder1]进入...")
		start := time.Now().Nanosecond()
		next(ctx) //next就是上一层传入的 被增强和封装的filter，在next(ctx)执行前后进行方法增强
		end := time.Now().Nanosecond()
		println("[Filter1][[MetricsFilterBuilder1]退出...")

		fmt.Printf("执行用时：%v", end-start)
	}
}

// MetricsFilterBuilder2 是一个FilterBuilder
func MetricsFilterBuilder2(next Filter) Filter {
	return func(ctx *Context) {
		println("[Filter2][[MetricsFilterBuilder2]进入...")

		next(ctx) //next就是上一层传入的 被增强和封装的filter，在next(ctx)执行前后进行方法增强

		println("[Filter2][[MetricsFilterBuilder2]退出...")
	}
}

var _ FilterBuilder = MetricsFilterBuilder1 //判断 右侧  是否是实现目标的type 或者 接口
