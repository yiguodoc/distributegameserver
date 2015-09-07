package controllers

import (
	// "github.com/ssor/fauxgaux"
	// "github.com/gorilla/websocket"
	"encoding/json"
	"time"
	// "strings"
	"fmt"
	// "reflect"
	"math/rand"
)

var default_time_of_one_loop = 5 * 60

var (
	g_UnitCenter       *DistributorProcessUnitCenter
	g_room_viewer      *WsRoom            //= NewRoom(eventReceiver)
	g_distributorStore = DistributorList{ //配送员列表
		NewDistributor("d01", "张军", 2, color_orange),
		NewDistributor("d02", "刘晓莉", 2, color_red),
		NewDistributor("d03", "桑鸿庆", 2, color_purple),
	}
	g_regions = RegionList{
		NewRegion("1", "255,128,128", 39.928935, 39.944789, 116.614041, 116.618676),
		NewRegion("2", "255,179,128", 39.928935, 39.944789, 116.618676, 116.625898),
		NewRegion("3", "255,255,128", 39.928935, 39.944789, 116.625898, 116.639373),
	}
)

func init() {
	i := 4
	j := float64(6.5)
	fmt.Println(fmt.Sprintf("i = %v j = %v", i, j))

	clientMessageTypeCodeCheck()

	//加载地图数据
	mapData := loadMapData()

	orders := mapData.Points.filter(func(pos *Position) bool {
		return pos.HasOrder
	}).Map(OrderList{}, func(pos *Position, list interface{}) interface{} {
		o := NewOrder(generateOrderID(), pos)
		return append(list.(OrderList), o)
	}).(OrderList) //.random()
	// orders := OrderList{} //所有的订单
	DebugPrintList_Info(orders)
	orders = orders.random(rand.New(rand.NewSource(time.Now().UnixNano())))
	DebugPrintList_Info(orders)

	room := NewRoom().start()
	room.addEventSubscriber(distributorRoomEventHandler, WsRoomEventCode_Online, WsRoomEventCode_Offline, WsRoomEventCode_Other)

	filter := func(d *Distributor) bool {
		l := []string{"d01", "d02", "d03"}
		for _, s := range l[:2] {
			if s == d.ID {
				return true
			}
		}
		return false
	}
	g_UnitCenter = NewDistributorProcessUnitCenter(room, g_distributorStore.clone(filter), orders, mapData, default_time_of_one_loop)
	g_UnitCenter.start()

	// g_UnitCenter.Process(NewMessageWithClient(pro_order_select_response, "", map[string]interface{}{"OrderID": "900100001", "DistributorID": "d01"}))
	// g_UnitCenter.Process(NewMessageWithClient(pro_order_select_response, "", map[string]interface{}{"OrderID": "900100002", "DistributorID": "d01"}))

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
