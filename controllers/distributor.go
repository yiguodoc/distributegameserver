package controllers

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	// "math/rand"
	// "time"
)

const (
	color_orange = "orange"
	color_red    = "red"
	color_grey   = "grey"
	color_purple = "purple"
)

var error_no_websocket_connection = errors.New("error_no_websocket_connection")

var (
	checkpoint_flag_origin           = 0
	checkpoint_flag_order_select     = 1
	checkpoint_flag_order_distribute = 2
)

// 配送员
type Distributor struct {
	ID                     string
	Name                   string
	AcceptedOrders         OrderList
	CheckPoint             int //所处的关卡
	Online                 bool
	Color                  string          //地图上marker颜色
	StartPos, DestPos      *Position       //配送时设置的出发和目的路径点
	CurrentPos             *Position       //配送时实时所在的路径
	MaxAcceptedOrdersCount int             `json:"-"` //配送员可以接收的最大订单数量
	Conn                   *websocket.Conn `json:"-"` // Only for WebSocket users; otherwise nil.
	// preparedForOrderSelect bool            `json:"-"` //为分发订单做好准备了
	// RejectedOrders         OrderList
}

func NewDistributor(id, name string, maxCount int, color string) *Distributor {
	return &Distributor{
		ID:                     id,
		Name:                   name,
		AcceptedOrders:         OrderList{},
		MaxAcceptedOrdersCount: maxCount,
		CheckPoint:             checkpoint_flag_origin,
		Color:                  color,
	}
}
func (this *Distributor) String() string {
	return fmt.Sprintf("ID: %-5s  Name: %-10s 游戏进程：%d  可接收新订单：%t  接收的订单：%2d   online:%t   ", this.ID, this.Name, this.CheckPoint, !this.full(), len(this.AcceptedOrders), this.IsOnline())
}
func (this *Distributor) full() bool {
	return len(this.AcceptedOrders) >= this.MaxAcceptedOrdersCount
}
func (this *Distributor) acceptOrder(order *Order) {
	this.AcceptedOrders = append(this.AcceptedOrders, order)
}
func (d *Distributor) IsOnline() bool {
	return d.Conn != nil
}
func (d *Distributor) GetID() string {
	return d.ID
}
func (d *Distributor) SetConn(conn *websocket.Conn) {
	d.Conn = conn
	d.Online = true
}
func (d *Distributor) SendBinaryMessage(msg []byte) error {
	if d.Conn != nil {
		// DebugTraceF("distributor [%s] send message to client", d.ID)
		return d.Conn.WriteMessage(websocket.TextMessage, msg)
	}
	return error_no_websocket_connection
}
func (d *Distributor) IdEqals(id string) bool {
	return d.ID == id
}
func (d *Distributor) SetOffline() error {
	defer func() {
		d.Online = false
	}()
	if d.Conn != nil {
		if err := d.Conn.Close(); err == nil {
			d.Conn = nil
			DebugInfoF("[%s] OffLine WebSocket closed", d.ID)
		} else {
			DebugMustF("[%s] OffLine,But close websocket err: %s", d.ID, err)
			return err
		}
	}
	// d.preparedForOrderSelect = false
	return nil
}

// func (this *Distributor) startListening() {
// 	go func() {
// 		msg := <-this.chBroadcastOrder
// 		if msg.msgType == BROADCASTMSGTYPE_ORDER_DISTRIBUTE {
// 			if this.full() == false {
// 				time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
// 				msg.chMsg <- this
// 				// DebugTraceF("%s try to accept order %s", this.ID, msg.order.ID)
// 			}
// 		}
// 	}()
// 	go func() {
// 		msg := <-this.chBroadcastOrderResult
// 		if msg.msgType == BROADCASTMSGTYPE_ORDER_DISTRIBUTE_RESULT {
// 			resultDistributor := msg.distributor
// 			// order := msg.order
// 			if resultDistributor.ID == this.ID {
// 				// DebugTraceF("ME %s Get order %s", this.ID, order.ID)
// 				// this.AcceptedOrders = append(this.AcceptedOrders, order)
// 			} else {
// 				// DebugTraceF("NOT ME %s But %s get order %s", this.ID, resultDistributor.ID, order.ID)
// 			}
// 		}
// 	}()
// }

// //订单分配通知
// func (this *Distributor) orderComing(order *Order, chMsg chan *Distributor) {
// 	if this.full() == false {
// 		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
// 		chMsg <- this
// 		// DebugTraceF("%s try to accept order %s", this.ID, order.ID)
// 		resultDistributor := <-chMsg
// 		if resultDistributor.ID == this.ID {
// 			// DebugTraceF("ME %s Get order %s", this.ID, order.ID)
// 		} else {
// 			// DebugTraceF("NOT ME But %s get order %s", resultDistributor.ID, order.ID)
// 		}
// 	}
// }

type DistributorList []*Distributor

func (l DistributorList) preparedForOrderSelect(id string) {
	d := l.find(id)
	if d != nil {
		d.CheckPoint = checkpoint_flag_order_select
	}
}
func (l DistributorList) allPreparedForOrderSelect() bool {
	for _, d := range l {
		if d.CheckPoint < checkpoint_flag_order_select {
			return false
		}
	}
	return true
}
func (d DistributorList) find(id string) *Distributor {
	for _, d := range d {
		if d.ID == id {
			return d
		}
	}
	return nil
}

// func (this DistributorList) startListening() {
// 	for _, d := range this {
// 		d.startListening()
// 	}
// }
// func (this DistributorList) orderComing(order *Order, chMsg chan *Distributor) {
// 	for _, d := range this {
// 		go d.orderComing(order, chMsg)
// 	}
// }
func (this DistributorList) notFull() (list DistributorList) {
	for _, d := range this {
		if d.full() == false {
			list = append(list, d)
		}
	}
	return
}

// func (this DistributorList) setBroadcastChannel(chOrder, chResult <-chan *broadcastMsg) {
// 	for _, d := range this {
// 		d.chBroadcastOrder = chOrder
// 		d.chBroadcastOrderResult = chResult
// 	}

// }

func (this DistributorList) ListName() string {
	return "配送员信息"
}
func (this DistributorList) InfoList() (list []string) {
	for _, d := range this {
		list = append(list, d.String())
	}
	return
}
