package controllers

import (
	// "errors"
	// "fmt"
	// "encoding/json"
	// "math"
	"time"
)

type DataWithID interface {
	GetID() string
}
type DistributorProcessUnitList map[string]*DistributorProcessUnit

func (l DistributorProcessUnitList) forEach(f func(*DistributorProcessUnit)) {
	if f == nil {
		return
	}
	for _, v := range l {
		f(v)
	}
}

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
	wsRoom            *WsRoom
	GameTimeMaxLength int //游戏最大时长
	TimeElapse        int //运行时间
	gameStarted       bool
}

func NewDistributorProcessUnitCenter(distributors DistributorList, orders OrderList, mapData *MapData, timeMaxLength int) *DistributorProcessUnitCenter {
	center := &DistributorProcessUnitCenter{
		units:        DistributorProcessUnitList{},
		chanEvent:    make(chan *MessageWithClient),
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
		},
		GameTimeMaxLength: timeMaxLength,
	}
	center.processors = handler_map.generateHandlerMap(center.supportPro, center)
	center.wsRoom = NewRoom().addEventSubscriber(distributorRoomEventHandlerGenerator(center), WsRoomEventCode_Online, WsRoomEventCode_Offline, WsRoomEventCode_Other)
	return center
}
func (dpc *DistributorProcessUnitCenter) stop() {
	if dpc.wsRoom != nil {
		dpc.wsRoom.stop()
	}
	if dpc.chanStop != nil {
		dpc.chanStop <- true
		dpc.chanStop = nil
	}
}
func (dpc *DistributorProcessUnitCenter) start() *DistributorProcessUnitCenter {
	dpc.distributors.forEach(func(distributor *Distributor) {
		distributor.GameTimeMaxLength = dpc.GameTimeMaxLength
		dpc.newUnit(distributor)
	})
	if dpc.mapDataLoader != nil {
		dpc.mapData = dpc.mapDataLoader()
	}
	dpc.wsRoom.start()
	go func() {
		timer := time.Tick(1 * time.Second) //计时器功能
		for {
			select {
			case msg := <-dpc.chanEvent:
				DebugInfoF("%v", msg)
				if processor, ok := dpc.processors[msg.MessageType]; ok { //首先自行处理
					go processor(msg)
				} else {
					if unit, ok := dpc.units[msg.TargetID]; ok { //之后交于处理单位处理
						go unit.process(msg)
					} else {
						DebugSysF("未找到消息处理单位：%s", msg)
					}
				}
			case <-timer:
				f := func(code ClientMessageTypeCode) {
					dpc.units.forEach(func(unit *DistributorProcessUnit) {
						go unit.process(NewMessageWithClient(pro_game_time_tick, unit.distributor.ID, unit))
					})
				}
				if dpc.TimeElapse < dpc.GameTimeMaxLength && dpc.gameStarted == true { //尚处于单局游戏时间内
					dpc.TimeElapse++
					f(pro_game_time_tick)
					// for _, unit := range dpc.units {
					// 	go unit.process(NewMessageWithClient(pro_game_time_tick, "", nil))
					// }
				} else { //游戏时间到达最终时限
					f(pro_end_game_request)
				}
			case <-dpc.chanStop:
				break
			}
		}
	}()
	return dpc
}
func (dpc *DistributorProcessUnitCenter) startGameTiming() {
	dpc.gameStarted = true
}
func (dpc *DistributorProcessUnitCenter) Process(msg *MessageWithClient) {
	dpc.chanEvent <- msg
}

func (dpc *DistributorProcessUnitCenter) newUnit(distributor *Distributor) *DistributorProcessUnit {
	if u, ok := dpc.units[distributor.ID]; ok {
		DebugInfoF("配送处理单元 %s 重复添加", distributor.ID)
		return u
	} else {
		unit := &DistributorProcessUnit{
			center:      dpc,
			distributor: distributor,
			chanStop:    make(chan bool),
			processors:  make(ProHandlerMap),
			supportPro: []ClientMessageTypeCode{
				pro_game_time_tick,
				pro_reset_destination_request,
				pro_change_state_request,
				pro_sign_order_request,
				pro_distributor_info_request,
				// pro_prepared_for_select_order,
			},
			// processors:  processors,
		}
		dpc.units[distributor.ID] = unit
		unit.processors = handler_map.generateHandlerMap(unit.supportPro, unit)
		// unit.start()
		return unit
	}
}
func (dpc *DistributorProcessUnitCenter) stopUnit(id string) {
	if u, ok := dpc.units[id]; ok {
		u.stop()
	} else {
		DebugSysF("配送处理单元 %s 不存在", id)
	}
}

// func (dpc *DistributorProcessUnitCenter) removeUnit(id string) {
// 	if u, ok := dpc.units[id]; ok {
// 		u.stop()
// 		delete(dpc.units, id)
// 	} else {
// 		DebugSysF("配送处理单元 %s 不存在，无法移除", id)
// 	}
// }
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

type DistributorProcessUnit struct {
	center      *DistributorProcessUnitCenter
	processors  ProHandlerMap
	chanEvent   (chan *MessageWithClient)
	distributor *Distributor
	chanStop    chan bool
	supportPro  []ClientMessageTypeCode
	// processors  []func(*MessageWithClient, *DistributorProcessUnit)
	// chanEvent   (chan []byte)
}

func (u *DistributorProcessUnit) addProcessor(generators ProHandlerGeneratorMap) {
	u.processors = handler_map.generateHandlerMap(u.supportPro, u)
}
func (u *DistributorProcessUnit) process(data *MessageWithClient) {
	if u.chanEvent != nil {
		u.chanEvent <- data
	}
}
func (u *DistributorProcessUnit) stop() {
	if u.chanStop != nil {
		u.chanStop <- true
	}
	u.chanEvent = nil
}
func (u *DistributorProcessUnit) start() {
	DebugInfoF("处理单元 %s %s 启动", u.distributor.ID, u.distributor.Name)

	u.chanEvent = make(chan *MessageWithClient)
	f := func() {
		// timer := time.Tick(1 * time.Second) //计时器功能
		for {
			select {
			case <-u.chanStop:
				break
			case msg := <-u.chanEvent:
				if processor, ok := u.processors[msg.MessageType]; ok {
					processor(msg)
				} else {
					DebugSysF("未找到消息处理单位：%s", msg)
				}
			}
		}
	}
	go f()
}
