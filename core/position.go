package core

import (
// "github.com/astaxie/beego"
// "errors"
// "fmt"
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
	POSITION_TYPE_WAREHOUSE   = 0
	POSITION_TYPE_ORDER_ROUTE = 1
	POSITION_TYPE_ROUTE_ONLY  = 2
)

//位置，订单的产生地
type Position struct {
	Lon, Lat     string
	Name         string
	PointType    int
	LinkedPoints PositionList //与该位置直接连接的点
}

func NewPosition(name, lon, lat string, ptype int) *Position {
	return &Position{
		Name:      name,
		Lon:       lon,
		Lat:       lat,
		PointType: ptype,
	}
}
func (this *Position) getLinks() PositionList {
	return this.LinkedPoints
}
func (this *Position) addLinks(list ...*Position) {
	// this.LinkedPoints = list
	for _, pos := range list {
		if this.LinkedPoints.contains(pos) == false {
			this.LinkedPoints = append(this.LinkedPoints, pos)
		}
	}
}

type PositionList []*Position

//从一系列的地点中产生一组订单
func (this PositionList) commitOrders() (list OrderList) {
	return
}
func (this PositionList) find(pos *Position) *Position {
	for _, p := range this {
		if p.Lat == pos.Lat && p.Lon == pos.Lon {
			return p
		}
	}
	return nil
}
func (this PositionList) contains(pos *Position) bool {
	return this.find(pos) != nil
}
