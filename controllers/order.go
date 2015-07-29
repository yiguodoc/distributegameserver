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

func (o *Order) String() string {
	return fmt.Sprintf("ID: %s  Address: %s", o.ID, o.GeoSrc.Address)
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
