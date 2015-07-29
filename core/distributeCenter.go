package core

import (
	// "github.com/astaxie/beego"
	"errors"
	// "fmt"
	// "time"
)

/*

* 分配中心告诉配送员一个频道，向收听者（并不一定是全部配送员）发送广播
* 广播分两类：分发订单和订单分配结果
* 分发订单的广播中含有一个特殊频道，用于接收配送员申请

 */
const (
	BROADCASTMSGTYPE_ORDER_DISTRIBUTE        = 0
	BROADCASTMSGTYPE_ORDER_DISTRIBUTE_RESULT = 1
)

var (
	ERR_NO_ENOUGH_ORDER                    = errors.New("订单数量太少")
	ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE      = errors.New("没有订单可以分配了")
	ERR_NO_SUCH_DISTRIBUTOR                = errors.New("未查找到配送员")
	ERR_DISTRIBUTOR_FULL                   = errors.New("配送员接收的订单已达到最大值")
	ERR_CANNOT_FIND_ORDER_FOR_DISTRIBUTION = errors.New("未找到要分配的订单")
)

// type broadcastMsg struct {
// 	msgType     int
// 	order       *Order
// 	distributor *Distributor
// 	chMsg       chan<- *Distributor //对于distributor只能往里写
// }

// //分发订单的广播消息
// func NewDistributeBroadcastMsg(order *Order, ch chan<- *Distributor) *broadcastMsg {
// 	return &broadcastMsg{
// 		msgType: BROADCASTMSGTYPE_ORDER_DISTRIBUTE,
// 		order:   order,
// 		chMsg:   ch,
// 	}
// }

// func NewDistributeResultBroadcastMsg(order *Order, distributor *Distributor) *broadcastMsg {
// 	return &broadcastMsg{
// 		msgType:     BROADCASTMSGTYPE_ORDER_DISTRIBUTE_RESULT,
// 		order:       order,
// 		distributor: distributor,
// 	}
// }

// 订单分配中心，管理着参与的订单和配送员的分配过程
type DistributeCenter struct {
	OrdersDistributed   OrderList       //已经分配的订单
	OrdersUndistributed OrderList       //尚未分配的订单
	Distributors        DistributorList //所有配送员
	// chBroadcastDistributeOrder       chan<- *broadcastMsg //分配中心发布订单分配消息的频道，每个配送员如果需要接收订单必须订阅该频道
	// chBroadcastDistributeOrderResult chan<- *broadcastMsg //分配中心发布订单分配结果消息的频道
}

func NewDistributeCenter(orders OrderList, distributors DistributorList) *DistributeCenter {
	// ch1 := make(chan *broadcastMsg, len(distributors)) //支持缓冲，同时让所有的配送员接收到消息
	// ch2 := make(chan *broadcastMsg, len(distributors)) //支持缓冲，同时让所有的配送员接收到消息
	// chReadOnly1 := (<-chan *broadcastMsg)(ch1)
	// chReadOnly2 := (<-chan *broadcastMsg)(ch2)
	// distributors.setBroadcastChannel(chReadOnly1, chReadOnly2)
	return &DistributeCenter{
		OrdersUndistributed: orders,
		Distributors:        distributors,
		// chBroadcastDistributeOrder:       chan<- *broadcastMsg(ch1), //发送要分配的订单
		// chBroadcastDistributeOrderResult: chan<- *broadcastMsg(ch2), //发送订单分配完成的结果
	}
}

// func (this *DistributeCenter) start() {
// 	this.Distributors.startListening()
// }

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

func (l OrderDistributionList) add(ods ...OrderDistribution) OrderDistributionList {
	for _, od := range ods {
		l = append(l, od)
	}
	return l
}

//开始进行订单的分配
//只是生成一个分配建议，不是最终的分配结果
func (d *DistributeCenter) createDistributionProposal() (list OrderDistributionList, err error) {
	if len(d.OrdersUndistributed) <= 0 {
		err = ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE
		return
		// return nil, ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE
	}
	distributorsNotFull := d.Distributors.notFull()
	if len(d.OrdersUndistributed) < len(distributorsNotFull) {
		DebugMustF("There is %d orders and %d distributors", len(d.OrdersUndistributed), len(distributorsNotFull))
		err = ERR_NO_ENOUGH_ORDER
		return
	}
	// d.Distributors.startListening()
	order := d.OrdersUndistributed[0]
	for _, distributor := range d.Distributors {
		list = list.add(NewOrderDistribution(order.ID, distributor.ID))
	}
	// chMsg := make(chan *Distributor)
	// distributeMsg := NewDistributeBroadcastMsg(order, chan<- *Distributor(chMsg))
	// for i := 0; i < len(d.Distributors); i++ {
	// 	d.chBroadcastDistributeOrder <- distributeMsg
	// }
	// broadcastDistributors := d.Distributors.notFull()
	// broadcastDistributors.orderComing(order, chMsg)
	// resultDistributor := <-chMsg
	// distributeResultMsg := NewDistributeResultBroadcastMsg(order, resultDistributor)
	// for i := 0; i < len(d.Distributors); i++ {
	// 	d.chBroadcastDistributeOrderResult <- distributeResultMsg
	// }
	// d.OrdersDistributed = append(d.OrdersDistributed, order)
	// d.OrdersUndistributed = d.OrdersUndistributed[1:]
	// DebugTraceF("接收订单(%s)的配送员信息：%s", order.ID, resultDistributor.ID)
	// resultDistributor.acceptOrder(order)
	return
}

//接收订单分配的反馈
//如果反馈被接受，则成为分配结果，相应地订单也会发生变化
func (d *DistributeCenter) acceptDistributionResponse(od *OrderDistribution) (*OrderDistribution, error) {
	//首先确保配送员满足订单分配条件，当前条件是已分配的订单未达到最大可接收数量
	distributor := d.Distributors.find(od.DistributorID)
	if distributor == nil {
		return nil, ERR_NO_SUCH_DISTRIBUTOR
	}
	if distributor.full() {
		return nil, ERR_DISTRIBUTOR_FULL
	}
	//其次订单满足分配条件，当前的条件是尚未分配
	order := d.OrdersUndistributed.find(od.OrderID)
	if order == nil {
		return nil, ERR_CANNOT_FIND_ORDER_FOR_DISTRIBUTION
	}
	//确定结果
	d.OrdersUndistributed.remove(order)
	d.OrdersDistributed = append(d.OrdersDistributed, order)
	distributor.acceptOrder(order)
	return od, nil
}
