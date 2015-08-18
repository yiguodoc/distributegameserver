package controllers

import (
	// "github.com/gorilla/websocket"
	"encoding/json"
	// "time"
	// "strings"
	// "fmt"
)

/*
事件触发和事件处理作为隔离的两部分
事件处理负责处理单一的事件
事件触发上有单独
*/
var (
	// g_mapData          *MapData
	// g_room_distributor *WsRoom = NewRoom()
	g_UnitCenter  *DistributorProcessUnitCenter
	g_room_viewer *WsRoom //= NewRoom(eventReceiver)
	// g_orders                   = OrderList{} //所有的订单

	// g_distributors DistributorList = DistributorList{ //配送员列表
	// 	NewDistributor("d01", "张军", 2, color_orange),
	// 	NewDistributor("d02", "刘晓莉", 2, color_red),
	// 	// NewDistributor("d03", "桑鸿庆", 3, color_purple),
	// }
	// g_sysEventSubscribeList         = SysEventSubscribeList{} //对各种事件的订阅函数列表
	// g_chanEvents                    = make(chan *SysEvent)    //系统事件中转站
	// g_chanEventPkgOver                                    = make(chan int64)        //接收事件包的id，表明该事件包执行完毕，每个事件包监测到自己没有后续事件执行时，发送自己的id
	// g_distributorUnits                                    = make(DistributorProcessUnitList)
)

func init() {
	// sysEventCodeCheck()
	clientMessageTypeCodeCheck()
	// initEventSubscribe()
	// initSysEventRoute() //系统事件的分发处理

	// g_room_viewer.init()

	distributors := DistributorList{ //配送员列表
		NewDistributor("d01", "张军", 2, color_orange),
		NewDistributor("d02", "刘晓莉", 2, color_red),
		// NewDistributor("d03", "桑鸿庆", 3, color_purple),
	}
	mapData := loadMapData()
	//加载地图数据
	orders := mapData.Points.createSimulatedOrders(generateOrderID) //生成模拟订单
	DebugPrintList_Info(orders)
	// orders := OrderList{} //所有的订单

	room := NewRoom()
	room.init()
	room.addEventSubscriber(distributorRoomEventHandler,
		WsRoomEventCode_Online, WsRoomEventCode_Offline, WsRoomEventCode_Other)

	g_UnitCenter = NewDistributorProcessUnitCenter(room, distributors, orders, mapData, handler_map)
	g_UnitCenter.start()

	//测试用
	//将订单分配给配送员
	// g_distributors[0].AcceptedOrders = g_orders[0:]
	// g_distributors[0].setCheckPoint(checkpoint_flag_order_distribute)
	// g_distributors[1].AcceptedOrders = g_orders[2:]
	// g_distributors[1].setCheckPoint(checkpoint_flag_order_distribute)
	// g_ordersDistributed = g_ordersUndistributed[:]
	// g_ordersUndistributed = g_ordersUndistributed[0:0]
	//--------------------------------------------------------------------------
}

func distributorRoomEventHandler(code, data interface{}) {
	c := WsRoomEventCode(code.(int))

	switch c {
	case WsRoomEventCode_Online:
		msg := data.(*MessageWithClient)
		distributor := g_UnitCenter.distributors.find(msg.TargetID)
		if distributor == nil {
			DebugSysF("未查找到配送员 %s", msg.Data.(string))
			return
		}
		// distributor := data.(*Distributor)
		// unit := g_UnitCenter.newUnit(distributor)
		// unit.addProcessor(handler_map)
		// msg.TargetID = distributor.ID
		g_UnitCenter.Process(msg)
		// g_UnitCenter.singleUnitprocess(distributor.ID, msg)
		DebugTraceF("配送员上线 ：%s", distributor.String())

	case WsRoomEventCode_Offline:
		msg := data.(*MessageWithClient)
		// id := msg.Data.(string)
		// msg := NewMessageWithClient(pro_off_line, id, nil)
		// subscriber := data.(Subscriber)
		// g_UnitCenter.removeUnit(id)
		g_UnitCenter.Process(msg)
		// DebugTraceF("配送员下线 ：%s", id)

	case WsRoomEventCode_Other:
		roommsg := data.(*roomMessage)
		var msg MessageWithClient
		err := json.Unmarshal(roommsg.content.([]byte), &msg)
		if err != nil {
			DebugSysF("解析数据出错：%s", err)
			return
		}
		msg.TargetID = roommsg.targetID
		g_UnitCenter.Process(&msg)
		// g_UnitCenter.singleUnitprocess(roommsg.targetID, &msg)
	}
}

func initEventSubscribe() {
	// codes := getSysEventCodeList()
	// for _, v := range codes {
	// 	g_sysEventSubscribeList = append(g_sysEventSubscribeList, NewSysEventSubscribe(v))
	// }
	// //公共事件
	// g_sysEventSubscribeList.addEventSubscriber(onUserOnlineChange, sys_event_user_online)
	// g_sysEventSubscribeList.addEventSubscriber(onUserOfflineChange, sys_event_user_offline)
	// g_sysEventSubscribeList.addEventSubscriber(onCountDown, sys_event_count_down)
	// g_sysEventSubscribeList.addEventSubscriber(onCountDownSilent, sys_event_count_down_silent)
	// g_sysEventSubscribeList.addEventSubscriber(onMsgBroadcast, sys_event_message_broadcast)

	// //订单分发
	// g_sysEventSubscribeList.addEventSubscriber(NewProcessNode(orderSelectionProcess).init().NodeEventProcessor(),
	// 	sys_event_start_select_order,
	// 	sys_event_give_out_order,
	// 	sys_event_give_order_select_result,
	// 	sys_event_order_select_additional_msg,
	// 	sys_event_select_order_prepared,
	// 	sys_event_order_select_response)

	// g_sysEventSubscribeList.addEventSubscriber(NewProcessNode(orderDistributionProcess).init().NodeEventProcessor(),
	// 	sys_event_start_order_distribution,
	// 	sys_event_distribution_prepared,
	// 	sys_event_order_distribute_additional_msg,
	// 	sys_event_reset_destination_request,
	// 	sys_event_change_state_request,
	// 	sys_event_order_distribute_result)

	// if ls := g_sysEventSubscribeList.checkEventWithNoSubscriber(); len(ls) > 0 {
	// 	DebugMust("有定义的事件没有发挥作用")
	// 	for _, code := range ls {
	// 		DebugTraceF("事件定义：%3d : %s", code, code.name())
	// 	}
	// }
}

// func onOrderDistributeResultToViewer(event *SysEvent) {
// 	g_room_viewer.broadcastMsgToSubscribers(pro_order_distribute_result, event.data)
// }

// func eventReceiver(id string, data interface{}) {
// 	// func eventReceiver(code sysEventCode, para interface{}) {
// 	// func eventReceiver(eventName string, para interface{}) {
// 	// g_chanEvents <- NewSysEvent(code, para)
// 	// g_chanEvents <- NewSysEvent(eventName, para)
// 	g_distributorUnits.process(id, data)
// }

//配送员上线，需要创建处理单元，并根据进度初始化
// func onDistributorOnlineChange(data interface{}) {
// 	// func onDistributorOnlineChange(event *SysEvent) {
// 	distributor := data.(*Distributor)
// 	// distributor := event.data.(*Distributor)
// 	g_distributorUnits.newUnit(distributor)
// 	g_distributorUnits.process(distributor.ID, nil)
// 	DebugTraceF("配送员上线 ：%s", distributor.String())
// 	//用户上线，应该根据其上线前状态，（可能第一次上线，也可能中间掉线），推送其必要的信息
// 	//如果在分配订单中，应该推送给其正在选择的订单
// 	// switch distributor.CheckPoint {
// 	// case checkpoint_flag_origin:

// 	// case checkpoint_flag_order_select:
// 	// 	triggerSysEvent(NewSysEvent(sys_event_order_select_additional_msg, distributor.GetID()))

// 	// case checkpoint_flag_order_distribute:
// 	// 	triggerSysEvent(NewSysEvent(sys_event_start_order_distribution, distributor))

// 	// }
// }

// func onDistributorOfflineChange(data interface{}) {
// 	// func onDistributorOfflineChange(event *SysEvent) {
// 	subscriber := data.(Subscriber)
// 	// subscriber := event.data.(Subscriber)
// 	g_distributorUnits.removeUnit(subscriber.GetID())
// 	DebugTraceF("配送员下线 ：%s", subscriber.String())

// }

// //上线时的事件处理
// func onUserOnlineChange(event *SysEvent) {
// 	// if g_room.subscribers.onlineCountGreaterThan(1) == true {
// 	// 	DebugTraceF("在线人数超过1，可以开始订单分发流程")
// 	// 	triggerSysEvent(NewSysEvent(("start_order_distribution_countdown"), nil))
// 	// }
// 	switch event.data.(type) {
// 	case *Distributor:
// 		distributor := event.data.(*Distributor)
// 		DebugTraceF("配送员上线 ：%s", distributor.String())
// 		//用户上线，应该根据其上线前状态，（可能第一次上线，也可能中间掉线），推送其必要的信息
// 		//如果在分配订单中，应该推送给其正在选择的订单
// 		switch distributor.CheckPoint {
// 		case checkpoint_flag_origin:
// 		case checkpoint_flag_order_select:
// 			// requestOrderSelectAdditionalMsg(distributor.GetID())
// 			triggerSysEvent(NewSysEvent(sys_event_order_select_additional_msg, distributor.GetID()))

// 		case checkpoint_flag_order_distribute:
// 			// triggerSysEvent(NewSysEvent(sys_event_order_distribute_additional_msg, distributor.GetID()))
// 			triggerSysEvent(NewSysEvent(sys_event_start_order_distribution, distributor))

// 		}
// 		g_room_viewer.broadcastMsgToSubscribers(pro_distributor_on_line, event.data)
// 	case *Viewer:
// 	}
// }
// func onUserOfflineChange(event *SysEvent) {
// 	subscriber := event.data.(Subscriber)
// 	DebugTraceF("配送员下线 ：%s", subscriber.String())
// 	g_room_viewer.broadcastMsgToSubscribers(pro_distributor_off_line, subscriber)
// }

//专门针对客户端准备好动作的反应
//只针对第一次准备好的情况
// func onDistributorPrepared(event *SysEvent) {
// 	DebugTraceF("配送员[%s]准备好订单的分发了", event.data)
// 	g_distributors.preparedForOrderSelect(event.data.(string))
// 	if g_distributors.allPreparedForOrderSelect() == true {
// 		DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
// 		// messageEvent := NewSysEvent(("message_broadcast"), "订单分发即将开始")
// 		// countDownEvent := NewCountDownEvent(3)
// 		startOrderDistributionEvent := NewSysEvent(sys_event_start_select_order, nil)
// 		triggerSysEvent(startOrderDistributionEvent)
// 		// pkg := NewSysEventPkg(messageEvent, countDownEvent)
// 		// addEventPkg(pkg)
// 	}
// }

//向配送员发布订单分配结果
// func onOrderDistributeResult(event *SysEvent) {
// 	g_room_distributor.broadcastMsgToSubscribers(pro_order_distribute_result, event.data)
// }

//通知客户端，允许配送员发送开始抢订单
// func onBeginSelectOrder(event *SysEvent) {
// 	g_room_distributor.broadcastMsgToSubscribers(pro_begin_select_order, nil)
// }

//向客户端广播消息
// func onMsgBroadcast(event *SysEvent) {
// 	g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, event.data)
// }

// //持续向客户端发送倒计时消息
// func onCountDown(event *SysEvent) {
// 	if event.data == nil {
// 		DebugMustF("倒计时参数错误")
// 		return
// 	}
// 	count := event.data.(int)
// 	if count < 0 {
// 		DebugInfoF("倒计时参数太小：%d", count)
// 	}
// 	timer := time.Tick(1 * time.Second)
// 	DebugInfo("start timer...")
// 	for {
// 		<-timer
// 		DebugTraceF("timer count : %d", count)
// 		if count <= 0 {
// 			break
// 		}
// 		g_room_distributor.broadcastMsgToSubscribers(pro_timer_count_down, count)
// 		count--
// 	}
// 	if event.nextEvent != nil {
// 		triggerSysEvent(event.nextEvent)
// 	}
// }

// //持续向客户端发送倒计时消息
// func onCountDownSilent(event *SysEvent) {
// 	if event.data == nil {
// 		DebugMustF("倒计时参数错误")
// 		return
// 	}
// 	count := event.data.(int)
// 	if count < 0 {
// 		DebugInfoF("倒计时参数太小：%d", count)
// 	}
// 	timer := time.Tick(1 * time.Second)
// 	DebugInfo("start timer...")
// 	for {
// 		<-timer
// 		DebugTraceF("timer count : %d", count)
// 		if count <= 0 {
// 			break
// 		}
// 		// g_room_distributor.broadcastMsgToSubscribers(pro_timer_count_down, count)
// 		count--
// 	}
// 	if event.nextEvent != nil {
// 		triggerSysEvent(event.nextEvent)
// 	}
// }

// //系统事件监听
// func initSysEventRoute() {
// 	// isTherePkgRunning := false //有没有正在处理的事件包，有的话不开启新的
// 	go func() {
// 		for {
// 			select {
// 			case event := <-g_chanEvents:
// 				DebugTraceF("接收到系统事件：%d: %s", event.eventCode, event.eventCode.name())
// 				// g_sysEventSubscribeList.notifyEventSubscribers(event)
// 			}

// 		}
// 	}()
// }

// //触发系统事件
// func triggerSysEvent(event *SysEvent) {
// 	// DebugTraceF("触发系统事件：%d  %s", event.eventCode, event.eventCode.name())
// 	g_chanEvents <- event
// }
