package controllers

import (
	"errors"
	"fmt"
	// "math"
	"github.com/gorilla/websocket"
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

//整个游戏结束
func pro_game_timeout_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		gu.gameStarted = false
		gu.Distributors.forEach(func(d *Distributor) {
			d.caculateScore()
		})
		gu.Distributors.Rank().forEach(func(d *Distributor) {
			gu.sendMsgToSpecialSubscriber(d, pro_2c_rank_change, d)
		})
		DebugPrintList_Info(gu.Distributors)
		gu.Distributors.filter(func(d *Distributor) bool { return d.whetherHasEndTheGame() == false }).
			forEach(func(d *Distributor) {
			d.setCheckPoint(checkpoint_flag_game_over)
			gu.sendMsgToSpecialSubscriber(d, pro_2c_end_game, d)
		})
	}
	return f
}

//单独申请配送结束
func pro_end_game_request_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		go func() {
			//计算得分
			//签收完一个订单得到该订单对应的分数，没有完成的订单减去惩罚分数
			msg.Target.caculateScore()
			//计算排名
			gu.Distributors.Rank()
			DebugPrintList_Info(gu.Distributors)
			msg.Target.setCheckPoint(checkpoint_flag_game_over)
			gu.sendMsgToSpecialSubscriber(msg.Target, pro_2c_end_game, msg.Target)
			gu.Distributors.forEach(func(d *Distributor) {
				if d.ID != msg.Target.ID {
					gu.sendMsgToSpecialSubscriber(d, pro_2c_rank_change, d)
				}
			})
		}()
	}
	return f
}

func pro_move_from_node_to_route_handlerGenerator(center *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		go func() {
			line := msg.Data.(*Line)
			if line.DistributorsCount() >= 2 {
				line.busy()
				DebugInfoF("line BUSY %s", line)
				line.DistributorsOn.forEach(func(d *Distributor) {
					d.CurrentSpeed = d.NormalSpeed / 2
					center.sendMsgToSpecialSubscriber(d, pro_2c_speed_change, d)
				})
			}
		}()
	}
	return f
}

func pro_move_from_route_to_node_handlerGenerator(center *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		go func() {
			line := msg.Data.(*Line)
			if line.DistributorsCount() < 2 { //这里只是可能会变得不堵
				line.nobusy()
				DebugInfoF("line NOBUSY %s ", line)
				line.DistributorsOn.forEach(func(d *Distributor) {
					d.CurrentSpeed = d.NormalSpeed
					center.sendMsgToSpecialSubscriber(d, pro_2c_speed_change, d)
				})
			}
		}()
	}
	return f
}

//订单选择的请求，可以对请求作出一些限制，例如当前的位置必须处于仓库等等
func pro_order_select_response_handlerGenerator(center *GameUnit) MessageWithClientHandler {
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
		distributor := msg.Target
		//系统设定，配送员需要在仓库位置才能选择订单
		warehousePointlist := center.mapData.Points.filter(func(pos *Position) bool { return pos.PointType == POSITION_TYPE_WAREHOUSE })
		if warehousePointlist.contains(func(pos *Position) bool { return distributor.CurrentPos.equals(pos) }) == false {
			DebugInfoF("配送员 %s 所处位置无法选择订单", distributor.Name)
			center.sendMsgToSpecialSubscriber(distributor, pro_2c_order_select_result, nil, "所处位置无法选择订单", strconv.Itoa(distributor.TimeElapse))
			return
		}

		if err := disposeOrderSelectResponse(orderID, distributor, center.Distributors, center.orders); err != nil {
			DebugInfoF("处理订单分配时的信息提示：%s", err)
			center.sendMsgToSpecialSubscriber(distributor, pro_2c_order_select_result, nil, err.Error(), strconv.Itoa(distributor.TimeElapse))
			return
		}
		//将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
		center.sendMsgToSpecialSubscriber(distributor, pro_2c_order_select_result, distributor)

		log := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", orderID, distributor.Name)
		// center.wsRoom.broadcastMsgToSubscribers(pro_2c_message_broadcast, msg)
		DebugInfoF(log)

		sendOrderProposal(center)
	}
	return f
}

func pro_game_start_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		msgList := []string{"配送员全部准备完毕", "请前往配送中心选择订单", "3", "2", "1"}
		for _, msg := range msgList {
			gu.broadcastMsgToSubscribers(pro_2c_message_broadcast_before_game_start, msg)
			time.Sleep(1 * time.Second)
		}
		gu.Distributors.forEach(func(d *Distributor) {
			d.setCheckPoint(checkpoint_flag_game_started)
			gu.sendMsgToSpecialSubscriber(d, pro_2c_game_start, d)
		})
		sendOrderProposal(gu)
		gu.startGameTiming()
	}
	return f
}
func sendOrderProposal(gu *GameUnit) {
	proposals := getOrderSelectProposal(gu.Distributors, gu.orders)
	gu.broadcastMsgToSubscribers(pro_2c_order_distribution_proposal, proposals)
}

func pro_prepared_for_select_order_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		DebugInfoF("配送员[%s]准备好订单的分发了", msg.Target.Name)

		msg.Target.setCheckPoint(checkpoint_flag_prepared_for_game)
		//提醒尚未准备好进入游戏者
		gu.broadcastMsgToSubscribers(pro_2c_on_line_user_change, gu.Distributors.filter(func(d *Distributor) bool { return d.CheckPoint < checkpoint_flag_prepared_for_game }))

		if gu.Distributors.every(func(d *Distributor) bool { return d.CheckPoint >= checkpoint_flag_prepared_for_game }) {
			DebugInfoF("所有配送员准备完毕，游戏开始")
			gu.Process(NewMessageWithClient(pro_game_start, msg.Target, nil))
		} else {
			DebugInfoF("还有 %d 个配送员未准备完毕", len(gu.Distributors.filter(func(d *Distributor) bool { return d.CheckPoint < checkpoint_flag_prepared_for_game })))
		}
	}
	return f
}
func pro_off_line_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		msg.Target.SetOffline()
		DebugTraceF("%s", msg)
		DebugInfoF("配送员 %s 离线", msg.Target.Name)

		//提醒尚未准备好进入游戏者
		distributorsOfflineAndNotprepared := gu.Distributors.filter(func(d *Distributor) bool {
			return d.CheckPoint < checkpoint_flag_prepared_for_game
		})
		gu.broadcastMsgToSubscribers(pro_2c_on_line_user_change, distributorsOfflineAndNotprepared)
	}
	return f
}
func pro_on_line_handlerGenerator(gu *GameUnit) MessageWithClientHandler {
	f := func(msg *MessageWithClient) {
		go func() {
			DebugTraceF("处理消息 %s", msg)
			msg.Target.SetConn(msg.Data.(*websocket.Conn))
			gu.sendMsgToSpecialSubscriber(msg.Target, pro_2c_map_data, gu.mapData)
			gu.sendMsgToSpecialSubscriber(msg.Target, pro_2c_distributor_info, msg.Target)
			onReconnect(gu, msg.Target)
		}()
	}
	return f
}
func onReconnect(gu *GameUnit, distributor *Distributor) {
	//如果在分配订单中，应该推送给其正在选择的订单
	switch distributor.CheckPoint {
	case checkpoint_flag_origin:
		DebugTraceF("配送员 %s 上线，状态 %d 初始化", distributor.Name, checkpoint_flag_origin)
	case checkpoint_flag_prepared_for_game:
		DebugTraceF("配送员 %s 上线，状态 %d 准备好游戏了", distributor.Name, checkpoint_flag_prepared_for_game)
		//如果之前配送员已经提交准备好的请求，现在是掉线重连状态，那么客户端方面就不会重新请求，因此如果是最后一个
	case checkpoint_flag_game_started:
		DebugTraceF("配送员 %s 上线，状态 %d 游戏进行中", distributor.Name, checkpoint_flag_game_started)
		// broadOrderSelectProposal(gu.distributors, gu.orders)
		proposals := getOrderSelectProposal(gu.Distributors, gu.orders)
		gu.sendMsgToSpecialSubscriber(distributor, pro_2c_order_distribution_proposal, proposals)
	case checkpoint_flag_game_over:
		DebugTraceF("配送员 %s 上线，状态 %d 配送完成", distributor.Name, checkpoint_flag_game_over)
	}
}

func getOrderSelectProposal(distributors DistributorList, orders OrderList) (list OrderList) {
	ordersUndistributed := orders.Filter(func(o *Order) bool { return o.Distributed == false })
	if len(ordersUndistributed) >= 5 {
		list = ordersUndistributed[0:5]
	} else {
		list = ordersUndistributed
	}
	return
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

	//确定结果
	order.distribute(distributor.TimeElapse)
	distributor.acceptOrder(order)
	DebugTraceF("未分配订单减少到 %d 个", len(orders.Filter(func(o *Order) bool { return o.Distributed == false })))
	return nil
}
