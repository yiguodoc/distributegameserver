package controllers

import (
// "container/list"
// "github.com/astaxie/beego"
// "github.com/gorilla/websocket"
// "time"
// "encoding/json"
// "strings"
// "fmt"
)

type OrderDistribution struct {
	OrderID       string
	DistributorID string
}

func NewOrderDistribution(orderID, distributorID string) *OrderDistribution {
	return &OrderDistribution{
		OrderID:       orderID,
		DistributorID: distributorID,
	}
}

type OrderDistributionList []*OrderDistribution

func (l OrderDistributionList) add(ods ...*OrderDistribution) OrderDistributionList {
	for _, od := range ods {
		l = append(l, od)
	}
	return l
}
