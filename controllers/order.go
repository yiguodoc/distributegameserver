package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	// "time"
)

var orderCount = 0

func generateOrderID() string {
	orderCount++
	return fmt.Sprintf("9001%05d", orderCount)
}

// type Any interface{}

// func (a Any) convertToOrderList() (l OrderList) {
// 	list := a.([]interface{})
// 	for _, in := range list {
// 		if in != nil {
// 			l = append(l, in.(*Order))
// 		}
// 	}
// 	return
// }
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
	// if len(aa) <= 0 {
	// 	return nil
	// }	// a0 := aa[0]
	// switch a0.(type) {
	// case *Order:
	// 	f := func(source AnyArray) (l OrderList) {
	// 		for _, s := range source {
	// 			l = append(l, s.(*Order))
	// 		}
	// 		return
	// 	}
	// 	return f(aa)
	// default:
	// 	panic(fmt.Sprintf("类型 %T 没有定义，无法转换", a0))

	// }
	// return nil
}

func convertToOrderList(list AnyArray) (l OrderList) {
	// func convertToOrderList(o interface{}) (l OrderList) {
	// list := o.([]interface{})

	for _, in := range list {
		if in != nil {
			l = append(l, in.(*Order))
		}
	}
	return
}

// 订单
type Order struct {
	ID          string
	GeoSrc      *Position
	Distributed bool //分配状态
	Signed      bool //签收状态
}
type orderDistributeFilter func(*Order) bool

func newOrderDistributeFilter(distributed bool) orderDistributeFilter {
	f := func(o *Order) bool {
		return o.Distributed == distributed
	}
	return f
}
func NewOrder(id string, pos *Position) *Order {
	return &Order{
		ID:     id,
		GeoSrc: pos,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("ID: %s Signed: %t  Distributed: %t Address: %s ", o.ID, o.Signed, o.Distributed, o.GeoSrc.Address)
}
func (o *Order) isDistributed() bool {
	return o.Distributed
}
func (o *Order) isSigned() bool {
	return o.Signed
}
func (o *Order) sign() {
	o.Signed = true
}
func (o *Order) distribute() {
	o.Distributed = true
}

type OrderList []*Order

func (l OrderList) findByID(id string) *Order {
	for _, o := range l {
		if o.ID == id {
			return o
		}
	}
	return nil
}
func (l OrderList) contains(id string) bool {
	return l.findByID(id) != nil
}

//
func (l OrderList) remove(order *Order) (list OrderList) {
	for _, o := range l {
		if o.ID != order.ID {
			list = append(list, o)
		}
	}
	return
}
func (ol OrderList) ListName() string {
	return "订单"
}
func (ol OrderList) InfoList() (list []string) {
	for _, o := range ol {
		list = append(list, o.String())
	}
	return
}
func (ol OrderList) Filter(filter orderDistributeFilter) (l OrderList) {
	for _, o := range ol {
		if filter(o) == true {
			l = append(l, o)
		}
	}
	return
}
