package controllers

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	// "math/rand"
	// "time"
	"encoding/json"
)

type distributorPredictor func(*Distributor) bool

const (
	color_orange = "orange"
	color_red    = "red"
	color_grey   = "grey"
	color_purple = "purple"
)

var error_no_websocket_connection = errors.New("error_no_websocket_connection")

type CheckPoint int

func getCheckPointMap() map[string]CheckPoint {
	m := make(map[string]CheckPoint)
	m["checkpoint_flag_origin"] = checkpoint_flag_origin
	m["checkpoint_flag_order_select"] = checkpoint_flag_order_select
	m["checkpoint_flag_order_distribute"] = checkpoint_flag_order_distribute
	m["checkpoint_flag_order_distribute_over"] = checkpoint_flag_order_distribute_over
	return m
}

var (
	checkpoint_flag_origin                CheckPoint = 0
	checkpoint_flag_order_select          CheckPoint = 1
	checkpoint_flag_order_distribute      CheckPoint = 2
	checkpoint_flag_order_distribute_over CheckPoint = 3
	checkpoint_max                        CheckPoint = 4
)

// 配送员
type Distributor struct {
	ID                     string
	Name                   string
	AcceptedOrders         OrderList
	CheckPoint             CheckPoint //所处的关卡
	Online                 bool
	Color                  string          //地图上marker颜色
	StartPos, DestPos      *Position       //配送时设置的出发和目的路径点
	CurrentPos             *Position       //配送时实时所在的路径
	MaxAcceptedOrdersCount int             `json:"-"` //配送员可以接收的最大订单数量
	Conn                   *websocket.Conn `json:"-"` // Only for WebSocket users; otherwise nil.
	NormalSpeed            float64         //运行速度 km/h
	CurrentSpeed           float64         //当前运行速度，0表示停止
	Distance               float64         //所在或者将要行驶的路径长度
	line                   *Line
	TimeElapse             int64 //运行时间
	// RejectedOrders         OrderList
}

func NewDistributorFromJson(bytes []byte) *Distributor {
	var Distributor Distributor
	err := json.Unmarshal(bytes, &Distributor)
	if err != nil {
		DebugSysF("解析JSON生成 Distributor 时出错：%s", err)
		return nil
	} else {
		return &Distributor
	}
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
	return fmt.Sprintf("ID: %-3s  Name: %-4s 游戏进程：%d 可接收新订单：%t 接收的订单：%2d online:%t", this.ID, this.Name, this.CheckPoint, !this.full(), len(this.AcceptedOrders), this.IsOnline())
}
func (d *Distributor) PosString() string {
	if d.CurrentPos == nil {
		return fmt.Sprintf("ID: %-3s  Name: %-4s  未设定当前位置", d.ID, d.Name)
	}
	if d.DestPos == nil {
		return fmt.Sprintf("ID: %-3s  Name: %-4s  (%f, %f) %fkm/h", d.ID, d.Name, d.CurrentPos.Lng, d.CurrentPos.Lat, d.CurrentSpeed)
	}
	return fmt.Sprintf("ID: %-3s  Name: %-4s  (%f, %f) => (%f, %f) %fkm/h", d.ID, d.Name, d.CurrentPos.Lng, d.CurrentPos.Lat, d.DestPos.Lng, d.DestPos.Lat, d.CurrentSpeed)
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
	// return d.Is(func(dr *Distributor) bool { return dr.ID == id })
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
func (d *Distributor) setCheckPoint(check CheckPoint) {
	d.CheckPoint = check
}
func (d *Distributor) copyAll() *Distributor {
	return &Distributor{
		ID:                     d.ID,
		Name:                   d.Name,
		AcceptedOrders:         OrderList{},
		MaxAcceptedOrdersCount: d.MaxAcceptedOrdersCount,
		CheckPoint:             d.CheckPoint,
		Color:                  d.Color,
	}
}

func (d *Distributor) Is(f distributorPredictor) bool {
	return f(d)
}

type DistributorList []*Distributor

func (dl DistributorList) clone(f distributorPredictor) (l DistributorList) {
	for _, d := range dl {
		if f == nil || f(d) {
			l = append(l, d.copyAll())
		}
	}
	return
}
func (dl DistributorList) forEach(f func(*Distributor)) {
	for _, d := range dl {
		if f != nil {
			f(d)
		}
	}
}
func (dl DistributorList) filter(f distributorPredictor) (l DistributorList) {
	for _, d := range dl {
		if f == nil || f(d) {
			l = append(l, d)
		}
	}
	return
}
func (dl DistributorList) findOne(f distributorPredictor) *Distributor {
	for _, p := range dl {
		if f(p) {
			return p
		}
	}
	return nil
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

// func (l DistributorList) preparedForOrderSelect(id string) {
// 	d := l.find(id)
// 	if d != nil {
// 		d.CheckPoint = checkpoint_flag_order_select
// 	}
// }
// func (l DistributorList) setCheckPoint(id string, checkPoint CheckPoint) {
// 	if checkPoint != checkpoint_flag_origin &&
// 		checkPoint != checkpoint_flag_order_select &&
// 		checkPoint != checkpoint_flag_order_distribute {
// 		return
// 	}
// 	d := l.find(id)
// 	if d != nil {
// 		d.CheckPoint = checkPoint
// 	}
// }

// func (l DistributorList) allPreparedForOrderSelect() bool {
// 	for _, d := range l {
// 		if d.CheckPoint < checkpoint_flag_order_select {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (d DistributorList) find(id string) *Distributor {
// 	for _, d := range d {
// 		if d.ID == id {
// 			return d
// 		}
// 	}
// 	return nil
// }

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
// func (this DistributorList) notFull() (list DistributorList) {
// 	for _, d := range this {
// 		if d.full() == false {
// 			list = append(list, d)
// 		}
// 	}
// 	return
// }

// func (this DistributorList) setBroadcastChannel(chOrder, chResult <-chan *broadcastMsg) {
// 	for _, d := range this {
// 		d.chBroadcastOrder = chOrder
// 		d.chBroadcastOrderResult = chResult
// 	}

// }
