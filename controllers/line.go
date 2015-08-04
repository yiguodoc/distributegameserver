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
	Start, End *Position
	Distance   float64
}

func (l *Line) withEnd(pos1, pos2 *Position) bool {
	if l.Start.equals(pos1) && l.End.equals(pos2) {
		return true
	}
	if l.End.equals(pos1) && l.Start.equals(pos2) {
		return true
	}
	return false
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
	for _, line := range ll {
		if line.withEnd(pos1, pos2) == true {
			return line
		}
	}
	return nil
}
