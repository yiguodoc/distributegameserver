package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	// "time"
	// "reflect"
)

// type objectOperate interface {
// 	GetClone(id interface{}) interface{}

// }

/*

* 作为系统中的位置点，设计上分为三类：仓库位置和路径点，而路径点包括能产生订单的点和不能产生订单的点
* 仓库点单独设定
* 路径点作为配送路线的关键节点存在
* 可以产生订单的点，既是路线关键节点，也作为订单发源地
* 配送路线以仓库作为起点，路径关键节点时，如果只有一个节点作为该节点连接，直接运行至下一个节点，否则需要做出选择


 */

const (
	POSITION_TYPE_TEMP        = -1 //临时点的标志
	POSITION_TYPE_WAREHOUSE   = 0  //仓库
	POSITION_TYPE_ORDER_ROUTE = 1  //路径节点
	POSITION_TYPE_ROUTE_ONLY  = 2  //途经点
	POSITION_TYPE_ROUTE_TEMP  = 3  //计算得出的临时点
)

// type positionTypeFilter func(*Position) bool
type predictor func(interface{}) bool
type positionPredictor func(*Position) bool

//位置，订单的产生地
type Position struct {
	ID        int
	Lng, Lat  float64
	Address   string
	PointType int
	HasOrder  bool
	mutable   bool
	// LinkedPoints PositionList //与该位置直接连接的点
}
type PositionList []*Position

func (p *Position) String() string {
	return fmt.Sprintf("ID: %3d  类型：%d  订单：%t  位置: (%f, %f)  %s", p.ID, p.PointType, p.HasOrder, p.Lng, p.Lat, p.Address)
}
func (p *Position) SimpleString() string {
	return fmt.Sprintf("(%f, %f)", p.Lng, p.Lat)
}
func (p *Position) checkMutable() {
	if p.mutable == false {
		panic("position imutable")
	}
}
func (p *Position) setMutabel(b bool) {
	p.mutable = b
}
func (p *Position) copyTemp(mutable bool) *Position {
	temp := p.copyAll(mutable)
	temp.PointType = POSITION_TYPE_TEMP
	return temp
}
func (p *Position) copyAll(mutable bool) *Position {
	return &Position{
		ID:        p.ID,
		Lng:       p.Lng,
		Lat:       p.Lat,
		Address:   p.Address,
		PointType: p.PointType,
		HasOrder:  p.HasOrder,
		mutable:   mutable,
	}
}
func (p *Position) equals(pos *Position) bool {
	if p.Lat == pos.Lat && p.Lng == pos.Lng {
		return true
	}
	return false
}
func (p *Position) setLngLat(lng, lat float64) {
	p.checkMutable()
	p.Lng = lng
	p.Lat = lat
}
func (p *Position) addLngLat(lng, lat float64) {
	p.checkMutable()
	p.Lng += lng
	p.Lat += lat
}
func (p *Position) minus(pos *Position) (lng, lat float64) {
	return p.Lng - pos.Lng, p.Lat - pos.Lat
}

func (pl PositionList) clone(f predictor) (l PositionList) {
	for _, p := range pl {
		if f == nil || f(p) {
			l = append(l, p.copyAll(true))
		}
	}
	return
}

func (pl PositionList) filter(f positionPredictor) (l PositionList) {
	for _, p := range pl {
		if f(p) {
			l = append(l, p)
		}
	}
	return
}
func (pl PositionList) findOne(f positionPredictor) *Position {
	for _, p := range pl {
		if f(p) {
			return p
		}
	}
	return nil
}
func (pl PositionList) contains(f positionPredictor) bool {
	for _, p := range pl {
		if f(p) {
			return true
		}
	}
	return false
}

// func (pl PositionList) findOne(f predictor) *Position {
// 	for _, p := range pl {
// 		if f(p) {
// 			return p
// 		}
// 	}
// 	return nil
// }

func (pl PositionList) Map(list interface{}, f func(*Position, interface{}) interface{}) interface{} {
	// func (pl PositionList) Map(f func(interface{}) interface{}) AnyArray {
	// list := []interface{}{}
	// for _, p := range pl {
	// 	list = append(list, f(p))
	// }
	// return list

	if len(pl) > 0 {
		return (pl[1:]).reduce(f(pl[0], list), f)
	} else {
		return list
	}
}

func (pl PositionList) reduce(list interface{}, f func(*Position, interface{}) interface{}) interface{} {
	if len(pl) > 0 {
		return (pl[1:]).reduce(f(pl[0], list), f)
	} else {
		return list
	}
}

// func (pl PositionList) findByID1(id int) *Position {
// 	for _, p := range pl {
// 		if p.ID == id {
// 			return p
// 		}
// 	}
// 	return nil
// }
// func (this PositionList) findLngLat1(lng, lat float64) *Position {
// 	for _, p := range this {
// 		if p.Lat == lat && p.Lng == lng {
// 			return p
// 		}
// 	}
// 	return nil
// }

// func (this PositionList) contains(pos *Position) bool {
// 	return this.findLngLat(pos.Lng, pos.Lat) != nil
// }
func (pl PositionList) ListName() string {
	return "关键点"
}
func (pl PositionList) InfoList() (l []string) {
	for _, p := range pl {
		l = append(l, p.String())
	}
	return
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

// func createPositionFilter(pointType int) positionTypeFilter {
// 	f := func(p *Position) bool {
// 		return p.PointType == pointType
// 	}
// 	return f
// }
