package controllers

import (
	// "github.com/ssor/fauxgaux"
	// "github.com/gorilla/websocket"
	"encoding/json"
	// "time"
	// "strings"
	// "fmt"
	// "reflect"
)

var (
	g_UnitCenter       *DistributorProcessUnitCenter
	g_room_viewer      *WsRoom            //= NewRoom(eventReceiver)
	g_distributorStore = DistributorList{ //配送员列表
		NewDistributor("d01", "张军", 2, color_orange),
		NewDistributor("d02", "刘晓莉", 2, color_red),
		// NewDistributor("d03", "桑鸿庆", 3, color_purple),
	}
	g_regions = RegionList{
		NewRegion("1", "255,128,128", 39.932725, 39.934789, 116.622593, 116.628198),
		NewRegion("2", "255,179,128", 39.932725, 39.934789, 116.628198, 116.632007),
		NewRegion("3", "255,255,128", 39.932725, 39.934789, 116.632007, 116.639374),
	}
)

func init() {
	clientMessageTypeCodeCheck()

	mapData := loadMapData()
	//加载地图数据
	f := func(o interface{}) interface{} {
		pos := o.(*Position)
		if pos.HasOrder {
			return NewOrder(generateOrderID(), pos)
		}
		return nil
	}
	orders := mapData.Points.Map(f).transform(Sys_Type_Order).(OrderList)
	// orders := OrderList{} //所有的订单
	DebugPrintList_Info(orders)

	room := NewRoom()
	room.init()
	room.addEventSubscriber(distributorRoomEventHandler,
		WsRoomEventCode_Online, WsRoomEventCode_Offline, WsRoomEventCode_Other)

	filter := func(d *Distributor) bool {
		l := []string{"d01", "d02"}
		for _, s := range l {
			if s == d.ID {
				return true
			}
		}
		return false
	}
	g_UnitCenter = NewDistributorProcessUnitCenter(room, g_distributorStore.clone(filter), orders, mapData)
	// g_UnitCenter.start()
	startCenterRunning(g_UnitCenter)

	// g_UnitCenter.Process(NewMessageWithClient(pro_order_select_response, "", map[string]interface{}{"OrderID": "900100001", "DistributorID": "d01"}))
	// g_UnitCenter.Process(NewMessageWithClient(pro_order_select_response, "", map[string]interface{}{"OrderID": "900100002", "DistributorID": "d01"}))
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
		distributor := g_UnitCenter.distributors.findOne(func(d *Distributor) bool { return d.ID == msg.TargetID })
		// distributor := g_UnitCenter.distributors.find(msg.TargetID)
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
