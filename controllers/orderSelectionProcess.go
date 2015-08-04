package controllers

import (
	"errors"
	// "fmt"
	"time"
)

var (
	ERR_no_enough_order               = errors.New("订单数量太少")
	ERR_no_enough_order_to_distribute = errors.New("没有订单可以分配了")
	ERR_cannot_find_order             = errors.New("未找到订单")
	ERR_distributor_full              = errors.New("配送员接收的订单已达到最大值")
	ERR_no_such_distributor           = errors.New("未查找到配送员")
	ERR_order_already_selected        = errors.New("订单已经被分配过")
)

//处理具体的客户端事件
//code ClientMessageTypeCode
func orderSelectionProcess(msg *MessageWithClient, unit *DistributorProcessUnit) {
	// c := sysEventCode(code.(int))
	// c := ClientMessageTypeCode(code.(int))
	switch msg.MessageType {
	case pro_on_line:
		distributor := unit.distributor
		if distributor != nil {
			//如果在分配订单中，应该推送给其正在选择的订单
			switch distributor.CheckPoint {
			case checkpoint_flag_origin:
				DebugTraceF("配送员上线，状态 %d 初始化", checkpoint_flag_origin)
			case checkpoint_flag_order_select:
				// triggerSysEvent(NewSysEvent(sys_event_order_select_additional_msg, distributor.GetID()))
				DebugTraceF("配送员上线，状态 %d 订单选择中", checkpoint_flag_order_select)
				broadOrderSelectProposal()
			case checkpoint_flag_order_distribute:
				DebugTraceF("配送员上线，状态 %d 配送中", checkpoint_flag_order_distribute)
				// triggerSysEvent(NewSysEvent(sys_event_start_order_distribution, distributor))
			}
		}

	// case sys_event_order_select_additional_msg:
	// sendOrderSelectProposal(data.(string))
	// proposals, err := createDistributionProposal(g_orders.Filter(newOrderDistributeFilter(false)), g_distributors)
	// if err != nil {
	// 	DebugMustF("%s", err.Error())
	// } else {
	// 	g_room_distributor.sendMsgToSpecialSubscriber(data.(string), pro_order_distribution_proposal, proposals)
	// }
	// case sys_event_select_order_prepared:
	case pro_prepared_for_select_order:
		m := msg.Data.(map[string]interface{})
		// m := event.data.(map[string]interface{})
		distributorID, ok := m["DistributorID"]
		if !ok {
			DebugMustF("客户端数据格式错误，无法获取配送员编号")
			return
		}
		DebugTraceF("配送员[%s %s]准备好订单的分发了", distributorID, unit.distributor.Name)
		g_distributors.preparedForOrderSelect(distributorID.(string))
		if g_distributors.allPreparedForOrderSelect() == true {
			DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
			unit.center.allUnitsProcess(NewMessageWithClient(pro_order_distribution_proposal_first, nil))
			// unit.process(NewMessageWithClient())
			// startOrderDistributionEvent := NewSysEvent(sys_event_start_select_order, nil)
			// triggerSysEvent(startOrderDistributionEvent)

		}
	case pro_order_select_response:
		unit.center.allUnitsProcess(msg)
	// case sys_event_start_select_order:
	//倒数
	case pro_timer_count_down:
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
		triggerSysEvent(NewSysEvent(sys_event_give_out_order, nil))

		// case sys_event_give_order_select_result:
		// 	distributionResult := data.(*OrderDistribution)
		// 	// distributionResult := event.data.(*OrderDistribution)
		// 	distributor := g_distributors.find(distributionResult.DistributorID)
		// 	g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_order_select_result, distributionResult)

		// 	msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", distributionResult.OrderID, distributor.Name)
		// 	g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, msg)
		// 	DebugInfo(msg)

		// case sys_event_order_select_response:
		// 	// chanOrderSelectResponse <- event.data.(*OrderDistribution)
		// 	m := data.(map[string]interface{})
		// 	// m := event.data.(map[string]interface{})
		// 	if list, err := mappedValue(m).Getter("OrderID", "DistributorID"); err == nil {
		// 		disposeOrderSelectResponse(list[0].(string), list[1].(string))
		// 	} else {
		// 		DebugMustF("客户端数据格式错误: %s", err)
		// 	}

		// case sys_event_give_out_order:
		// 	proposals, err := createDistributionProposal(g_orders.Filter(newOrderDistributeFilter(false)), g_distributors)
		// 	if err != nil {
		// 		DebugMustF("%s", err.Error())
		// 	} else {
		// 		g_room_distributor.broadcastMsgToSubscribers(pro_order_distribution_proposal, proposals)
		// 	}

	}
}
func broadOrderSelectProposal() {
	proposals, err := createDistributionProposal(g_orders.Filter(newOrderDistributeFilter(false)), g_distributors)
	if err != nil {
		DebugMustF("%s", err.Error())
	} else {
		// g_room_distributor.sendMsgToSpecialSubscriber(id, pro_order_distribution_proposal, proposals)
		g_room_distributor.broadcastMsgToSubscribers(pro_order_distribution_proposal, proposals)
	}
}
func disposeOrderSelectResponse(orderID, distributorID string) error {
	order := g_orders.findByID(orderID)
	if order == nil {
		DebugMustF("系统异常，分配不存在的订单：%s", orderID)
		return ERR_cannot_find_order
	}
	if order.distributed == true {
		DebugInfoF("订单[%s]已经被分配", orderID)
		return ERR_order_already_selected
	}

	distributor := g_distributors.find(distributorID) //首先确保配送员满足订单分配条件，当前条件是已分配的订单未达到最大可接收数量
	if distributor == nil {
		DebugInfoF("没有找到配送员[%s]的信息", distributorID)
		return ERR_no_such_distributor
	}
	if distributor.full() {
		DebugInfoF("配送员[%s]已经满载", distributorID)
		return ERR_distributor_full
	}

	//确定结果
	order.distributed = true
	distributor.acceptOrder(order)
	DebugTraceF("未分配订单减少到 %d 个", len(g_orders.Filter(newOrderDistributeFilter(false))))
	return nil
	// //将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
	// distributionResult := NewOrderDistribution(order.ID, distributor.ID)
	// triggerSysEvent(NewSysEvent(sys_event_give_order_select_result, distributionResult))
	// // g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_order_distribute_result, distributionResult)

	// // msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", od.OrderID, distributor.Name)
	// // g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, msg)

	// if distributor.full() == true { //配送员的订单满载了
	// 	g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, fmt.Sprintf("配送员 %s 订单满载", distributor.Name))
	// 	// g_room_distributor.broadcastMsgToSubscribers(pro_distribution_prepared, distributor.ID)
	// 	triggerSysEvent(NewSysEvent(sys_event_distribution_prepared, distributor))
	// }
	// triggerSysEvent(NewSysEvent(sys_event_give_out_order, nil))
}

//只是生成一个分配建议，不是最终的分配结果
//分配原则：开始分配相同的N（N=配送员的数量）个，每当有订单被选定时，补充一个新的订单
func createDistributionProposal(ordersUndistributed OrderList, distributors DistributorList) (list OrderList, err error) {
	if len(ordersUndistributed) <= 0 {
		err = ERR_no_enough_order_to_distribute
		return
	}
	if len(ordersUndistributed) <= 0 {
		err = ERR_no_enough_order
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
