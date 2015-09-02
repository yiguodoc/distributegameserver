// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	// "container/list"
	// "github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	// "time"
	"encoding/json"
	// "strings"
	// "errors"
	// "fmt"
)

type Subscriber interface {
	SetConn(conn *websocket.Conn)
	SendBinaryMessage(msg []byte) error
	IdEqals(id string) bool
	SetOffline() error
	String() string
	GetID() string
}
type SubscriberList []Subscriber

func (sl SubscriberList) findByID(id string) Subscriber {
	for _, sub := range sl {
		if sub.IdEqals(id) == true {
			return sub
		}
	}
	return nil
}
func (sl SubscriberList) remove(id string) (l SubscriberList) {
	for _, sub := range sl {
		if sub.IdEqals(id) == false {
			l = append(l, sub)
		}
	}
	return
}
func (s SubscriberList) ListName() string {
	return "在线人员"
}
func (s SubscriberList) InfoList() (l []string) {
	for _, sub := range s {
		l = append(l, sub.String())
	}
	return
}

type roomMessage struct {
	targetID string
	content  interface{}
	// msg      *MessageWithClient
}
type WsRoomEventCode int

var (
	WsRoomEventCode_Online  WsRoomEventCode = 1
	WsRoomEventCode_Offline WsRoomEventCode = 2
	WsRoomEventCode_Other   WsRoomEventCode = 3
)

//对应房间来说，只需要注意房间中人的变化及消息的分发和接收，至于人的变化和消息的处理则需要由系统来处理
//房间对外开放接收消息的接口，并通过外部导入的处理对象完成消息向外部环境的分发
//房间的人也是抽象的满足一系列接口的对象，配送员、观察者都可以进行抽象
type WsRoom struct {
	chanSubscribe            (chan Subscriber)   // Channel for new join users.
	chanUnsubscribe          (chan string)       // Channel for exit users.
	chanPublishToSubscribers (chan *roomMessage) // Send events here to publish them.
	chanMessage              (chan *roomMessage) // new message comes in and need to do something
	subscribers              SubscriberList      //配送人员列表，在线或者不在线
	eventSubscribers         SysEventSubscribeList
}

//******************************************************************
//wsroom helper function

//新 user 加入room
func (w *WsRoom) join(sub Subscriber, conn *websocket.Conn) {
	if sub == nil {
		return
	}
	sub.SetConn(conn)
	w.chanSubscribe <- sub
}

//离开或者掉线
func (w *WsRoom) leave(id string) {
	w.chanUnsubscribe <- id
}

//新消息处理
func (w *WsRoom) newMessage(id string, content []byte) {
	DebugTraceF("<= : %s", string(content))
	w.chanMessage <- &roomMessage{
		targetID: id,
		content:  content,
	}
	// var msg MessageWithClient
	// err := json.Unmarshal(content, &msg)
	// if err != nil {
	// 	DebugMustF("解析客户端信息时出错：%s", err)
	// 	DebugTrace(string(content))
	// } else {
	// 	w.chanMessage <- &msg
	// }
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func (w *WsRoom) broadcastMsgToSubscribers(protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	// DebugTraceF("broadcastWebSocket => %s", event)
	// msg := &MessageWithClient{protocal, "", obj, ""}
	msg := NewMessageWithClient(protocal, "", obj, err...)
	w.chanPublishToSubscribers <- &roomMessage{content: msg}
}

// send messages to WebSocket special user.
func (w *WsRoom) sendMsgToSpecialSubscriber(id string, protocal ClientMessageTypeCode, obj interface{}, err ...string) {
	// DebugTraceF("broadcastWebSocket => %s", event)
	msg := NewMessageWithClient(protocal, id, obj, err...)
	// msg := &MessageWithClient{protocal, id, obj}
	w.chanPublishToSubscribers <- &roomMessage{targetID: id, content: msg}
	if protocal != pro_2c_sys_time_elapse {

		DebugTraceF("=>  %s : %v", id, msg)
	}
}

//******************************************************************

func NewRoom() *WsRoom {
	return &WsRoom{
		chanSubscribe:            make(chan Subscriber),   // Channel for new join users.
		chanUnsubscribe:          make(chan string),       // Channel for exit users.
		chanPublishToSubscribers: make(chan *roomMessage), // Send events here to publish them.
		chanMessage:              make(chan *roomMessage), // new message comes in and need to do something
		subscribers:              SubscriberList{},        //配送人员列表，在线或者不在线
		eventSubscribers:         SysEventSubscribeList{},
	}
}
func (w *WsRoom) init() {
	go func() {
		for {
			select {
			case sub := <-w.chanSubscribe:
				w.newUserOnline(sub)
			case id := <-w.chanUnsubscribe:
				w.setUserOffline(id)
			case msg := <-w.chanMessage: //有信息从客户端传入，
				// if w.EventReceiver != nil {
				// 	w.EventReceiver(msg.id, msg.content)
				// }
				w.eventSubscribers.notifyEventSubscribers(int(WsRoomEventCode_Other), msg)
				// DebugTraceF("message => %s", msg)
				// if code, err := msg.MessageType.mapToSysEventCode(); err == nil {
				// 	w.triggerSysEvent(code, msg.Data)
				// } else {
				// 	DebugInfoF("客户端事件 %s 未处理", msg.MessageType.name())
				// }
			case msg := <-w.chanPublishToSubscribers:
				// DebugTraceF("publish distributors => %s", msg.msg)
				// distributors := w.subscribers.filtedSubscribers(subscriber_type_distributor)
				w.writeMsgToSubscribers(msg)
			}
		}
	}()
}
func (w *WsRoom) writeMsgToSubscribers(roommsg *roomMessage) {
	msg := roommsg.content
	data, err := json.Marshal(msg)
	if err != nil {
		DebugMustF("Fail to marshal event: %s", err)
		return
	}
	if len(roommsg.targetID) > 0 {
		s := w.subscribers.findByID(roommsg.targetID)
		if s != nil {
			if s.SendBinaryMessage(data) != nil {
				// User disconnected.
				w.leave(s.GetID()) //统一函数调用
			}
		}
	} else {
		for _, s := range w.subscribers {
			// DebugTraceF("send msg  to [%s]", s.GetID())
			if s.SendBinaryMessage(data) != nil {
				// User disconnected.
				w.leave(s.GetID()) //统一函数调用
			}
			// }
		}
	}
}

//为该房间提供事件订阅功能
func (w *WsRoom) addEventSubscriber(f eventProcessor, codes ...WsRoomEventCode) {
	// w.eventSubscribers.addEventSubscriber(f, codes...)
	for _, code := range codes {
		w.eventSubscribers.addEventSubscriber(f, int(code))
	}
}

//移除对该房间的事件的订阅
func (w *WsRoom) removeEventSubscriber(code WsRoomEventCode, process eventProcessor) {
	w.eventSubscribers.removeEventSubscriber(int(code), process)
}

func (w *WsRoom) setUserOffline(id string) {
	sub := w.subscribers.findByID(id)
	if sub != nil {
		if err := sub.SetOffline(); err == nil {
			// w.triggerSysEvent(sys_event_user_offline, sub)
			w.eventSubscribers.notifyEventSubscribers(int(WsRoomEventCode_Offline), NewMessageWithClient(pro_off_line, id, id))
		} else {
			DebugMustF("[%s] OffLine,But  err: %s", id, err)
		}
		w.subscribers = w.subscribers.remove(id)
	}
	DebugTraceF("下线后，最新用户人数：%d", len(w.subscribers))
}

func (w *WsRoom) newUserOnline(sub Subscriber) {
	if sub == nil {
		return
	}
	// DebugTraceF("当前用户人数：%d", len(w.subscribers))
	subTemp := w.subscribers.findByID(sub.GetID())
	if subTemp == nil {
		w.subscribers = append(w.subscribers, sub)
		// w.triggerSysEvent(sys_event_user_online, sub)
		w.eventSubscribers.notifyEventSubscribers(int(WsRoomEventCode_Online), NewMessageWithClient(pro_on_line, sub.GetID(), sub.GetID()))
	}
	DebugTraceF("上线后，最新用户人数：%d", len(w.subscribers))
}
