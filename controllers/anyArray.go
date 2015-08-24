package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	// "time"
)

type Sys_Type string

var (
	Sys_Type_Order Sys_Type = "*Order"
)

type AnyArray []interface{}

func (aa AnyArray) without(o interface{}) (l AnyArray) {
	for _, a := range aa {
		if a != o {
			l = append(l, a)
		}
	}
	return
}
func (aa AnyArray) transform(destType Sys_Type) interface{} {
	aa = aa.without(nil)
	switch destType {
	case Sys_Type_Order:
		f := func(source AnyArray) (l OrderList) {
			for _, s := range source {
				l = append(l, s.(*Order))
			}
			return
		}
		return f(aa)
	default:
		panic(fmt.Sprintf("类型 %s 没有定义，无法转换", destType))
	}
}
