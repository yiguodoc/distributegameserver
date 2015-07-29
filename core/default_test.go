package core

import (
	// "net/http"
	// "net/http/httptest"
	"testing"
	// "runtime"
	// "path/filepath"
	// _ "distributionGame/routers"

	// "github.com/astaxie/beego"
	// . "github.com/smartystreets/goconvey/convey"
	// "time"
)

func init() {
}

func TestRoute(t *testing.T) {
	warehouse := NewPosition("warehouse", "1", "1", POSITION_TYPE_WAREHOUSE)
	pA := NewPosition("A", "2", "1", POSITION_TYPE_ROUTE_ONLY)
	warehouse.addLinks(pA)
	list := warehouse.getLinks()
	if len(list) != 1 {
		panic("only one route")
	}
	pB := NewPosition("B", "2", "2", POSITION_TYPE_ROUTE_ONLY)
	pC := NewPosition("C", "3", "1", POSITION_TYPE_ROUTE_ONLY)
	pA.addLinks(pB, pC)
	list = pA.getLinks()
	if len(list) != 2 {
		panic("there are two route points")
	}

	// time.Sleep(1 * time.Second)
}
