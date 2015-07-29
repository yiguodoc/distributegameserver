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
	msg      *MessageWithClient
}

//对应房间来说，只需要注意房间中人的变化及消息的分发和接收，至于人的变化和消息的处理则需要由系统来处理
//房间对外开放接收消息的接口，并通过外部导入的处理对象完成消息向外部环境的分发
//房间的人也是抽象的满足一系列接口的对象，配送员、观察者都可以进行抽象
type WsRoom struct {
	chanSubscribe            (chan Subscriber)         // Channel for new join users.
	chanUnsubscribe          (chan string)             // Channel for exit users.
	chanPublishToSubscribers (chan *roomMessage)       // Send events here to publish them.
	chanMessage              (chan *MessageWithClient) // new message comes in and need to do something
	EventReceiver            func(code sysEventCode, para interface{})
	subscribers              SubscriberList //配送人员列表，在线或者不在线
	// EventReceiver            func(eventName string, para interface{})
	// chanMsgToSpecialSubscribers (chan *roomMessage)       // Send events here to send to special.
	// chanPublishToViewers      (chan *MessageWithClient) // Send events here to publish them.
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
	// chanMessage <- newEvent(EVENT_MESSAGE, uriFlag, content)
	// chanMessage <- &MessageWithClient{}
	DebugTraceF("newMessage: %s", string(content))
	var msg MessageWithClient
	err := json.Unmarshal(content, &msg)
	if err != nil {
		DebugMustF("解析客户端信息时出错：%s", err)
		DebugTrace(string(content))
	} else {
		w.chanMessage <- &msg
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func (w *WsRoom) broadcastMsgToSubscribers(protocal int, obj interface{}) {
	// DebugTraceF("broadcastWebSocket => %s", event)
	msg := &MessageWithClient{protocal, obj}
	w.chanPublishToSubscribers <- &roomMessage{msg: msg}
	// data, err := json.Marshal(msg)
	// if err != nil {
	// 	DebugMustF("Fail to marshal event: %s", err)
	// 	return
	// }
	// for _, s := range g_subscribers {
	// 	if s.IsOnline == true && s.Conn != nil {
	// 		if s.Conn.WriteMessage(websocket.TextMessage, data) != nil {
	// 			// User disconnected.
	// 			leaveRoom(s.ID) //统一函数调用
	// 		}
	// 	}
	// }
	/*	for ele := subscribers.Front(); ele != nil; ele = ele.Next() {
		sub := ele.Value.(*Subscriber)
		// for _, sub := range subscribers {
		// Immediately send event to WebSocket users.
		var ws *websocket.Conn = nil
		ws = sub.Conn

		// if len(namePrefix) > 0 && strings.HasPrefix(sub.Name, namePrefix) {
		// 	// if len(namePrefix) > 0 && strings.HasPrefix(sub.Value.(Subscriber).Name, namePrefix) {
		// 	ws = sub.Conn
		// 	// ws = sub.Value.(Subscriber).Conn

		// } else {
		// 	ws = sub.Conn
		// 	// beego.Warn("no prefix")
		// }
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				leaveRoom(sub.Name) //统一函数调用
				// chanUnsubscribe <- sub.Name
				// removeUserFromRoom(sub.Name)
			}
		}

	}*/
}

// send messages to WebSocket special user.
func (w *WsRoom) sendMsgToSpecialSubscriber(id string, protocal int, obj interface{}) {
	// DebugTraceF("broadcastWebSocket => %s", event)
	msg := &MessageWithClient{protocal, obj}
	w.chanPublishToSubscribers <- &roomMessage{targetID: id, msg: msg}
}

//******************************************************************

func NewRoom(eventReceiver func(code sysEventCode, para interface{})) *WsRoom {
	// func NewRoom(eventReceiver func(eventName string, para interface{})) *WsRoom {
	return &WsRoom{
		chanSubscribe:            make(chan Subscriber),         // Channel for new join users.
		chanUnsubscribe:          make(chan string),             // Channel for exit users.
		chanPublishToSubscribers: make(chan *roomMessage),       // Send events here to publish them.
		chanMessage:              make(chan *MessageWithClient), // new message comes in and need to do something
		subscribers:              SubscriberList{},              //配送人员列表，在线或者不在线
		EventReceiver:            eventReceiver,
		// chanMsgToSpecialSubscribers: make(chan *roomMessage),       // Send events here to publish them.
		// chanPublishToViewers:      make(chan *MessageWithClient), // Send events here to publish them.
	}
}
func (w *WsRoom) init() {
	go func() {
		for {
			select {
			case sub := <-w.chanSubscribe:
				w.newUserOnline(sub)
				/*			if !isUserExist(sub.Name) {
								newUserToRoom(sub)
							} else {
								DebugInfoF("Old user: [%s]", sub.Name)
							}*/
			case id := <-w.chanUnsubscribe:
				w.setUserOffline(id)
			case msg := <-w.chanMessage: //有信息从客户端传入，将与客户端交互的消息转化成系统的消息进行处理
				DebugTraceF("message => %s", msg)
				//do something on message
				switch msg.MessageType {
				case pro_order_select_response:
					DebugTraceF("response: %s", msg.Data)
					m := msg.Data.(map[string]interface{})
					orderID, ok := m["OrderID"]
					if !ok {
						DebugMustF("客户端返回数据格式错误，无法获取订单编号")
						return
					}
					distributorID, ok := m["DistributorID"]
					if !ok {
						DebugMustF("客户端返回数据格式错误，无法获取配送员编号")
						return
					}
					// event := NewSysEvent(, )
					// triggerSysEvent(event)
					if orderID != nil && distributorID != nil {
						w.triggerSysEvent(sys_event_order_select_response, NewOrderDistribution(orderID.(string), distributorID.(string)))
					}
				case pro_distributor_prepared:
					m := msg.Data.(map[string]interface{})
					distributorID, ok := m["DistributorID"]
					if !ok {
						DebugMustF("客户端返回数据格式错误，无法获取配送员编号")
						return
					}
					DebugTraceF("配送员[%s]准备完毕", distributorID)
					w.triggerSysEvent(sys_event_distributor_prepared, distributorID)
				}
			case msg := <-w.chanPublishToSubscribers:
				// DebugTraceF("publish distributors => %s", msg.msg)
				// distributors := w.subscribers.filtedSubscribers(subscriber_type_distributor)
				w.writeMsgToSubscribers(msg)
			}
		}
	}()
}
func (w *WsRoom) writeMsgToSubscribers(roommsg *roomMessage) {
	msg := roommsg.msg
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
func (w *WsRoom) triggerSysEvent(code sysEventCode, para interface{}) {
	// func (w *WsRoom) triggerSysEvent(eventName string, para interface{}) {
	if w.EventReceiver != nil {
		w.EventReceiver(code, para)
	}
}
func (w *WsRoom) setUserOffline(id string) {
	sub := w.subscribers.findByID(id)
	if sub != nil {

		if err := sub.SetOffline(); err == nil {

			// if sub.Conn != nil {
			// 	if err := sub.Conn.Close(); err == nil {
			// 		sub.Conn = nil
			// 		DebugInfoF("[%s] OffLine WebSocket closed", sub.ID)
			// 	} else {
			// 		DebugMustF("[%s] OffLine,But close websocket err: %s", id, err)
			// 	}
			// }
			// event := NewSysEvent(getSysEventDefValue("user_offline"), sub)
			// triggerSysEvent(event)
			w.triggerSysEvent(sys_event_user_offline, sub)
		} else {
			DebugMustF("[%s] OffLine,But  err: %s", id, err)

		}
		w.subscribers = w.subscribers.remove(id)
	}
	DebugTraceF("下线后，最新用户人数：%d", len(w.subscribers))

	// for _, s := range g_subscribers {
	// 	if s.ID == id {
	// 		s.setOffline()
	// 		if s.Conn != nil {
	// 			if err := s.Conn.Close(); err == nil {
	// 				s.Conn = nil
	// 				DebugInfoF("[%s] OffLine WebSocket closed", s.ID)
	// 			} else {
	// 				DebugMustF("[%s] OffLine,But close websocket err: %s", id, err)
	// 			}
	// 		}
	// 	}
	// }
	// DebugPrintList_Trace(g_subscribers)

	/*	for ele := subscribers.Front(); ele != nil; ele = ele.Next() {
		sub := ele.Value.(*Subscriber)
		if sub.Name == id {
			subscribers.Remove(ele)
			// Close connection.
			ws := sub.Conn
			if ws != nil {
				if err := ws.Close(); err == nil {
					DebugInfoF("WebSocket closed: %s", id)
				} else {
					DebugMustF("close websocket err when %s leave: %s", id, err)
				}
			}
			// chanMessage <- newEvent(EVENT_LEAVE, id, "") // Publish a LEAVE event.
			break
		}
	}*/
}

func (w *WsRoom) newUserOnline(sub Subscriber) {
	if sub == nil {
		return
	}
	// DebugTraceF("当前用户人数：%d", len(w.subscribers))
	subTemp := w.subscribers.findByID(sub.GetID())
	if subTemp == nil {
		// sub.setOnline()
		w.subscribers = append(w.subscribers, sub)
		// DebugTraceF("观察者上线 %s", sub.ID)
		// } else {
		// subTemp.setOnline()
		// subTemp.Conn = sub.Conn
		// event := NewSysEvent(getSysEventDefValue("user_online"), subTemp)
		// triggerSysEvent(event)
		w.triggerSysEvent(sys_event_user_online, sub)
	}
	DebugTraceF("上线后，最新用户人数：%d", len(w.subscribers))
	// for _, s := range g_subscribers {
	// 	if s.ID == user.ID {
	// 		s.setOnline()
	// 		s.Conn = user.Conn
	// 	}
	// }
}

// func isUserExist(name string) bool {
// 	// for _, sub := range subscribers {
// 	for ele := subscribers.Front(); ele != nil; ele = ele.Next() {
// 		sub := ele.Value.(*Subscriber)
// 		if sub.Name == name {
// 			return true
// 		}
// 	}
// 	return false
// }
