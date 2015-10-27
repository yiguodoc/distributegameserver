package controllers

import (
	// "errors"
	"encoding/json"
	"fmt"
	// "math"
	"github.com/gorilla/websocket"
	"github.com/ungerik/go-dry"
	"math/rand"
	"time"
)

var game_index = 0

func getGameUnitID() string {
	game_index++
	return fmt.Sprintf("game_%d", game_index)
}

type GameUnitPreditor func(*GameUnit) bool
type GameUnitList []*GameUnit

func (gul GameUnitList) containsDistributor(distributorID string) (*GameUnit, *Distributor) {
	if len(gul) <= 0 {
		return nil, nil
	}
	if d := gul[0].containsDistributor(distributorID); d != nil {
		return gul[0], d
	} else {
		return gul[1:].containsDistributor(distributorID)
	}
}

func (gl GameUnitList) findOne(p GameUnitPreditor) *GameUnit {
	if len(gl) <= 0 {
		return nil
	}
	if p(gl[0]) {
		return gl[0]
	} else {
		return gl[1:].findOne(p)
	}
}
func (gl GameUnitList) find(p GameUnitPreditor) GameUnitList {
	return gl.findRecursive(p, GameUnitList{})
}
func (gl GameUnitList) findRecursive(p GameUnitPreditor, l GameUnitList) GameUnitList {
	if len(gl) <= 0 {
		return l
	}
	if p(gl[0]) {
		l = append(l, gl[0])
	}
	return gl[1:].findRecursive(p, l)
}

type GameUnit struct {
	ID                string
	Distributors      DistributorList
	orders            OrderList
	MapName           string
	GameTimeMaxLength int                                                //游戏最大时长
	TimeElapse        int                                                //运行时间
	chanEvent         chan *MessageWithClient                            `json:"-"`
	chanStop          chan bool                                          `json:"-"`
	processors        map[ClientMessageTypeCode]MessageWithClientHandler `json:"-"`
	supportPro        []ClientMessageTypeCode                            `json:"-"`
	distributorIDList []string                                           `json:"-"`
	mapData           *MapData                                           `json:"-"`
	gameStarted       bool                                               `json:"-"`
	// units             DistributorProcessUnitList
	// mapDataLoader     func() *MapData
	// wsRoom            *WsRoom
}

func NewGameUnit(distributorIDList []string, mapName string, timeMaxLength int) *GameUnit {
	unit := &GameUnit{
		ID:                getGameUnitID(),
		chanEvent:         make(chan *MessageWithClient, 128),
		chanStop:          make(chan bool),
		MapName:           mapName,
		processors:        make(ProHandlerMap),
		distributorIDList: distributorIDList,
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
			pro_reset_destination_request,
			pro_game_time_tick,
			pro_change_state_request,
			pro_sign_order_request,
			pro_distributor_info_request,
		},
		GameTimeMaxLength: timeMaxLength,
	}
	unit.processors = handler_map.generateHandlerMap(unit.supportPro, unit)
	return unit
}
func (gu *GameUnit) containsDistributor(id string) *Distributor {
	return gu.Distributors.findOne(func(d *Distributor) bool { return d.UserInfo.ID == id })
}

// 上线
func (gu *GameUnit) distributorOnLine(distributor *Distributor, conn *websocket.Conn) {
	// distributor.SetConn(conn)
	//处理上线事件
	DebugInfoF("配送员 %s 上线", distributor.UserInfo.Name)
	gu.Process(NewMessageWithClient(pro_on_line, distributor, conn))
}

func (gu *GameUnit) distributorOffLine(distributor *Distributor) {
	// distributor.SetOffline()
	DebugInfoF("配送员 %s 离线", distributor.UserInfo.Name)
	//处理下线事件
	gu.Process(NewMessageWithClient(pro_off_line, distributor, distributor))
}
func (gu *GameUnit) distributorMessageIn(distributor *Distributor, content []byte) {
	var msg MessageWithClient
	err := json.Unmarshal(content, &msg)
	if err != nil {
		DebugSysF("解析数据出错：%s", err)
		return
	}
	msg.Target = distributor
	gu.Process(&msg)
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func (gu *GameUnit) broadcastMsgToSubscribers(protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	msg := NewMessageWithClient(protocal, nil, obj, err...)
	data, e := json.Marshal(msg)
	if e != nil {
		DebugMustF("Fail to marshal event: %s", e)
		return
	}
	gu.Distributors.forEach(func(d *Distributor) {
		if d.SendBinaryMessage(data) != nil {
			// User disconnected.
			// dpc.distributorOffLine(d)
		}
	})
}

// send messages to WebSocket special user.
func (dpc *GameUnit) sendMsgToSpecialSubscriber(distributor *Distributor, protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	// func (dpc *GameUnit) sendMsgToSpecialSubscriber(id string, protocal ClientMessageTypeCode, obj interface{}, err ...string) {
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
		// dpc.distributorOffLine(distributor)
	}
	// } else {
	// 	DebugSysF("系统异常，无法向 %d 发送消息", id)
	// }
	if protocal != pro_2c_sys_time_elapse {
		DebugTraceF("=>  %s : %v", distributor.UserInfo.ID, msg)
	}
}

func (gu *GameUnit) stop() {
	// gu.stopAllUnits()
	if gu.chanStop != nil {
		gu.chanStop <- true
		gu.chanStop = nil
	}
	gu.Distributors.forEach(func(d *Distributor) {
		d.GameID = ""
	})
	// time.Sleep(2 * time.Second)
}
func (gu *GameUnit) start() *GameUnit {
	gu.mapData = loadMapData(gu.MapName)

	gu.orders = gu.mapData.Points.filter(func(pos *Position) bool {
		return pos.PointType == POSITION_TYPE_ORDER
	}).Map(OrderList{}, func(pos *Position, list interface{}) interface{} {
		o := NewOrder(generateOrderID(), pos)
		return append(list.(OrderList), o)
	}).(OrderList).random(rand.New(rand.NewSource(time.Now().UnixNano())), OrderList{})

	gu.Distributors = g_var.distributors.filter(func(d *Distributor) bool { return dry.StringListContains(gu.distributorIDList, d.UserInfo.ID) })
	DebugSysF("%d distributors from %s", len(gu.Distributors), gu.distributorIDList)
	bornPoints := gu.mapData.Points.filter(func(p *Position) bool { return p.IsBornPoint }).random(rand.New(rand.NewSource(time.Now().UnixNano())), PositionList{})
	// i := len(bornPoints)
	// DebugInfoF("出生点数量 => %d", len(bornPoints))
	// positionGenerator := func() *Position {
	// 	i--
	// 	if i < 0 {
	// 		i = len(bornPoints) - 1
	// 	}
	// 	return bornPoints[i]
	// }
	positionGenerator := func(bornPoints PositionList) func() *Position {
		i := len(bornPoints)
		return func() *Position {
			i--
			if i < 0 {
				i = len(bornPoints) - 1
			}
			return bornPoints[i]
		}
	}
	gu.Distributors.forEach(func(distributor *Distributor) {
		distributor.setCheckPoint(checkpoint_flag_origin)
		distributor.GameTimeMaxLength = gu.GameTimeMaxLength
		distributor.StartPos = positionGenerator(bornPoints)()
		distributor.CurrentPos = distributor.StartPos.copyTemp(true)
		distributor.NormalSpeed = defaultSpeed
		distributor.CurrentSpeed = defaultSpeed
		distributor.AcceptedOrders = OrderList{}
		distributor.GameID = gu.ID
		// newUnit(gu, distributor)
	})

	go func() {
		timer := time.Tick(1 * time.Second) //计时器功能
		breakLoop := false
		for {
			select {
			case msg := <-gu.chanEvent:
				DebugTraceF("<- %v", msg)
				if processor, ok := gu.processors[msg.MessageType]; ok { //首先自行处理
					processor(msg)
				} else {
					DebugSysF("未找到消息处理单位：%s", msg)
				}
			case <-timer:
				if gu.TimeElapse < gu.GameTimeMaxLength && gu.gameStarted == true { //尚处于单局游戏时间内
					gu.TimeElapse++
					gu.Distributors.forEach(func(distributor *Distributor) {
						gu.Process(NewMessageWithClient(pro_game_time_tick, distributor, nil))
					})
				} else if gu.gameStarted == true && gu.TimeElapse >= gu.GameTimeMaxLength { //游戏时间到达最终时限
					DebugSysF("游戏到达最终时限，开始统计成绩")
					gu.Process(NewMessageWithClient(pro_game_timeout, nil, gu))
				} else {
					// DebugSysF("没有逻辑处理")
				}
			case <-gu.chanStop:
				breakLoop = true
			}
			if breakLoop {
				break
			}
		}
		DebugSysF("跳出计时循环")

	}()
	DebugInfoF("配送系统处理中心开始运行...")
	return gu
}

//重置游戏
//参与者状态清零，通知客户端游戏重置，客户端采取相应的措施
func (dpc *GameUnit) restart() {
	DebugInfoF("游戏重新启动...")
	dpc.broadcastMsgToSubscribers(pro_2c_restart_game, nil)
	dpc.stop()
	dpc.start()
	DebugInfoF("游戏重新启动完成")
}
func (dpc *GameUnit) startGameTiming() {
	dpc.gameStarted = true
}
func (dpc *GameUnit) Process(msg *MessageWithClient) {
	DebugInfoF("<- %s", msg)
	if dpc.chanEvent != nil {

		dpc.chanEvent <- msg
	}
}

// func (dpc *GameUnit) stopAllUnits() {
// 	for id, _ := range dpc.units {
// 		dpc.stopUnit(id)
// 		// go u.stop()
// 	}
// }
// func (dpc *GameUnit) stopUnit(id string) {
// 	if u, ok := dpc.units[id]; ok {
// 		go u.stop()
// 		// delete(dpc.units, id)
// 	} else {
// 		DebugSysF("配送处理单元 %s 不存在", id)
// 	}
// }

// func (dpc *GameUnit) startUnit(id string) {
// 	if unit, ok := dpc.units[id]; ok {
// 		unit.start()
// 	} else {
// 		DebugSysF("启动配送处理单元出错，指定的单元 %s 不存在", id)
// 	}
// }
// func (dpc *GameUnit) startAlltUnit() {
// 	for _, unit := range dpc.units {
// 		unit.start()
// 	}
// }
