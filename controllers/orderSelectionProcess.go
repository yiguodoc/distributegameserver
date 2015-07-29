package controllers

import (
	"errors"
	"fmt"
	"time"
)

var (
	ERR_NO_ENOUGH_ORDER               = errors.New("订单数量太少")
	ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE = errors.New("没有订单可以分配了")
	// ERR_CANNOT_FIND_ORDER_FOR_DISTRIBUTION = errors.New("未找到要分配的订单")
	// ERR_DISTRIBUTOR_FULL                   = errors.New("配送员接收的订单已达到最大值")
	// ERR_NO_SUCH_DISTRIBUTOR                = errors.New("未查找到配送员")
)
var (
	chanOrderSelectResponse      = make(chan *OrderDistribution, 16) //接收订单分配结果的channel，保证数据通道的唯一
	chanStartOrderSelection      = make(chan int, 16)                //接收开始订单分发的消息，为后面的订单分发做准备工作
	chanCreateOrderProposal      = make(chan int, 16)                //接收分发订单的消息时，分发可以选择的订单到客户端
	chanOrderSelectAdditionalMsg = make(chan string)                 //需要订单分发补充该过程的信息时使用，一般用于掉线后重新连接，与其余用户同步数据，其中接收的内容为接收者的id，单独发送，非群发
)

func init() {
	startOrderSelectionChannel()
}
func requestOrderSelectAdditionalMsg(id string) {
	chanOrderSelectAdditionalMsg <- id
}
func orderSelectionProcess(event *SysEvent) {
	switch event.eventCode {
	case sys_event_order_select_response:
		chanOrderSelectResponse <- event.data.(*OrderDistribution)
	case sys_event_start_order_selection:
		chanStartOrderSelection <- 1
	}
}

func startOrderSelectionChannel() {
	go func() {
		for {
			select {
			case <-chanStartOrderSelection: //开始订单分发
				//通知客户端的消息
				// triggerSysEvent(NewMessageBroadcastEvent("一大波订单即将到来"))
				g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, "配送员全部准备完毕")
				g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, "一大波订单即将到来")
				//倒计时
				timer := time.Tick(1 * time.Second)
				count := 3
				DebugInfo("start timer...")
				for {
					<-timer
					DebugTraceF("timer count : %d", count)
					if count <= 0 {
						break
					}
					g_room_distributor.broadcastMsgToSubscribers(pro_timer_count_down, count)
					count--
				}
				//给负责分发订单的channel发消息
				chanCreateOrderProposal <- 1
			case <-chanCreateOrderProposal: //建立人与分发的订单的关系列表，发送到客户端
				proposals, err := createDistributionProposal(g_ordersUndistributed, g_distributors)
				if err != nil {
					DebugMustF("%s", err.Error())
				} else {
					g_room_distributor.broadcastMsgToSubscribers(pro_order_distribution_proposal, proposals)
				}
			case od := <-chanOrderSelectResponse:
				order := g_ordersDistributed.findByID(od.OrderID)
				if order != nil {
					DebugInfoF("订单[%s]已经被分配", od.OrderID)
					return
				}
				//其次订单满足分配条件，当前的条件是尚未分配
				order = g_ordersUndistributed.findByID(od.OrderID)
				if order == nil {
					DebugInfoF("没有在未分配订单列表中找到订单[%s]的信息", od.OrderID)
					// return nil, ERR_CANNOT_FIND_ORDER_FOR_DISTRIBUTION
					continue
				}
				distributor := g_distributors.find(od.DistributorID) //首先确保配送员满足订单分配条件，当前条件是已分配的订单未达到最大可接收数量
				if distributor == nil {
					DebugInfoF("没有找到配送员[%s]的信息", od.DistributorID)
					// return nil, ERR_NO_SUCH_DISTRIBUTOR
					continue
				}
				if distributor.full() {
					// return nil, ERR_DISTRIBUTOR_FULL
					DebugInfoF("配送员[%s]已经满载", od.DistributorID)
					continue
				}

				//确定结果
				// DebugTraceF("未分配订单有 %d 个", len(g_ordersUndistributed))
				g_ordersUndistributed = g_ordersUndistributed.remove(order)
				DebugTraceF("未分配订单减少到 %d 个", len(g_ordersUndistributed))
				g_ordersDistributed = append(g_ordersDistributed, order)
				distributor.acceptOrder(order)

				//将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
				distributionResult := NewOrderDistribution(order.ID, distributor.ID)
				triggerSysEvent(NewSysEvent(sys_event_order_select_result, distributionResult))
				g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_order_distribute_result, distributionResult)

				msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", od.OrderID, distributor.Name)
				g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, msg)

				if distributor.full() == true { //配送员的订单满载了
					g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, fmt.Sprintf("配送员 %s 订单满载，可以开始配送了", distributor.Name))
					g_room_distributor.broadcastMsgToSubscribers(pro_distribution_prepared, distributor.ID)
				}
				DebugInfo(msg)
				chanCreateOrderProposal <- 1
			case id := <-chanOrderSelectAdditionalMsg:
				proposals, err := createDistributionProposal(g_ordersUndistributed, g_distributors)
				if err != nil {
					DebugMustF("%s", err.Error())
				} else {
					g_room_distributor.sendMsgToSpecialSubscriber(id, pro_order_distribution_proposal, proposals)
				}
			}

		}
	}()
}

//只是生成一个分配建议，不是最终的分配结果
//分配原则：开始分配相同的N（N=配送员的数量）个，每当有订单被选定时，补充一个新的订单
func createDistributionProposal(ordersUndistributed OrderList, distributors DistributorList) (list OrderList, err error) {
	if len(ordersUndistributed) <= 0 {
		err = ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE
		return
	}
	if len(ordersUndistributed) <= 0 {
		err = ERR_NO_ENOUGH_ORDER
		return
	}
	// distributorsNotFull := distributors.notFull()
	// if len(ordersUndistributed) < len(distributorsNotFull) { //每次分配的订单数等于待接受订单的配送员的数量
	// 	DebugMustF("There is %d orders and %d distributors", len(ordersUndistributed), len(distributorsNotFull))
	// 	err = ERR_NO_ENOUGH_ORDER
	// 	return
	// }
	// var list OrderList
	if len(ordersUndistributed) >= len(distributors) {
		list = ordersUndistributed[0:len(distributors)]
	} else {
		list = ordersUndistributed
	}
	return
}
