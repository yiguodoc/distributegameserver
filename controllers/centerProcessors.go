package controllers

import (
	"errors"
	"fmt"
	// "math"
	"strconv"
	"time"
)

/*
经纬度变化1秒/帧
路径两点的经纬度和长度是确定的，配送员的速度确定，那么运行时间就可以确定，根据系统时间与现实时间的比例，确定系统时间每秒行走的距离和改变的经纬度

*/

var (
	ERR_no_enough_order               = errors.New("订单数量太少")
	ERR_no_enough_order_to_distribute = errors.New("没有订单可以分配了")
	ERR_cannot_find_order             = errors.New("未找到订单")
	ERR_distributor_full              = errors.New("配送员接收的订单已达到最大值")
	ERR_no_such_distributor           = errors.New("未查找到配送员")
	ERR_order_already_selected        = errors.New("订单已经被分配过")
)

// func pro_rank_changed_handlerGenerator(o interface{}) MessageWithClientHandler {
// 	center := o.(*DistributorProcessUnitCenter)
// 	f := func(msg *MessageWithClient) {
// 		center.distributors.forEach(func(d *Distributor) { center.wsRoom.sendMsgToSpecialSubscriber(d.ID, pro_2c_end_game, d) })
// 	}
// 	return f
// }
func pro_end_game_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		if distributor := center.distributors.findOne(func(d *Distributor) bool { return d.ID == msg.TargetID }); distributor != nil {
			//计算得分
			//签收完一个订单得到该订单对应的分数，没有完成的订单减去惩罚分数
			unSignedOrders := distributor.AcceptedOrders.Filter(func(o *Order) bool { return o.Signed == false })
			distributor.Score -= (len(unSignedOrders))
			//计算排名
			center.distributors.Rank()
			DebugPrintList_Info(center.distributors)
			distributor.setCheckPoint(checkpoint_flag_order_distribute_over)
			center.stopUnit(distributor.ID)
			center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_end_game, distributor)
			center.distributors.forEach(func(d *Distributor) {
				if d.ID != distributor.ID {
					center.wsRoom.sendMsgToSpecialSubscriber(d.ID, pro_2c_rank_change, d)
				}
			})
		}

	}
	return f
}

// func stopUnit(center *DistributorProcessUnitCenter, )
func pro_move_from_node_to_route_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		line := msg.Data.(*Line)
		if line.DistributorsCount() >= 2 {
			line.busy()
			DebugInfoF("line BUSY %s", line)
			line.DistributorsOn.forEach(func(d *Distributor) {
				d.CurrentSpeed = d.NormalSpeed / 2
				center.wsRoom.sendMsgToSpecialSubscriber(d.ID, pro_2c_speed_change, d)
			})
			// for id, d := range line.DistributorsOn {
			// 	center.wsRoom.sendMsgToSpecialSubscriber(id, pro_2c_speed_change, d)
			// }
		}
	}
	return f
}

func pro_move_from_route_to_node_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		line := msg.Data.(*Line)
		if line.DistributorsCount() < 2 {
			line.nobusy()
			DebugInfoF("line NOBUSY %s ", line)
			line.DistributorsOn.forEach(func(d *Distributor) {
				d.CurrentSpeed = d.NormalSpeed
				center.wsRoom.sendMsgToSpecialSubscriber(d.ID, pro_2c_speed_change, d)
			})
			// for id, d := range line.DistributorsOn {
			// 	center.wsRoom.sendMsgToSpecialSubscriber(id, pro_2c_speed_change, d)
			// }

		}
	}
	return f
}

//订单选择的请求，可以对请求作出一些限制，例如当前的位置必须处于仓库等等
func pro_order_select_response_handlerGenerator(o interface{}) MessageWithClientHandler {
	// unit := o.(*DistributorProcessUnit)
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		if msg.Data == nil {
			return
		}
		m := msg.Data.(map[string]interface{})
		var list []interface{}

		list, err := mappedValue(m).Getter("OrderID", "DistributorID")
		if err != nil {
			DebugMustF("客户端数据格式错误: %s", err)
			return
		}

		orderID := list[0].(string)
		distributorID := list[1].(string)
		distributor := center.distributors.findOne(func(d *Distributor) bool { return d.ID == distributorID })
		if distributor == nil {
			DebugSysF("没有找到配送员[%s]的信息", distributorID)
			return
		}
		//系统设定，配送员需要在仓库位置才能选择订单
		qualifier := func(d *Distributor) bool {
			list := center.mapData.Points.filter(func(pos *Position) bool { return pos.PointType == POSITION_TYPE_WAREHOUSE })
			curPos := d.CurrentPos
			return list.contains(func(pos *Position) bool { return curPos.equals(pos) })
		}
		if distributor.Is(qualifier) == false {
			DebugInfoF("配送员 %s 所处位置无法选择订单", distributor.Name)
			return
		}

		if err := disposeOrderSelectResponse(orderID, distributor, center.distributors, center.orders); err != nil {
			DebugInfoF("处理订单分配时的信息提示：%s", err)
			center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_order_select_result, nil, err.Error(), strconv.Itoa(distributor.TimeElapse))
			return
		}
		//将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
		center.wsRoom.sendMsgToSpecialSubscriber(distributorID, pro_2c_order_select_result, distributor)

		log := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", orderID, distributor.Name)
		// center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, msg)
		DebugInfoF(log)

		// if distributor.fullyLoaded() == true { //配送员的订单满载了
		// 	log = fmt.Sprintf("配送员 %s 订单满载", distributor.Name)
		// 	center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, log)
		// 	DebugInfoF(log)
		// 	// distributor.setCheckPoint(checkpoint_flag_order_distribute)
		// 	center.distributors.forOne(func(d *Distributor) bool {
		// 		if d.ID == distributorID {
		// 			d.CheckPoint = checkpoint_flag_order_distribute
		// 			DebugInfoF("配送员 %s 状态变化 => 配送环节", d.Name)
		// 			return true
		// 		}
		// 		return false
		// 	})
		// 	center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_order_full, distributor)
		// }
		sendOrderProposal(center)

	}
	return f
}

func pro_game_start_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {

		msgList := []string{"配送员全部准备完毕进入订单选择环节", "一大波订单即将到来"}
		for _, msg := range msgList {
			center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast_before_game_start, msg)
			time.Sleep(2 * time.Second)
		}
		//倒计时
		timer := time.Tick(1 * time.Second)
		count := 3
		// DebugInfo("start timer...")
		for {
			<-timer
			DebugTraceF("timer count : %d", count)
			if count <= 0 {
				break
			}
			center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast_before_game_start, count)
			count--
		}
		center.wsRoom.broadcastMsgToSubscribers(pro_2c_all_prepared_4_select_order, nil)
		sendOrderProposal(center)
		center.startAlltUnit()
		center.startGameTiming()
	}
	return f
}
func sendOrderProposal(center *DistributorProcessUnitCenter) {
	// if len(center.distributors.filter(func(d *Distributor) bool { return d.fullyLoaded() == false })) > 0 {
	// if len(center.distributors.notFull()) > 0 {
	// broadOrderSelectProposal(center.distributors, center.orders)
	if proposals, err := getOrderSelectProposal(center.distributors, center.orders); err == nil {
		center.wsRoom.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
	} else {
		DebugInfoF("%s", err)
	}
	// }
}

func pro_prepared_for_select_order_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	// unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		m := msg.Data.(map[string]interface{})
		// m := event.data.(map[string]interface{})
		distributorID, ok := m["DistributorID"]
		if !ok {
			DebugMustF("客户端数据格式错误，无法获取配送员编号")
			return
		}
		distributor := center.distributors.findOne(func(d *Distributor) bool { return d.ID == distributorID.(string) })
		if distributor == nil {
			DebugSysF("配送员 %s 不应该出现在该中心内，系统错误", distributorID.(string))
			return
		}
		DebugInfoF("配送员[%s %s]准备好订单的分发了", distributorID, distributor.Name)
		// center.distributors.preparedForOrderSelect(distributorID.(string))
		setCheckPoint := func(d *Distributor) bool {
			if d.ID == distributorID.(string) {
				d.CheckPoint = checkpoint_flag_order_select
				return true
			}
			return false
		}
		center.distributors.forOne(setCheckPoint)
		// distributorsPrepared := center.distributors.filter(func(d *Distributor) bool {
		// 	return d.CheckPoint >= checkpoint_flag_order_select
		// })
		// if len(distributorsPrepared) >= len(center.distributors) {
		// 	// if center.distributors.allPreparedForOrderSelect() == true {
		// 	DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
		// 	center.Process(NewMessageWithClient(pro_game_start, "", nil))
		// } else {
		// 	DebugInfoF("还有 %d 个配送员未准备完毕", len(center.distributors)-len(distributorsPrepared))
		// }
		if center.distributors.every(func(d *Distributor) bool { return d.CheckPoint >= checkpoint_flag_order_select }) {
			DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
			center.Process(NewMessageWithClient(pro_game_start, "", nil))
		} else {
			DebugInfoF("还有 %d 个配送员未准备完毕", len(center.distributors.filter(func(d *Distributor) bool {
				return d.CheckPoint < checkpoint_flag_order_select
			})))

		}
	}
	return f
}
func pro_off_line_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		DebugInfoF("%s", msg)
		center.stopUnit(msg.TargetID)
	}
	return f
}
func pro_on_line_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		DebugInfoF("处理消息 %s", msg)
		distributor := center.distributors.findOne(func(d *Distributor) bool { return d.ID == msg.TargetID })
		if distributor != nil {
			if distributor.IsOriginal() { //掉线后重新上线自动启动
				DebugTraceF("配送员上线，状态 %d 初始化", checkpoint_flag_origin)
				//设置默认起始点
				filter := func(pos *Position) bool {
					return pos.PointType == POSITION_TYPE_WAREHOUSE
				}
				warehouses := center.mapData.Points.filter(filter)
				if len(warehouses) > 0 {
					distributor.StartPos = warehouses[0] //
					distributor.CurrentPos = distributor.StartPos.copyTemp(true)
				} else {
					DebugSysF("无法设置出发点")
				}
				distributor.NormalSpeed = defaultSpeed
				distributor.CurrentSpeed = defaultSpeed
			}
			center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_map_data, center.mapData)
			center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_distributor_info, distributor)
			onReconnect(center, distributor)
		}
	}
	return f
}
func onReconnect(center *DistributorProcessUnitCenter, distributor *Distributor) {
	center.startUnit(distributor.ID)
	//如果在分配订单中，应该推送给其正在选择的订单
	switch distributor.CheckPoint {
	case checkpoint_flag_origin:
		DebugTraceF("配送员上线，状态 %d 初始化", checkpoint_flag_origin)
	case checkpoint_flag_order_select:
		DebugTraceF("配送员上线，状态 %d 订单选择中", checkpoint_flag_order_select)
		// broadOrderSelectProposal(center.distributors, center.orders)
		if proposals, err := getOrderSelectProposal(center.distributors, center.orders); err == nil {
			// center.wsRoom.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
			center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_order_distribution_proposal, proposals)
		} else {
			DebugInfoF("%s", err)
		}
	case checkpoint_flag_order_distribute:
		DebugTraceF("配送员上线，状态 %d 配送中", checkpoint_flag_order_distribute)
	case checkpoint_flag_order_distribute_over:
		DebugTraceF("配送员上线，状态 %d 配送完成", checkpoint_flag_order_distribute_over)
	}
}

func getOrderSelectProposal(distributors DistributorList, orders OrderList) (list OrderList, err error) {
	proposals, err := createDistributionProposal(orders.Filter(func(o *Order) bool { return o.Distributed == false }), distributors)
	// proposals, err := createDistributionProposal(orders.Filter(newOrderDistributeFilter(false)), distributors)
	if err != nil {
		return nil, err
	}
	return proposals, nil
}
func disposeOrderSelectResponse(orderID string, distributor *Distributor, distributors DistributorList, orders OrderList) error {
	order := orders.findOne(func(o interface{}) bool { return o.(*Order).ID == orderID })
	if order == nil {
		DebugMustF("系统异常，分配不存在的订单：%s", orderID)
		return ERR_cannot_find_order
	}
	if order.isDistributed() == true {
		DebugInfoF("订单[%s]已经被分配", orderID)
		return ERR_order_already_selected
	}

	// if distributor.fullyLoaded() {
	// 	DebugInfoF("配送员[%s]已经满载", distributor.Name)
	// 	return ERR_distributor_full
	// }

	//确定结果
	order.distribute(distributor.TimeElapse)
	distributor.acceptOrder(order)
	DebugTraceF("未分配订单减少到 %d 个", len(orders.Filter(func(o *Order) bool { return o.Distributed == false })))
	return nil
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
