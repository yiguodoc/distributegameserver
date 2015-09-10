package controllers

import (
// "errors"
// "fmt"
// "encoding/json"
// "math"
// "time"
)

// type DataWithID interface {
// 	GetID() string
// }
type DistributorProcessUnitList map[string]*DistributorProcessUnit

func (l DistributorProcessUnitList) forEach(f func(*DistributorProcessUnit)) {
	if f == nil {
		return
	}
	for _, v := range l {
		f(v)
	}
}

func newUnit(dpc *DistributorProcessUnitCenter, distributor *Distributor) *DistributorProcessUnit {
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
		}
		dpc.units[distributor.ID] = unit
		unit.processors = handler_map.generateHandlerMap(unit.supportPro, unit)
		return unit
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

	u.chanEvent = make(chan *MessageWithClient, 128)
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
