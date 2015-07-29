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
	"time"
)

func init() {
}

func testDistributeOrders(t *testing.T) {
	orders := OrderList{
		NewOrder("1", nil),
		NewOrder("2", nil),
		NewOrder("3", nil),
		NewOrder("4", nil),
		NewOrder("5", nil),
		NewOrder("6", nil),
	}
	distributors := DistributorList{
		NewDistributor("distributor1", "0d1", 2),
		NewDistributor("distributor2", "0d2", 2),
		NewDistributor("distributor3", "0d3", 2),
	}
	distributeCenter := NewDistributeCenter(orders, distributors)

	proposalList, err := distributeCenter.createDistributionProposal()
	if err != nil {
		panic(err)
	}
	if len(proposalList) != len(distributors) {
		panic("订单分配建议应等于配送员数量")
	}
	DebugPrintList_Trace(distributors)
}
