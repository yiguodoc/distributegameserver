package core

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	"math/rand"
	"time"
)

// 配送员
type Distributor struct {
	ID                     string
	Name                   string
	AcceptedOrders         OrderList
	MaxAcceptedOrdersCount int                  //配送员可以接收的最大订单数量
	chBroadcastOrder       <-chan *broadcastMsg //分配中心发送订单消息频道
	chBroadcastOrderResult <-chan *broadcastMsg //分配中心发送订单分配结果消息频道
	// RejectedOrders         OrderList
}

func NewDistributor(id, name string, maxCount int) *Distributor {
	return &Distributor{
		ID:                     id,
		Name:                   name,
		AcceptedOrders:         OrderList{},
		MaxAcceptedOrdersCount: maxCount,
		// RejectedOrders:         OrderList{},
	}
}
func (this *Distributor) String() string {
	return fmt.Sprintf("ID: %-5s  Name: %-10s 可接收新订单：%t     接收的订单：%2d      ", this.ID, this.Name, !this.full(), len(this.AcceptedOrders))
}
func (this *Distributor) full() bool {
	return len(this.AcceptedOrders) >= this.MaxAcceptedOrdersCount
}
func (this *Distributor) acceptOrder(order *Order) {
	this.AcceptedOrders = append(this.AcceptedOrders, order)
}
func (this *Distributor) startListening() {
	go func() {
		msg := <-this.chBroadcastOrder
		if msg.msgType == BROADCASTMSGTYPE_ORDER_DISTRIBUTE {
			if this.full() == false {
				time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
				msg.chMsg <- this
				// DebugTraceF("%s try to accept order %s", this.ID, msg.order.ID)
			}
		}
	}()
	go func() {
		msg := <-this.chBroadcastOrderResult
		if msg.msgType == BROADCASTMSGTYPE_ORDER_DISTRIBUTE_RESULT {
			resultDistributor := msg.distributor
			// order := msg.order
			if resultDistributor.ID == this.ID {
				// DebugTraceF("ME %s Get order %s", this.ID, order.ID)
				// this.AcceptedOrders = append(this.AcceptedOrders, order)
			} else {
				// DebugTraceF("NOT ME %s But %s get order %s", this.ID, resultDistributor.ID, order.ID)
			}
		}
	}()
}

//订单分配通知
func (this *Distributor) orderComing(order *Order, chMsg chan *Distributor) {
	if this.full() == false {
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		chMsg <- this
		// DebugTraceF("%s try to accept order %s", this.ID, order.ID)
		resultDistributor := <-chMsg
		if resultDistributor.ID == this.ID {
			// DebugTraceF("ME %s Get order %s", this.ID, order.ID)
		} else {
			// DebugTraceF("NOT ME But %s get order %s", resultDistributor.ID, order.ID)
		}
	}
}

type DistributorList []*Distributor

func (this DistributorList) startListening() {
	for _, d := range this {
		d.startListening()
	}
}
func (this DistributorList) orderComing(order *Order, chMsg chan *Distributor) {
	for _, d := range this {
		go d.orderComing(order, chMsg)
	}
}
func (this DistributorList) notFull() (list DistributorList) {
	for _, d := range this {
		if d.full() == false {
			list = append(list, d)
		}
	}
	return
}
func (this DistributorList) setBroadcastChannel(chOrder, chResult <-chan *broadcastMsg) {
	for _, d := range this {
		d.chBroadcastOrder = chOrder
		d.chBroadcastOrderResult = chResult
	}

}

func (this DistributorList) ListName() string {
	return "配送员信息"
}
func (this DistributorList) InfoList() (list []string) {
	for _, d := range this {
		list = append(list, d.String())
	}
	return
}
