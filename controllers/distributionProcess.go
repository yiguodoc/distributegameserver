package controllers

import (
// "errors"
// "fmt"
// "math"
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

			/*
			* 当前位置有两种情况，在节点和在节点之间
			* 如果在节点上，判断目标点与当前节点是否在同一条路径上，是则可以设为终点
			* 如果在节点之间，那么终点只可以设为两个节点之一，这里只需要注意方向即可
			 */
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
