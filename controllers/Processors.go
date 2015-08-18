package controllers

import (
	"errors"
	"fmt"
	"math"
	"time"
)

/*
经纬度变化1秒/帧
路径两点的经纬度和长度是确定的，配送员的速度确定，那么运行时间就可以确定，根据系统时间与现实时间的比例，确定系统时间每秒行走的距离和改变的经纬度

*/

var (
	realityToSystemTimeRatio float64 = 5  //现实时间与系统时间的比例
	defaultSpeed             float64 = 20 //默认配送员行驶速度 20km/h
	timePerFrame             float64 = 1
)
var (
	ERR_no_enough_order               = errors.New("订单数量太少")
	ERR_no_enough_order_to_distribute = errors.New("没有订单可以分配了")
	ERR_cannot_find_order             = errors.New("未找到订单")
	ERR_distributor_full              = errors.New("配送员接收的订单已达到最大值")
	ERR_no_such_distributor           = errors.New("未查找到配送员")
	ERR_order_already_selected        = errors.New("订单已经被分配过")
)

func pro_move_from_node_to_route_handlerGenerator(o interface{}) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		line := msg.Data.(*Line)
		if line.DistributorsCount() >= 2 {
			line.busy()
			DebugInfoF("line BUSY %s", line)
		}
	}
	return f
}

func pro_move_from_route_to_node_handlerGenerator(o interface{}) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		line := msg.Data.(*Line)
		if line.DistributorsCount() < 2 {
			line.nobusy()
			DebugInfoF("line NOBUSY %s ", line)
		}
	}
	return f
}
func pro_order_select_response_handlerGenerator(o interface{}) MessageWithClientHandler {
	// unit := o.(*DistributorProcessUnit)
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {
		m := msg.Data.(map[string]interface{})
		if list, err := mappedValue(m).Getter("OrderID", "DistributorID"); err == nil {
			orderID := list[0].(string)
			distributorID := list[1].(string)
			if err := disposeOrderSelectResponse(orderID, distributorID, center.distributors, center.orders); err != nil {
				DebugInfoF("处理订单分配时的信息提示：%s", err)
			} else {
				//将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
				distributionResult := NewOrderDistribution(orderID, distributorID)
				distributor := center.distributors.find(distributorID)
				center.wsRoom.sendMsgToSpecialSubscriber(distributorID, pro_2c_order_select_result, distributionResult)

				msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", distributionResult.OrderID, distributor.Name)
				center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, msg)
				DebugInfo(msg)

				if distributor.full() == true { //配送员的订单满载了
					center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, fmt.Sprintf("配送员 %s 订单满载", distributor.Name))
					center.wsRoom.broadcastMsgToSubscribers(pro_2c_distribution_prepared, distributor.ID)
					distributor.setCheckPoint(checkpoint_flag_order_distribute)
					// triggerSysEvent(NewSysEvent(sys_event_distribution_prepared, distributor))
				}
				if len(center.distributors.notFull()) > 0 {
					// broadOrderSelectProposal(center.distributors, center.orders)
					if proposals, err := getOrderSelectProposal(center.distributors, center.orders); err == nil {
						center.wsRoom.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
					} else {
						DebugInfoF("%s", err)
					}
				}
			}
		} else {
			DebugMustF("客户端数据格式错误: %s", err)
		}
	}
	return f
}

func pro_order_distribution_proposal_first_handlerGenerator(o interface{}) MessageWithClientHandler {
	center := o.(*DistributorProcessUnitCenter)
	f := func(msg *MessageWithClient) {

		center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, "配送员全部准备完毕")
		center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, "一大波订单即将到来")
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
			center.wsRoom.broadcastMsgToSubscribers(pro_timer_count_down, count)
			count--
		}
		if len(center.distributors.notFull()) > 0 {
			// broadOrderSelectProposal(center.distributors, center.orders)
			if proposals, err := getOrderSelectProposal(center.distributors, center.orders); err == nil {
				center.wsRoom.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
			} else {
				DebugInfoF("%s", err)
			}
		}
	}
	return f
}

// func pro_order_select_response_handlerGenerator(unit *DistributorProcessUnit) MessageWithClientHandler {
// 	f := func(msg *MessageWithClient) {
// 		unit.center.Process(msg)
// 	}
// 	return f
// }
func pro_prepared_for_select_order_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		m := msg.Data.(map[string]interface{})
		// m := event.data.(map[string]interface{})
		distributorID, ok := m["DistributorID"]
		if !ok {
			DebugMustF("客户端数据格式错误，无法获取配送员编号")
			return
		}
		DebugTraceF("配送员[%s %s]准备好订单的分发了", distributorID, unit.distributor.Name)
		unit.center.distributors.preparedForOrderSelect(distributorID.(string))
		if unit.center.distributors.allPreparedForOrderSelect() == true {
			DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
			unit.center.Process(NewMessageWithClient(pro_order_distribution_proposal_first, "", nil))
			// unit.process(NewMessageWithClient())
			// startOrderDistributionEvent := NewSysEvent(sys_event_start_select_order, nil)
			// triggerSysEvent(startOrderDistributionEvent)

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
	// unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		DebugInfoF("处理消息 %s", msg)
		center.startUnit(msg.TargetID)
		distributor := center.distributors.find(msg.TargetID)
		if distributor != nil {
			//如果在分配订单中，应该推送给其正在选择的订单
			switch distributor.CheckPoint {
			case checkpoint_flag_origin:
				DebugTraceF("配送员上线，状态 %d 初始化", checkpoint_flag_origin)
			case checkpoint_flag_order_select:
				// triggerSysEvent(NewSysEvent(sys_event_order_select_additional_msg, distributor.GetID()))
				DebugTraceF("配送员上线，状态 %d 订单选择中", checkpoint_flag_order_select)
				// broadOrderSelectProposal(center.distributors, center.orders)
				if proposals, err := getOrderSelectProposal(center.distributors, center.orders); err == nil {
					center.wsRoom.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
				} else {
					DebugInfoF("%s", err)
				}
			case checkpoint_flag_order_distribute:
				DebugTraceF("配送员上线，状态 %d 配送中", checkpoint_flag_order_distribute)
				warehouses := center.mapData.Points.filter(createPositionFilter(POSITION_TYPE_WAREHOUSE))
				if len(warehouses) > 0 {
					if distributor.StartPos == nil {
						distributor.StartPos = warehouses[0] //.copy(false)
					}
					if distributor.CurrentPos == nil { //没有保存的位置信息，设置仓库为默认的出发点
						distributor.CurrentPos = distributor.StartPos.copyTemp(true)
					}
				} else {
					DebugSysF("无法设置出发点")
				}
				type d struct {
					Distributor *Distributor
					MapData     *MapData
				}
				center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_distribution_prepared, &d{distributor, center.mapData})
			}
		}
	}
	return f
}
func pro_game_time_tick_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		distributor := unit.distributor
		//----------------------------------------------------------------------------
		//计算行走的坐标位置
		if distributor.NormalSpeed > 0 {
			if distributor.StartPos != nil && distributor.DestPos != nil && distributor.line != nil {
				if distributor.StartPos.equals(distributor.DestPos) == false {
					DebugTraceF("配送员 %s 运行路线 %s => %s", distributor.Name, distributor.StartPos.SimpleString(), distributor.DestPos.SimpleString())
					line := distributor.line
					crtSpeed := distributor.NormalSpeed
					if line.isBusy() == true {
						crtSpeed = crtSpeed * 0.5
					}
					distributor.CurrentSpeed = crtSpeed
					totalTime := distributor.Distance * 60 * 60 / (crtSpeed * 1000) / realityToSystemTimeRatio //系统中运行路程所花费的时间
					totalFrames := totalTime / timePerFrame                                                    //一共大约这么多帧就可以走完
					//使用绝对值差距大的作为分片的标准
					totalLng := distributor.DestPos.Lng - distributor.StartPos.Lng
					totalLat := distributor.DestPos.Lat - distributor.StartPos.Lat
					lngPerFrame := totalLng / totalFrames
					latPerFrame := totalLat / totalFrames
					// DebugTraceF("pos change per frame lng %f  lat %f", lngPerFrame, latPerFrame)
					lng, lat := distributor.DestPos.minus(distributor.CurrentPos) //是否已经足够接近目标点
					// DebugTraceF("pos gap lng %f  lat %f", lng, lat)
					if math.Abs(lng) < math.Abs(lngPerFrame) || math.Abs(lat) < math.Abs(latPerFrame) {
						distributor.CurrentPos.addLngLat(lng, lat)
						//已经到达目标点，运动停止
						// distributor.StartPos.setLngLat(distributor.DestPos.Lng, distributor.DestPos.Lat) //
						distributor.StartPos = unit.center.mapData.Points.findLngLat(distributor.DestPos.Lng, distributor.DestPos.Lat)
						distributor.DestPos = nil
						line.removeDistributor(distributor.ID)
						distributor.line = nil
						unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_move_to_new_position, distributor) //通知客户端移动到新坐标
						unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_reach_route_node, distributor)     //通知客户端移动到新坐标
						DebugInfoF("配送员已经行驶到目标点 %s", distributor)
						DebugTraceF("配送员实时位置：%s", distributor.PosString())
						//配送员从路上转移到节点
						unit.center.Process(NewMessageWithClient(pro_move_from_route_to_node, distributor.ID, line))
					} else {
						distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
						unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_move_to_new_position, distributor) //通知客户端移动到新坐标
						DebugTraceF("配送员实时位置：%s", distributor.PosString())
						//配送员从节点到路上
						line.addDistributor(distributor)
						unit.center.Process(NewMessageWithClient(pro_move_from_node_to_route, distributor.ID, line))
					}

				}
			}
		}
	}
	return f
}
func pro_reset_destination_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		m := msg.Data.(map[string]interface{})
		if list, err := mappedValue(m).Getter("PositionID", "DistributorID"); err == nil {
			posWanted := unit.center.mapData.Points.findByID(int(list[0].(float64)))
			if posWanted == nil {
				DebugMustF("重置目标点出错，不存在编号为 %d 的节点", int(list[0].(float64)))
				return
			}
			distributor := unit.center.distributors.find(list[1].(string))
			if distributor == nil {
				DebugMustF("重置目标点出错，不存在配送员[%s]", distributor.ID)
				return
			}
			if distributor.CurrentPos.equals(posWanted) {
				DebugMustF("重置目标点出错，不能和当前所在点相同")
				DebugTraceF("%s => %s", distributor.PosString(), posWanted.String())
				DebugPrintList_Info(unit.center.mapData.Points)

				return
			}

			/*
			* 当前位置有两种情况，在节点和在节点之间
			* 如果在节点上，判断目标点与当前节点是否在同一条路径上，是则可以设为终点
			* 如果在节点之间，那么终点只可以设为两个节点之一，这里只需要注意方向即可
			 */
			if distributor.CurrentPos.equals(distributor.StartPos) || distributor.CurrentPos.equals(distributor.DestPos) { //在节点上
				//两点需要在同一条线上
				line := unit.center.mapData.Lines.find(posWanted, distributor.StartPos)
				// line := g_mapData.Lines.find(posWanted, distributor.CurrentPos)
				if line == nil {
					DebugMustF("重置目标点出错，两点不属于同一条路径")
					DebugTraceF("%s => %s", distributor.PosString(), posWanted.String())
					return
				}
				distributor.DestPos = unit.center.mapData.Points.findLngLat(posWanted.Lng, posWanted.Lat)
				distributor.Distance = line.Distance
				distributor.line = line
				unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_reset_destination, distributor)
				DebugInfoF("配送员设置目标点为 %s , 与当前位置 %s 距离为 %f 米", distributor.DestPos.SimpleString(), distributor.CurrentPos.SimpleString(), distributor.Distance)
			} else { //在节点之间
				//如果已经设置了终点，那么新设的点如果和当前的终点相同直接返回
				if distributor.DestPos != nil && distributor.DestPos.equals(posWanted) {
					return
				}
				//如果新设置的点是起点，相当于掉头回去，需要将起始点和终点交换
				if distributor.StartPos.equals(posWanted) {
					p := distributor.StartPos
					distributor.StartPos = distributor.DestPos
					distributor.DestPos = p
					return
				}
				DebugInfoF("没有操作的飞过")
			}
		} else {
			DebugMustF("更改目的地时，客户端数据格式错误: %s", err)
		}
	}
	return f
}
func pro_change_state_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		if list, err := mappedValue(msg.Data.(map[string]interface{})).Getter("DistributorID", "State"); err == nil {
			distributor := unit.distributor
			if distributor.ID != list[0].(string) {
				DebugMustF("重置目标点出错，不存在配送员[%s]", list[0].(string))
				return
			}
			state := int(list[1].(float64))
			if state == 0 {
				distributor.NormalSpeed = 0
			} else {
				distributor.NormalSpeed = defaultSpeed
			}
			unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_change_state, distributor)
			DebugInfoF("配送员 %s 当前速度：%f", unit.distributor.Name, unit.distributor.NormalSpeed)
		} else {
			DebugMustF("更改运动状态时，客户端数据格式错误: %s", err)
		}
	}
	return f
}
func pro_sign_order_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		if list, err := mappedValue(msg.Data.(map[string]interface{})).Getter("DistributorID", "OrderID"); err == nil {
			distributor := unit.distributor
			if distributor.ID != list[0].(string) {
				DebugMustF("订单签收出错，不存在配送员[%s]", list[0].(string))
				return
			}
			orderID := list[1].(string)
			order := unit.center.orders.findByID(orderID)
			if order == nil {
				DebugSysF("不存在订单 %s", orderID)
				return
			}
			if distributor.AcceptedOrders.contains(orderID) == false {
				DebugSysF("订单 %s 必须由本人签收", orderID)
				return
			}
			order.sign()
			unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sign_order, distributor)
			DebugInfoF("配送员 %s 签收了订单 %s", unit.distributor.Name, orderID)
			// DebugPrintList_Info(g_orders)

		} else {
			DebugMustF("签收订单时，客户端数据格式错误: %s", err)
		}
	}
	return f
}

/*
func orderDistributionProcess(msg *MessageWithClient, unit *DistributorProcessUnit) {
	DebugTraceF("执行事件：%d: %s", msg.MessageType, msg.MessageType.name())
	switch msg.MessageType {
	case pro_reset_destination_request:
		m := msg.Data.(map[string]interface{})
		if list, err := mappedValue(m).Getter("PositionID", "DistributorID"); err == nil {
			posWanted := g_mapData.Points.findByID(int(list[0].(float64)))
			if posWanted == nil {
				DebugMustF("重置目标点出错，不存在编号为 %d 的节点", int(list[0].(float64)))
				return
			}
			distributor := g_distributors.find(list[1].(string))
			if distributor == nil {
				DebugMustF("重置目标点出错，不存在配送员[%s]", distributor.ID)
				return
			}
			if distributor.CurrentPos.equals(posWanted) {
				DebugMustF("重置目标点出错，不能和当前所在点相同")
				DebugTraceF("%s => %s", distributor.PosString(), posWanted.String())
				DebugPrintList_Info(g_mapData.Points)

				return
			}


			if distributor.CurrentPos.equals(distributor.StartPos) || distributor.CurrentPos.equals(distributor.DestPos) { //在节点上
				//两点需要在同一条线上
				line := g_mapData.Lines.find(posWanted, distributor.StartPos)
				// line := g_mapData.Lines.find(posWanted, distributor.CurrentPos)
				if line == nil {
					DebugMustF("重置目标点出错，两点不属于同一条路径")
					DebugTraceF("%s => %s", distributor.PosString(), posWanted.String())
					return
				}
				distributor.DestPos = g_mapData.Points.findLngLat(posWanted.Lng, posWanted.Lat)
				distributor.Distance = line.Distance
				distributor.line = line
				g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_reset_destination, distributor)
				DebugInfoF("配送员设置目标点为 %s , 与当前位置 %s 距离为 %f 米", distributor.DestPos.SimpleString(), distributor.CurrentPos.SimpleString(), distributor.Distance)
			} else { //在节点之间
				//如果已经设置了终点，那么新设的点如果和当前的终点相同直接返回
				if distributor.DestPos != nil && distributor.DestPos.equals(posWanted) {
					return
				}
				//如果新设置的点是起点，相当于掉头回去，需要将起始点和终点交换
				if distributor.StartPos.equals(posWanted) {
					p := distributor.StartPos
					distributor.StartPos = distributor.DestPos
					distributor.DestPos = p
					return
				}
				DebugInfoF("没有操作的飞过")
			}
		} else {
			DebugMustF("更改目的地时，客户端数据格式错误: %s", err)
		}
	case pro_change_state_request:
		if list, err := mappedValue(msg.Data.(map[string]interface{})).Getter("DistributorID", "State"); err == nil {
			distributor := unit.distributor
			if distributor.ID != list[0].(string) {
				DebugMustF("重置目标点出错，不存在配送员[%s]", list[0].(string))
				return
			}
			state := int(list[1].(float64))
			if state == 0 {
				distributor.NormalSpeed = 0
			} else {
				distributor.NormalSpeed = defaultSpeed
			}
			g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_change_state, distributor)
			DebugInfoF("配送员 %s 当前速度：%f", unit.distributor.Name, unit.distributor.NormalSpeed)
		} else {
			DebugMustF("更改运动状态时，客户端数据格式错误: %s", err)
		}
	case pro_sign_order_request:
		if list, err := mappedValue(msg.Data.(map[string]interface{})).Getter("DistributorID", "OrderID"); err == nil {
			distributor := unit.distributor
			if distributor.ID != list[0].(string) {
				DebugMustF("订单签收出错，不存在配送员[%s]", list[0].(string))
				return
			}
			orderID := list[1].(string)
			order := g_orders.findByID(orderID)
			if order == nil {
				DebugSysF("不存在订单 %s", orderID)
				return
			}
			if distributor.AcceptedOrders.contains(orderID) == false {
				DebugSysF("订单 %s 必须由本人签收", orderID)
				return
			}
			order.sign()
			g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_sign_order, distributor)
			DebugInfoF("配送员 %s 签收了订单 %s", unit.distributor.Name, orderID)
			// DebugPrintList_Info(g_orders)

		} else {
			DebugMustF("签收订单时，客户端数据格式错误: %s", err)
		}
	}
}
*/

func getOrderSelectProposal(distributors DistributorList, orders OrderList) (list OrderList, err error) {
	proposals, err := createDistributionProposal(orders.Filter(newOrderDistributeFilter(false)), distributors)
	if err != nil {
		// DebugMustF("%s", err.Error())
		return nil, err
	}
	// else {
	// 	// g_room_distributor.sendMsgToSpecialSubscriber(id, pro_2c_order_distribution_proposal, proposals)
	// 	g_room_distributor.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
	// }
	return proposals, nil
}
func disposeOrderSelectResponse(orderID, distributorID string, distributors DistributorList, orders OrderList) error {
	order := orders.findByID(orderID)
	if order == nil {
		DebugMustF("系统异常，分配不存在的订单：%s", orderID)
		return ERR_cannot_find_order
	}
	if order.isDistributed() == true {
		DebugInfoF("订单[%s]已经被分配", orderID)
		return ERR_order_already_selected
	}

	distributor := distributors.find(distributorID) //首先确保配送员满足订单分配条件，当前条件是已分配的订单未达到最大可接收数量
	if distributor == nil {
		DebugInfoF("没有找到配送员[%s]的信息", distributorID)
		return ERR_no_such_distributor
	}
	if distributor.full() {
		DebugInfoF("配送员[%s]已经满载", distributorID)
		return ERR_distributor_full
	}

	//确定结果
	order.distribute()
	distributor.acceptOrder(order)
	DebugTraceF("未分配订单减少到 %d 个", len(orders.Filter(newOrderDistributeFilter(false))))
	return nil
	// //将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
	// distributionResult := NewOrderDistribution(order.ID, distributor.ID)
	// triggerSysEvent(NewSysEvent(sys_event_give_order_select_result, distributionResult))
	// // g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_order_distribute_result, distributionResult)

	// // msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", od.OrderID, distributor.Name)
	// // g_room_distributor.broadcastMsgToSubscribers(pro_2c_message_broadcast, msg)

	// if distributor.full() == true { //配送员的订单满载了
	// 	g_room_distributor.broadcastMsgToSubscribers(pro_2c_message_broadcast, fmt.Sprintf("配送员 %s 订单满载", distributor.Name))
	// 	// g_room_distributor.broadcastMsgToSubscribers(pro_2c_distribution_prepared, distributor.ID)
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
