package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	// "time"
)

/*

* 作为系统中的位置点，设计上分为三类：仓库位置和路径点，而路径点包括能产生订单的点和不能产生订单的点
* 仓库点单独设定
* 路径点作为配送路线的关键节点存在
* 可以产生订单的点，既是路线关键节点，也作为订单发源地
* 配送路线以仓库作为起点，路径关键节点时，如果只有一个节点作为该节点连接，直接运行至下一个节点，否则需要做出选择


 */

const (
	POSITION_TYPE_WAREHOUSE   = 0 //仓库
	POSITION_TYPE_ORDER_ROUTE = 1 //路径节点
	POSITION_TYPE_ROUTE_ONLY  = 2 //途经点
	POSITION_TYPE_ROUTE_TEMP  = 3 //计算得出的临时点
)

//位置，订单的产生地
type Position struct {
	ID        int
	Lng, Lat  float64
	Address   string
	PointType int
	HasOrder  bool
	// LinkedPoints PositionList //与该位置直接连接的点
}

func (p *Position) String() string {
	return fmt.Sprintf("ID: %3d  类型：%d  订单：%t  位置: (%f, %f)  %s", p.ID, p.PointType, p.HasOrder, p.Lng, p.Lat, p.Address)
}
func NewPosition(id int, address string, lng, lat float64, ptype int, hasOrder bool) *Position {
	return &Position{
		ID:        id,
		Address:   address,
		Lng:       lng,
		Lat:       lat,
		PointType: ptype,
		HasOrder:  hasOrder,
	}
}

// func (this *Position) getLinks() PositionList {
// 	return this.LinkedPoints
// }
// func (this *Position) addLinks(list ...*Position) {
// 	// this.LinkedPoints = list
// 	for _, pos := range list {
// 		if this.LinkedPoints.contains(pos) == false {
// 			this.LinkedPoints = append(this.LinkedPoints, pos)
// 		}
// 	}
// }

type PositionList []*Position

// //从一系列的地点中产生一组订单
func (pl PositionList) createSimulatedOrders(idGenerator func() string) (list OrderList) {
	for _, p := range pl {
		if p.HasOrder {
			list = append(list, NewOrder(idGenerator(), p))
		}
	}
	return
}
func (this PositionList) find(pos *Position) *Position {
	for _, p := range this {
		if p.Lat == pos.Lat && p.Lng == pos.Lng {
			return p
		}
	}
	return nil
}
func (this PositionList) contains(pos *Position) bool {
	return this.find(pos) != nil
}
func (pl PositionList) ListName() string {
	return "关键点"
}
func (pl PositionList) InfoList() (l []string) {
	for _, p := range pl {
		l = append(l, p.String())
	}
	return
}
