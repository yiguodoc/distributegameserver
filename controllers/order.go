package controllers

import (
	// "github.com/astaxie/beego"
	// "errors"
	"fmt"
	"math/rand"
	// "time"
)

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
	SignTime     int  //签收时间
	SelectedTime int  //被选择的时间
	Score        int  //分值
	// Region       *Region
}

func NewOrder(id string, pos *Position) *Order {
	return &Order{
		ID:     id,
		GeoSrc: pos,
		Score:  pos.Score,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("ID: %s Signed: %t  Distributed: %t Address: %s ", o.ID, o.Signed, o.Distributed, o.GeoSrc.Address)
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
func (l OrderList) forEach(f func(*Order)) {
	if f == nil || len(l) <= 0 {
		return
	}
	f(l[0])
	l[1:].forEach(f)
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
func (ol OrderList) totalScore(score int) int {
	if len(ol) <= 0 {
		return score
	}
	return ol[1:].totalScore(score + ol[0].Score)
}
