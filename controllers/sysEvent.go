package controllers

import (
	// "container/list"
	// "github.com/astaxie/beego"
	// "github.com/gorilla/websocket"
	// "time"
	// "encoding/json"
	// "strings"
	"fmt"
)

func sysEventCodeCheck() {
	l := getSysEventCodeList()
	for _, code := range l {
		DebugInfoF("系统事件定义：%d %s", code, code.name())
	}
	// for i := 0; i < int(sys_event_max); i++ {
	// }
}
func getSysEventCodeList() (l []sysEventCode) {
	for i := 0; i < int(sys_event_max); i++ {
		l = append(l, sysEventCode(i))
	}
	return
}

type sysEventCode int

var (
	sys_event_user_online              sysEventCode = 0  //上线
	sys_event_user_offline             sysEventCode = 1  //下线
	sys_event_order_select_response    sysEventCode = 2  //订单选择的反馈
	sys_event_count_down               sysEventCode = 3  //倒计时
	sys_event_begin_select_order       sysEventCode = 4  //向配送员发送开始抢订单的命令
	sys_event_message_broadcast        sysEventCode = 5  //向配送员发布信息广播
	sys_event_order_select_result      sysEventCode = 6  //向配送员发布订单分配结果
	sys_event_start_order_selection    sysEventCode = 7  //开始配送员分发订单
	sys_event_distributor_prepared     sysEventCode = 8  //配送员为订单分配做好准备的通知
	sys_event_count_down_silent        sysEventCode = 9  //倒计时，但不会将倒计时发送给客户端
	sys_event_start_order_distribution sysEventCode = 10 //开始配送员分发订单
	sys_event_start_order_distribute   sysEventCode = 11 //向配送员发布订单分配结果
	sys_event_order_distribute_result  sysEventCode = 12 //向配送员发布订单分配结果
	sys_event_max                                   = 13 //最大值，保证比定义的事件的最大值大1
)

func (se sysEventCode) name() (s string) {
	switch se {
	case sys_event_user_online:
		s = "sys_event_user_online"
	case sys_event_user_offline:
		s = "sys_event_user_offline"
	case sys_event_order_select_response:
		s = "sys_event_order_select_response"
	case sys_event_count_down:
		s = "sys_event_count_down"
	case sys_event_begin_select_order:
		s = "sys_event_begin_select_order"
	case sys_event_message_broadcast:
		s = "sys_event_message_broadcast"
	case sys_event_order_select_result:
		s = "sys_event_order_select_result"
	case sys_event_start_order_selection:
		s = "sys_event_start_order_selection"
	case sys_event_distributor_prepared:
		s = "sys_event_distributor_prepared"
	case sys_event_count_down_silent:
		s = "sys_event_count_down_silent"
	case sys_event_start_order_distribution:
		s = "sys_event_start_order_distribution"
	case sys_event_order_distribute_result:
		s = "sys_event_order_distribute_result"
	default:
		if int(se) < sys_event_max {
			panic(fmt.Sprintf("事件(%d)定义描述不完全", int(se)))
		}
	}
	return
}

// var sysEventDefination = map[string]int{
// 	"user_online":             0, //上线
// 	"user_offline":            1, //下线
// 	"order_select_response":   2, //订单选择的反馈
// 	"count_down":              3, //倒计时
// 	"begin_select_order":      4, //向配送员发送开始抢订单的命令
// 	"message_broadcast":       5, //向配送员发布信息广播
// 	"order_distribute_result": 6, //向配送员发布订单分配结果
// 	"start_order_selection":   7, //开始配送员分发订单
// 	"distributor_prepared":    8, //配送员为订单分配做好准备的通知
// 	"count_down_silent":       9, //倒计时，但不会将倒计时发送给客户端
// 	// "start_order_distribution": 10, //开始配送员分发订单
// 	// "start_order_distribution_countdown": 8, //开始配送员分发订单前的倒计时，完成后出发事件 start_order_distribution
// }

// func getSysEventDefValue(key string) int {
// 	if value, ok := sysEventDefination[key]; ok {
// 		return value
// 	} else {
// 		panic("没有事件：" + key + " 的定义")
// 		return -1
// 	}
// }

var eventPkgCount int64 = 1

func getPkgCount() int64 {
	eventPkgCount += 2
	return eventPkgCount
}

//系统通用事件，整个系统内流通
type SysEvent struct {
	eventCode sysEventCode
	// eventName string
	data      interface{}
	nextEvent *SysEvent
	// pkgID     int64
}

func (s *SysEvent) setNextEvent(next *SysEvent) {
	s.nextEvent = next
}
func NewSysEvent(code sysEventCode, data interface{}) *SysEvent {
	// func NewSysEvent(eventName string, data interface{}) *SysEvent {
	// code, ok := sysEventDefination[eventName]
	// if !ok {
	// 	panic("没有该事件的定义")
	// }
	// code := getSysEventDefValue(eventName)
	return &SysEvent{
		eventCode: code,
		data:      data,
		// eventName: eventName,
	}
}
func NewCountDownEvent(count int) *SysEvent {
	return NewSysEvent(sys_event_count_down, count)
}
func NewMessageBroadcastEvent(msg string) *SysEvent {
	return NewSysEvent(sys_event_message_broadcast, msg)
}

// type SysEventPkg struct {
// 	id     int64
// 	events []*SysEvent
// }

// func NewSysEventPkg(args ...*SysEvent) *SysEventPkg {
// 	count := getPkgCount()
// 	for _, e := range args {
// 		e.pkgID = count
// 	}
// 	if len(args) > 1 {
// 		for i := 0; i < len(args)-1; i++ {
// 			args[i].setNextEvent(args[i+1])
// 		}
// 	}

// 	return &SysEventPkg{
// 		id:     count,
// 		events: args,
// 	}
// }

// type SysEventPkgList []*SysEventPkg

// func (s SysEventPkgList) add(pkg ...*SysEventPkg) SysEventPkgList {
// 	return append(s, pkg...)
// }
// func (s SysEventPkgList) remove(id int64) (list SysEventPkgList) {
// 	for _, pkg := range s {
// 		if pkg.id != id {
// 			list = append(list, pkg)
// 		}
// 	}
// 	return
// }
