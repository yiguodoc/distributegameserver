package controllers

import (
	// "errors"
	"encoding/json"
	// "fmt"
	// "math"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

type DistributorProcessUnitCenter struct {
	units             DistributorProcessUnitList
	chanEvent         chan *MessageWithClient
	chanStop          chan bool
	processors        map[ClientMessageTypeCode]MessageWithClientHandler
	supportPro        []ClientMessageTypeCode
	distributors      DistributorList
	orders            OrderList
	mapData           *MapData
	mapDataLoader     func() *MapData
	GameTimeMaxLength int //游戏最大时长
	TimeElapse        int //运行时间
	gameStarted       bool
	// wsRoom            *WsRoom
}

func NewDistributorProcessUnitCenter(distributors DistributorList, orders OrderList, mapData *MapData, timeMaxLength int) *DistributorProcessUnitCenter {
	if len(mapData.Points.filter(func(p *Position) bool { return p.IsBornPoint == true })) <= 0 {
		DebugSysF("地图中没有出生点信息")
		return nil
	}
	center := &DistributorProcessUnitCenter{
		units:        DistributorProcessUnitList{},
		chanEvent:    make(chan *MessageWithClient, 128),
		chanStop:     make(chan bool),
		orders:       orders,
		mapData:      mapData,
		processors:   make(ProHandlerMap),
		distributors: distributors,
		supportPro: []ClientMessageTypeCode{
			pro_game_start,
			pro_order_select_response,
			pro_move_from_route_to_node,
			pro_move_from_node_to_route,
			pro_on_line,
			pro_off_line,
			pro_prepared_for_select_order,
			pro_end_game_request,
			pro_game_timeout,
		},
		GameTimeMaxLength: timeMaxLength,
	}
	center.processors = handler_map.generateHandlerMap(center.supportPro, center)
	// center.wsRoom = NewRoom().addEventSubscriber(distributorRoomEventHandlerGenerator(center), WsRoomEventCode_Online, WsRoomEventCode_Offline, WsRoomEventCode_Other)
	return center
}
func (dpc *DistributorProcessUnitCenter) containsDistributor(id string) *Distributor {
	return dpc.distributors.findOne(func(d *Distributor) bool { return d.ID == id })
}

// 上线
func (dpc *DistributorProcessUnitCenter) distributorOnLine(distributor *Distributor, conn *websocket.Conn) {
	// func (dpc *DistributorProcessUnitCenter) distributorOnLine(id string, conn *websocket.Conn) {
	// DebugTraceF("配送员 %s 连接服务", distributor.Name)
	// distributor := dpc.distributors.findOne(func(d *Distributor) bool { return d.ID == id })
	// if distributor != nil {
	distributor.SetConn(conn)
	//处理上线事件
	DebugInfoF("配送员 %s 上线", distributor.Name)
	dpc.Process(NewMessageWithClient(pro_on_line, distributor, distributor))
	// } else {
	// 	DebugInfoF("不存在配送员 %s", id)
	// }
}

func (dpc *DistributorProcessUnitCenter) distributorOffLine(distributor *Distributor) {
	// func (dpc *DistributorProcessUnitCenter) distributorOffLine(id string) {
	// distributor := dpc.distributors.findOne(func(d *Distributor) bool { return d.ID == id })
	// if distributor != nil {
	distributor.SetOffline()
	//处理下线事件
	dpc.Process(NewMessageWithClient(pro_off_line, distributor, distributor))
	// }
}
func (dpc *DistributorProcessUnitCenter) distributorMessageIn(distributor *Distributor, content []byte) {
	var msg MessageWithClient
	err := json.Unmarshal(content, &msg)
	if err != nil {
		DebugSysF("解析数据出错：%s", err)
		return
	}
	msg.Target = distributor
	dpc.Process(&msg)
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func (dpc *DistributorProcessUnitCenter) broadcastMsgToSubscribers(protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	msg := NewMessageWithClient(protocal, nil, obj, err...)
	data, e := json.Marshal(msg)
	if e != nil {
		DebugMustF("Fail to marshal event: %s", e)
		return
	}
	dpc.distributors.forEach(func(d *Distributor) {
		if d.SendBinaryMessage(data) != nil {
			// User disconnected.
			dpc.distributorOffLine(d)
		}
	})
}

// send messages to WebSocket special user.
func (dpc *DistributorProcessUnitCenter) sendMsgToSpecialSubscriber(distributor *Distributor, protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	// func (dpc *DistributorProcessUnitCenter) sendMsgToSpecialSubscriber(id string, protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	msg := NewMessageWithClient(protocal, distributor, obj, err...)
	data, e := json.Marshal(msg)
	if e != nil {
		DebugMustF("Fail to marshal event: %s", e)
		return
	}
	// distributor := dpc.distributors.findOne(func(d *Distributor) bool { return id == d.ID })
	// if distributor != nil {
	if distributor.SendBinaryMessage(data) != nil {
		// User disconnected.
		dpc.distributorOffLine(distributor)
	}
	// } else {
	// 	DebugSysF("系统异常，无法向 %d 发送消息", id)
	// }
	if protocal != pro_2c_sys_time_elapse {
		DebugTraceF("=>  %s : %v", distributor.ID, msg)
	}
}

func (dpc *DistributorProcessUnitCenter) stop() {
	// if dpc.wsRoom != nil {
	// 	dpc.wsRoom.stop()
	// }
	if dpc.chanStop != nil {
		dpc.chanStop <- true
		dpc.chanStop = nil
	}
}
func (dpc *DistributorProcessUnitCenter) start() *DistributorProcessUnitCenter {
	bornPoints := dpc.mapData.Points.filter(func(p *Position) bool { return p.IsBornPoint }).random(rand.New(rand.NewSource(time.Now().UnixNano())), PositionList{})
	i := len(bornPoints)
	DebugInfoF("出生点数量 => %d", len(bornPoints))
	positionGenerator := func() *Position {
		i--
		if i < 0 {
			i = len(bornPoints) - 1
		}
		// fmt.Println("index => ", i)
		return bornPoints[i]
	}
	dpc.distributors.forEach(func(distributor *Distributor) {
		distributor.GameTimeMaxLength = dpc.GameTimeMaxLength
		distributor.StartPos = positionGenerator()
		distributor.CurrentPos = distributor.StartPos.copyTemp(true)
		distributor.NormalSpeed = defaultSpeed
		distributor.CurrentSpeed = defaultSpeed

		newUnit(dpc, distributor)
	})
	if dpc.mapDataLoader != nil {
		dpc.mapData = dpc.mapDataLoader()
	}
	// dpc.wsRoom.start()
	go func() {
		timer := time.Tick(1 * time.Second) //计时器功能
		for {
			select {
			case msg := <-dpc.chanEvent:
				DebugTraceF("<- %v", msg)
				if processor, ok := dpc.processors[msg.MessageType]; ok { //首先自行处理
					processor(msg)
					// go processor(msg)
				} else {
					if unit, ok := dpc.units[msg.Target.ID]; ok { //之后交于处理单位处理
						go unit.process(msg)
					} else {
						DebugSysF("未找到消息处理单位：%s", msg)
					}
				}
			case <-timer:
				if dpc.TimeElapse < dpc.GameTimeMaxLength && dpc.gameStarted == true { //尚处于单局游戏时间内
					dpc.TimeElapse++
					dpc.units.forEach(func(unit *DistributorProcessUnit) {
						go unit.process(NewMessageWithClient(pro_game_time_tick, unit.distributor, unit))
					})

				} else if dpc.gameStarted == true && dpc.TimeElapse >= dpc.GameTimeMaxLength { //游戏时间到达最终时限
					DebugSysF("游戏到达最终时限，开始统计成绩")
					go dpc.Process(NewMessageWithClient(pro_game_timeout, nil, dpc))
				} else {
					// DebugSysF("没有逻辑处理")
				}
			case <-dpc.chanStop:
				break
			}
		}
	}()
	DebugInfoF("配送系统处理中心开始运行...")
	return dpc
}
func (dpc *DistributorProcessUnitCenter) startGameTiming() {
	dpc.gameStarted = true
}
func (dpc *DistributorProcessUnitCenter) Process(msg *MessageWithClient) {
	// DebugInfoF("<- %s", msg)
	dpc.chanEvent <- msg
}

func (dpc *DistributorProcessUnitCenter) stopAllUnits() {
	for _, u := range dpc.units {
		go u.stop()
	}
}
func (dpc *DistributorProcessUnitCenter) stopUnit(id string) {
	if u, ok := dpc.units[id]; ok {
		go u.stop()
	} else {
		DebugSysF("配送处理单元 %s 不存在", id)
	}
}

func (dpc *DistributorProcessUnitCenter) startUnit(id string) {
	if unit, ok := dpc.units[id]; ok {
		unit.start()
	} else {
		DebugSysF("启动配送处理单元出错，指定的单元 %s 不存在", id)
	}
}
func (dpc *DistributorProcessUnitCenter) startAlltUnit() {
	for _, unit := range dpc.units {
		unit.start()
	}
}
