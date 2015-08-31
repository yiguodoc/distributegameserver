package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	// "time"
)

type Region struct {
	LatMin, LatMax float64
	LngMin, LngMax float64
	Color          string
	Code           string
}

func NewRegion(code, color string, latMin, latMax, lngMin, lngMax float64) *Region {
	return &Region{
		Code:   code,
		Color:  color,
		LatMin: latMin,
		LatMax: latMax,
		LngMin: lngMin,
		LngMax: lngMax,
	}
}
func (r *Region) in(pos *Position) bool {
	return pos != nil &&
		pos.Lat >= r.LatMin && pos.Lat < r.LatMax &&
		pos.Lng >= r.LngMin && pos.Lng < r.LngMax
}

type RegionList []*Region

func (rl RegionList) findRegion(pos *Position) *Region {
	for _, r := range rl {
		if r.in(pos) {
			return r
		}
	}
	return nil
}

var orderCount = 0

func generateOrderID() string {
	orderCount++
	return fmt.Sprintf("9001%05d", orderCount)
}

// 订单
type Order struct {
	ID          string
	GeoSrc      *Position
	Distributed bool //分配状态
	Signed      bool //签收状态
	Region      *Region
}

func NewOrder(id string, pos *Position) *Order {
	region := g_regions.findRegion(pos)
	if region == nil {
		panic(fmt.Sprintf("没有定义点所属的区域：%s", pos))
	}
	return &Order{
		ID:     id,
		GeoSrc: pos,
		Region: region,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("ID: %s Signed: %t  Distributed: %t Address: %s Region: %s", o.ID, o.Signed, o.Distributed, o.GeoSrc.Address, o.Region.Code)
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

func (l OrderList) findOne(f predictor) *Order {
	for _, o := range l {
		if f(o) {
			return o
		}
	}
	return nil
}

// func (l OrderList) findByID(id string) *Order {
// 	for _, o := range l {
// 		if o.ID == id {
// 			return o
// 		}
// 	}
// 	return nil
// }
func (l OrderList) all(f predictor) bool {
	for _, o := range l {
		if f(o) == false {
			return false
		}
	}
	return true
}
func (l OrderList) contains(f predictor) bool {
	return l.findOne(f) != nil
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
func (ol OrderList) Filter(f predictor) (l OrderList) {
	for _, o := range ol {
		if f(o) == true {
			l = append(l, o)
		}
	}
	return
}
