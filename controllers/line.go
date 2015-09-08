package controllers

import (
	// "errors"
	// "encoding/json"
	"fmt"
	// "github.com/BurntSushi/toml"
	// "github.com/astaxie/beego"
	// "github.com/ungerik/go-dry"
	// "os"
)

type LatLng struct {
	Lat, Lng float64
}

func (ll *LatLng) String() string {
	return fmt.Sprintf("(%f, %f)", ll.Lng, ll.Lat)
}

type Line struct {
	Start, End     *Position
	Distance       float64
	Busy           bool
	DistributorsOn DistributorList
	// DistributorsOn map[string]*Distributor
}

func (l *Line) isBusy() bool {
	return l.Busy
}
func (l *Line) busy() {
	l.Busy = true
}
func (l *Line) nobusy() {
	l.Busy = false
}
func (l *Line) DistributorsCount() int {
	return len(l.DistributorsOn)
}
func (l *Line) isDistributorOn(id string) bool {
	// _, ok := l.DistributorsOn[id]
	// return ok

	return l.DistributorsOn.findOne(func(d *Distributor) bool { return d.ID == id }) != nil
}
func (l *Line) addDistributor(d *Distributor) {
	if l.DistributorsOn == nil {
		l.DistributorsOn = DistributorList{}
	}
	if l.isDistributorOn(d.ID) == false {
		l.DistributorsOn = append(l.DistributorsOn, d)
	}
}
func (l *Line) removeDistributor(id string) {
	// delete(l.DistributorsOn, id)
	l.DistributorsOn = l.DistributorsOn.filter(func(d *Distributor) bool { return d.ID != id })
}
func (l *Line) withEnd(pos1, pos2 *Position) bool {
	f := func(pos *Position) positionPredictor {
		return func(pIn *Position) bool {
			return pIn.equals(pos)
		}
	}
	pl := PositionList{l.Start, l.End}
	return pl.contains(f(pos1)) && pl.contains(f(pos2))
	// if l.Start.equals(pos1) && l.End.equals(pos2) {
	// 	return true
	// }
	// if l.End.equals(pos1) && l.Start.equals(pos2) {
	// 	return true
	// }
	// return false
}
func (l *Line) String() string {
	return fmt.Sprintf("line: %f米 (%f, %f) => (%f, %f)", l.Distance, l.Start.Lng, l.Start.Lat, l.End.Lng, l.End.Lat)
}

type LineList []*Line

func (ll LineList) ListName() string {
	return "路径"
}
func (ll LineList) InfoList() (l []string) {
	for _, line := range ll {
		l = append(l, line.String())
	}
	return
}
func (ll LineList) find(pos1, pos2 *Position) *Line {
	if len(ll) <= 0 {
		return nil
	}
	if ll[0].withEnd(pos1, pos2) == true {
		return ll[0]
	} else {
		return ll[1:].find(pos1, pos2)
	}
	// for _, line := range ll {
	// 	if line.withEnd(pos1, pos2) == true {
	// 		return line
	// 	}
	// }
	// return nil
}
