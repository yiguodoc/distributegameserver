package controllers

import (
	// "errors"
	// "fmt"
	"math"
	// "strconv"
	// "time"
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

func pro_game_time_tick_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		distributor := unit.distributor
		distributor.TimeElapse = unit.center.TimeElapse
		// if distributor.TimeElapse < unit.center.GameTimeMaxLength { //如果时间超过了最长设定时间，此时，客户端应该发起结束游戏的提示

		// 	//运行时间增加
		// 	distributor.TimeElapse++
		// 	// DebugInfoF("运行时间+1 => %d", distributor.TimeElapse)
		unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sys_time_elapse, distributor.TimeElapse)
		// } else {
		// 	unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_end_game, distributor)
		// 	return
		// }
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
					if math.Abs(lng) <= math.Abs(lngPerFrame) && math.Abs(lat) <= math.Abs(latPerFrame) {
						distributor.CurrentPos.addLngLat(lng, lat)
						//已经到达目标点，运动停止
						// distributor.StartPos.setLngLat(distributor.DestPos.Lng, distributor.DestPos.Lat) //
						distributor.StartPos = unit.center.mapData.Points.findOne(func(pos *Position) bool {
							return pos.Lng == distributor.DestPos.Lng && pos.Lat == distributor.DestPos.Lat
						})
						// distributor.StartPos = unit.center.mapData.Points.findLngLat(distributor.DestPos.Lng, distributor.DestPos.Lat)
						distributor.DestPos = nil
						line.removeDistributor(distributor.ID)
						distributor.line = nil
						// unit.center.wsRoom.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_move_to_new_position, distributor) //通知客户端移动到新坐标
						unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_reach_route_node, distributor) //通知客户端移动到新坐标
						DebugInfoF("配送员已经行驶到目标点 %s", distributor)
						DebugTraceF("配送员实时位置：%s", distributor.PosString())
						//配送员从路上转移到节点
						distributor.CurrentSpeed = distributor.NormalSpeed
						unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_speed_change, distributor)
						unit.center.Process(NewMessageWithClient(pro_move_from_route_to_node, distributor.ID, line))
					} else {
						just_move_to_route := false //测算一下是否是从节点上路的第一步
						if distributor.CurrentPos.equals(distributor.StartPos) {
							just_move_to_route = true
							DebugInfoF("配送员 %s 上路了", distributor.Name)
						}
						distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
						DebugTraceF("配送员实时位置：%s", distributor.PosString())
						if just_move_to_route {
							//配送员从节点到路上
							line.addDistributor(distributor)
							unit.center.Process(NewMessageWithClient(pro_move_from_node_to_route, distributor.ID, line))
							unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_move_from_node, distributor) //通知客户端移动到新坐标
						} else {
							unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_move_to_new_position, distributor) //通知客户端移动到新坐标
						}
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
			posWanted := unit.center.mapData.Points.findOne(func(pos *Position) bool {
				return pos.ID == int(list[0].(float64))
			})
			// posWanted := unit.center.mapData.Points.findByID(int(list[0].(float64)))
			if posWanted == nil {
				DebugMustF("重置目标点出错，不存在编号为 %d 的节点", int(list[0].(float64)))
				return
			}
			distributor := unit.center.distributors.findOne(func(d *Distributor) bool { return d.ID == list[1].(string) })
			// distributor := unit.center.distributors.find(list[1].(string))
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
			if (PositionList{distributor.StartPos, distributor.DestPos}).contains(func(pos *Position) bool { return pos.equals(distributor.CurrentPos) }) {
				// if distributor.CurrentPos.equals(distributor.StartPos) || distributor.CurrentPos.equals(distributor.DestPos) { //在节点上
				//两点需要在同一条线上
				line := unit.center.mapData.Lines.find(posWanted, distributor.StartPos)
				// line := g_mapData.Lines.find(posWanted, distributor.CurrentPos)
				if line == nil {
					DebugMustF("重置目标点出错，两点不属于同一条路径")
					DebugTraceF("%s => %s", distributor.PosString(), posWanted.String())
					return
				}
				distributor.DestPos = unit.center.mapData.Points.findOne(func(pos *Position) bool {
					// return pos.Lng == posWanted.Lng && pos.Lat == posWanted.Lat
					return pos.equals(posWanted)
				})
				// distributor.DestPos = unit.center.mapData.Points.findLngLat(posWanted.Lng, posWanted.Lat)
				distributor.Distance = line.Distance
				distributor.line = line
				unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_reset_destination, distributor)
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
					unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_reset_destination, distributor)
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
			unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_change_state, distributor)
			DebugInfoF("配送员 %s 当前速度：%f", unit.distributor.Name, unit.distributor.NormalSpeed)
		} else {
			DebugMustF("更改运动状态时，客户端数据格式错误: %s", err)
		}
	}
	return f
}
func pro_distributor_info_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		distributor := unit.distributor
		unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_distributor_info, distributor)
	}
	return f
}

//签收订单，注意签收的位置必须是订单所在的位置
func pro_sign_order_request_handlerGenerator(o interface{}) MessageWithClientHandler {
	unit := o.(*DistributorProcessUnit)
	f := func(msg *MessageWithClient) {
		if list, err := mappedValue(msg.Data.(map[string]interface{})).Getter("DistributorID", "OrderID"); err == nil {
			distributor := unit.distributor
			if distributor.ID != list[0].(string) {
				DebugMustF("订单签收出错，不存在配送员[%s]", list[0].(string))
				unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sign_order, nil)
				return
			}
			orderID := list[1].(string)
			order := unit.center.orders.findOne(func(o interface{}) bool { return o.(*Order).ID == orderID })
			// order := unit.center.orders.findByID(orderID)
			if order == nil {
				DebugSysF("不存在订单 %s", orderID)
				unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sign_order, nil)
				return
			}
			if order.isSigned() {
				DebugSysF("订单已经签收过 %s", order)
				return
			}
			if distributor.AcceptedOrders.contains(func(o interface{}) bool { return o.(*Order).ID == orderID }) == false {
				DebugSysF("订单 %s 必须由本人签收", orderID)
				unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sign_order, nil)
				return
			}
			order.sign(distributor.TimeElapse)
			DebugInfoF("签收订单 %s , 时间 %d", order.ID, order.SignTime)
			unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_sign_order, distributor)
			DebugInfoF("配送员 %s 签收了订单 %s", unit.distributor.Name, orderID)
			distributor.Score++
			// DebugPrintList_Info(g_orders)
			if unit.distributor.AcceptedOrders.all(func(o interface{}) bool { return o.(*Order).Signed == true }) {
				// unit.distributor.setCheckPoint(checkpoint_flag_order_distribute_over)
				unit.center.sendMsgToSpecialSubscriber(distributor.ID, pro_2c_all_order_signed, distributor)
			}

		} else {
			DebugMustF("签收订单时，客户端数据格式错误: %s", err)
		}
	}
	return f
}
