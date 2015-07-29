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
	Start, End *LatLng
}

func (l *Line) String() string {
	return fmt.Sprintf("line: %s => %s", l.Start, l.End)
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
