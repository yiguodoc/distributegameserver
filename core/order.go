package core

import (
// "github.com/astaxie/beego"
// "errors"
// "fmt"
// "time"
)

// 订单
type Order struct {
	ID     string
	GeoSrc *Position
}

func NewOrder(id string, pos *Position) *Order {
	return &Order{
		ID:     id,
		GeoSrc: pos,
	}
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

//
func (l OrderList) remove(o *Order) (list OrderList) {
	for _, o := range l {
		if o.ID != o.ID {
			list = append(list, o)
		}
	}
	return
}
