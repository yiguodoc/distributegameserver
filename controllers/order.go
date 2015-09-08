package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	"math/rand"
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
	ID           string
	GeoSrc       *Position
	Distributed  bool //分配状态
	Signed       bool //签收状态
	Region       *Region
	SignTime     int //签收时间
	SelectedTime int //被选择的时间
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
func (o *Order) sign(time int) {
	o.Signed = true
	o.SignTime = time
}
func (o *Order) distribute(time int) {
	o.Distributed = true
	o.SelectedTime = time
}

type OrderList []*Order

func (l OrderList) random(rander *rand.Rand, list OrderList) OrderList {
	if len(l) <= 1 {
		return append(list, l...)
	}
	r := rander.Intn(len(l))
	last := append(append(OrderList{}, l[:r]...), l[r+1:]...)
	return last.random(rander, append(list, l[r]))
}

func (l OrderList) findOne(f predictor) *Order {
	if len(l) <= 0 || f == nil {
		return nil
	}
	if f(l[0]) {
		return l[0]
	} else {
		return l[1:].findOne(f)
	}

	// for _, o := range l {
	// 	if f(o) {
	// 		return o
	// 	}
	// }
	// return nil
}

func (l OrderList) all(f predictor) bool {
	if len(l) <= 0 || f == nil {
		return true
	}
	if f(l[0]) == false {
		return false
	} else {
		return l[1:].all(f)
	}
	// for _, o := range l {
	// 	if f(o) == false {
	// 		return false
	// 	}
	// }
	// return true
}
func (l OrderList) contains(f predictor) bool {
	return l.findOne(f) != nil
}

//
func (ll OrderList) remove(f func(*Order) bool, list ...OrderList) (l OrderList) {
	var ol OrderList
	if len(list) <= 0 {
		ol = OrderList{}
	} else {
		ol = list[0]
	}
	if len(ll) <= 0 {
		return ol
	}
	if f(ll[0]) {
		return ll[1:].remove(f, append(ol, ll[0]))
	} else {
		return ll[1:].remove(f, ol)
	}
	// for _, o := range ll {
	// 	if o.ID != order.ID {
	// 		l = append(l, o)
	// 	}
	// }
	// return
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
func (ol OrderList) Filter(f func(*Order) bool) (l OrderList) {
	for _, o := range ol {
		if f(o) == true {
			l = append(l, o)
		}
	}
	return
}
