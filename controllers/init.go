package controllers

import (
	// "container/list"
	// "github.com/astaxie/beego"
	// "github.com/gorilla/websocket"
	"time"
	// "encoding/json"
	// "strings"
	// "fmt"
)

/*
事件触发和事件处理作为隔离的两部分
事件处理负责处理单一的事件
事件触发上有单独
*/
var (
	g_mapData MapData

	g_distributors DistributorList = DistributorList{ //配送员列表
		NewDistributor("d01", "张军", 3, color_orange),
		NewDistributor("d02", "刘晓莉", 3, color_red),
		// NewDistributor("d03", "桑鸿庆", 3, color_purple),
	}
	// g_room.subscribers           = SubscriberList{}        //配送人员列表，在线或者不在线
	// g_chanEventStore        = make(chan *SysEventPkg) //接收事件包
	// g_eventpkgStore         = SysEventPkgList{}
	g_chanEvents            = make(chan *SysEvent)    //系统事件中转站
	g_chanEventPkgOver      = make(chan int64)        //接收事件包的id，表明该事件包执行完毕，每个事件包监测到自己没有后续事件执行时，发送自己的id
	g_sysEventSubscribeList = SysEventSubscribeList{} //对各种事件的订阅函数列表
	g_ordersDistributed     = OrderList{}             //已经分配的订单
	g_ordersUndistributed   = OrderList{              //尚未分配的订单
	// NewOrder("100001", NewPosition("北京市通州区物资学院", "116.643936", "39.936224", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100002", NewPosition("北京市朝阳区常营中路万象新天四区", "116.599632", "39.933292", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100003", NewPosition("北京市朝阳区四季星河中街星河湾", "116.527732", "39.933665", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100004", NewPosition("北京市朝阳区三里屯路22号", "116.462012", "39.941894", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100005", NewPosition("北京市东城区海运仓胡同2号", "116.436428", "39.943001", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100006", NewPosition("北京市东城区交道口北二条36号", "116.419288", "39.94903", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100007", NewPosition("北京市西城区四环胡同13号", "116.38465", "39.948035", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100008", NewPosition("北京市西城区北展北街17号", "116.355652", "39.947758", POSITION_TYPE_ORDER_ROUTE)),
	// NewOrder("100009", NewPosition("北京市海淀区蓝靛厂南路59号", "116.299742", "39.943001", POSITION_TYPE_ORDER_ROUTE)),
	}
	g_room_distributor *WsRoom = NewRoom(eventReceiver)
	g_room_viewer      *WsRoom = NewRoom(eventReceiver)
)

func init() {
	sysEventCodeCheck()
	codes := getSysEventCodeList()
	for _, v := range codes {
		g_sysEventSubscribeList = append(g_sysEventSubscribeList, NewSysEventSubscribe(v))
	}
	g_sysEventSubscribeList.addEventSubscriber(sys_event_user_online, "onUserOnlineChange", onUserOnlineChange)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_user_offline, "onUserOfflineChange", onUserOfflineChange)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_distributor_prepared, "onDistributorPrepared", onDistributorPrepared)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_count_down, "onCountDown", onCountDown)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_count_down_silent, "onCountDownSilent", onCountDownSilent)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_message_broadcast, "onMsgBroadcast", onMsgBroadcast)

	//订单分发
	g_sysEventSubscribeList.addEventSubscriber(sys_event_start_order_selection, "onStartOrderSelection", orderSelectionProcess)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_begin_select_order, "onBeginSelectOrder", orderSelectionProcess)
	g_sysEventSubscribeList.addEventSubscriber(sys_event_order_select_response, "onOrderSelectionResponse", orderSelectionProcess)

	// g_sysEventSubscribeList.addEventSubscriber(sys_event_order_select_result, "onOrderDistributeResultToViewer", onOrderDistributeResultToViewer)
	//
	// g_sysEventSubscribeList.addEventSubscriber(("user_online"), "onUserOnlineChangeForViewer", onUserOnlineChangeForViewer)
	// g_sysEventSubscribeList.addEventSubscriber(("user_offline"), "onUserOfflineChangeForViewer", onUserOfflineChangeForViewer)
	// g_sysEventSubscribeList.addEventSubscriber(("start_order_distribution_countdown"), "onStartOrderDistributionCountdown", onStartOrderDistributionCountdown)
	initSysEventChannel()
	// initUserList(g_distributors)

	g_room_viewer.init()
	g_room_distributor.init()

	loadMapData()                                                                   //加载地图数据
	g_ordersUndistributed = g_mapData.Points.createSimulatedOrders(generateOrderID) //生成模拟订单
	DebugPrintList_Info(g_ordersUndistributed)

	//测试用
	//将订单分配给配送员
	// g_distributors[0].AcceptedOrders = g_ordersUndistributed[0:3]
	// g_distributors[0].CheckPoint = checkpoint_flag_order_distribute
	// g_distributors[1].AcceptedOrders = g_ordersUndistributed[3:]
	// g_distributors[1].CheckPoint = checkpoint_flag_order_distribute
	// g_ordersDistributed = g_ordersUndistributed[:]
	// g_ordersUndistributed = g_ordersUndistributed[0:0]
	//--------------------------------------------------------------------------
}
func onOrderDistributeResultToViewer(event *SysEvent) {
	g_room_viewer.broadcastMsgToSubscribers(pro_order_distribute_result, event.data)
}

func eventReceiver(code sysEventCode, para interface{}) {
	// func eventReceiver(eventName string, para interface{}) {
	g_chanEvents <- NewSysEvent(code, para)
	// g_chanEvents <- NewSysEvent(eventName, para)
}
func onStartOrderDistributionCountdown(event *SysEvent) {
	go func() {
		triggerSysEvent(NewMessageBroadcastEvent("订单分发即将开始"))
		time.Sleep(1 * time.Second)
		eventCountDown := NewCountDownEvent(3)
		eventCountDown.setNextEvent(NewSysEvent(sys_event_start_order_selection, nil))
		triggerSysEvent(eventCountDown)
	}()
}

// func onUserOfflineChangeForViewer(event *SysEvent) {
// 	g_room.broadcastMsgToViewer(pro_distributor_off_line, event.data)
// }

// //上线时的事件处理，针对观察者
// func onUserOnlineChangeForViewer(event *SysEvent) {
// 	g_room.broadcastMsgToViewer(pro_distributor_on_line, event.data)
// }

//上线时的事件处理
func onUserOnlineChange(event *SysEvent) {
	// if g_room.subscribers.onlineCountGreaterThan(1) == true {
	// 	DebugTraceF("在线人数超过1，可以开始订单分发流程")
	// 	triggerSysEvent(NewSysEvent(("start_order_distribution_countdown"), nil))
	// }
	switch event.data.(type) {
	case *Distributor:
		distributor := event.data.(*Distributor)
		DebugTraceF("配送员上线 ：%s", distributor.String())
		//用户上线，应该根据其上线前状态，（可能第一次上线，也可能中间掉线），推送其必要的信息
		//如果在分配订单中，应该推送给其正在选择的订单
		switch distributor.CheckPoint {
		case checkpoint_flag_origin:
		case checkpoint_flag_order_select:
			requestOrderSelectAdditionalMsg(distributor.GetID())
		case checkpoint_flag_order_distribute:
		}
		g_room_viewer.broadcastMsgToSubscribers(pro_distributor_on_line, event.data)
	case *Viewer:
	}

}
func onUserOfflineChange(event *SysEvent) {
	subscriber := event.data.(Subscriber)
	DebugTraceF("配送员下线 ：%s", subscriber.String())
	g_room_viewer.broadcastMsgToSubscribers(pro_distributor_off_line, subscriber)
}

//专门针对客户端准备好动作的反应
//只针对第一次准备好的情况
func onDistributorPrepared(event *SysEvent) {
	DebugTraceF("配送员[%s]准备好订单的分发了", event.data)
	g_distributors.preparedForOrderSelect(event.data.(string))
	if g_distributors.allPreparedForOrderSelect() == true {
		DebugInfoF("所有配送员准备完毕，可以开始订单分发了")
		// messageEvent := NewSysEvent(("message_broadcast"), "订单分发即将开始")
		// countDownEvent := NewCountDownEvent(3)
		startOrderDistributionEvent := NewSysEvent(sys_event_start_order_selection, nil)
		triggerSysEvent(startOrderDistributionEvent)
		// pkg := NewSysEventPkg(messageEvent, countDownEvent)
		// addEventPkg(pkg)
	}
}

//向配送员发布订单分配结果
func onOrderDistributeResult(event *SysEvent) {
	g_room_distributor.broadcastMsgToSubscribers(pro_order_distribute_result, event.data)
}

//通知客户端，允许配送员发送开始抢订单
func onBeginSelectOrder(event *SysEvent) {
	g_room_distributor.broadcastMsgToSubscribers(pro_begin_select_order, nil)
}

//向客户端广播消息
func onMsgBroadcast(event *SysEvent) {
	g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, event.data)
}

//持续向客户端发送倒计时消息
func onCountDown(event *SysEvent) {
	if event.data == nil {
		DebugMustF("倒计时参数错误")
		return
	}
	count := event.data.(int)
	if count < 0 {
		DebugInfoF("倒计时参数太小：%d", count)
	}
	timer := time.Tick(1 * time.Second)
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
	if event.nextEvent != nil {
		triggerSysEvent(event.nextEvent)
	}
}

//持续向客户端发送倒计时消息
func onCountDownSilent(event *SysEvent) {
	if event.data == nil {
		DebugMustF("倒计时参数错误")
		return
	}
	count := event.data.(int)
	if count < 0 {
		DebugInfoF("倒计时参数太小：%d", count)
	}
	timer := time.Tick(1 * time.Second)
	DebugInfo("start timer...")
	for {
		<-timer
		DebugTraceF("timer count : %d", count)
		if count <= 0 {
			break
		}
		// g_room_distributor.broadcastMsgToSubscribers(pro_timer_count_down, count)
		count--
	}
	if event.nextEvent != nil {
		triggerSysEvent(event.nextEvent)
	}
}

// //通过配送员信息初始化参与人员列表
// func initUserList(distributors DistributorList) {
// 	for _, d := range distributors {
// 		g_room.subscribers = append(g_room.subscribers, NewSubscriber(d.ID, subscriber_type_distributor, nil))
// 	}
// 	DebugPrintList_Trace(g_room.subscribers)
// }

//系统事件监听
func initSysEventChannel() {
	// isTherePkgRunning := false //有没有正在处理的事件包，有的话不开启新的
	go func() {
		for {
			select {
			case event := <-g_chanEvents:
				DebugTraceF("接收到系统事件：%d: %s", event.eventCode, event.eventCode.name())
				// switch event.eventCode {
				// case sys_event_user_online:
				// 	DebugInfoF("新用户上线")
				// 	DebugPrintList_Trace(g_room_distributor.subscribers)
				// case getSysEventDefValue("user_offline"):
				// 	DebugInfoF("用户下线")
				// 	DebugPrintList_Trace(g_room_distributor.subscribers)
				// case getSysEventDefValue("count_down"):
				// 	DebugInfo("开始倒计时")
				// case getSysEventDefValue("begin_select_order"):
				// 	DebugInfo("开始选择订单")
				// }
				g_sysEventSubscribeList.notifyEventSubscribers(event)
				// case eventpkg := <-g_chanEventStore: //只负责收集任务
				// 	DebugTraceF("new event pkg %d", eventpkg.id)
				// 	g_eventpkgStore = g_eventpkgStore.add(eventpkg)
				// 	raiseEventPkg()
				// case id := <-g_chanEventPkgOver: //事件包完成后，开启新的事件包
				// 	g_eventpkgStore = g_eventpkgStore.remove(id)
				// 	raiseEventPkg()
			}

		}
	}()
}

// func addEventPkg(pkg *SysEventPkg) {
// 	g_chanEventStore <- pkg
// }

// //开启新的事件包
// func raiseEventPkg() {
// 	if len(g_eventpkgStore) > 0 {
// 		pkg0 := g_eventpkgStore[0]
// 		if len(pkg0.events) > 0 {
// 			// isTherePkgRunning = true
// 			go triggerSysEvent(pkg0.events[0])
// 		}
// 	}
// }

//触发系统事件
func triggerSysEvent(event *SysEvent) {
	DebugTraceF("触发系统事件：%d  %s", event.eventCode, event.eventCode.name())
	g_chanEvents <- event
	// if len(args) > 0 {
	// 	g_chanEvents <- NewSysEvent(event, args[0])
	// } else {
	// 	g_chanEvents <- NewSysEvent(event, nil)
	// }
}
