package controllers

import (
// "container/list"
// "github.com/astaxie/beego"
// "github.com/gorilla/websocket"
// "time"
// "strings"
// "fmt"
)

var (
// subscriber_type_distributor = 1 //distributor 配送员 1
// subscriber_type_viewer      = 2 // viewer 观察者 2
)

// type Subscriber struct {
// 	ID       string
// 	Conn     *websocket.Conn `json:"-"` // Only for WebSocket users; otherwise nil.
// 	IsOnline bool            `json:"-"`
// 	// myType   int
// 	// Name string
// }

// func (s *Subscriber) String() string {
// 	return fmt.Sprintf("ID: %-10s  isOnline: %t  ", s.ID, s.IsOnline)
// }

// func (s *Subscriber) setOnline() {
// 	s.IsOnline = true
// }
// func (s *Subscriber) setOffline() {
// 	s.IsOnline = false
// }
// func NewSubscriber(id string, conn *websocket.Conn) *Subscriber {
// 	return &Subscriber{
// 		ID:   id,
// 		Conn: conn,
// 		// myType: subscriberType,
// 	}
// }

// type SubscriberList []*Subscriber

// func (s SubscriberList) ListName() string {
// 	return "在线人员"
// }
// func (s SubscriberList) InfoList() (l []string) {
// 	for _, sub := range s {
// 		if sub.IsOnline == true {
// 			l = append(l, sub.ID)
// 		}
// 	}
// 	return
// }

// func (s SubscriberList) allOnline() bool {
// 	return len(s.onlineSubscribers()) >= len(s)
// }
// func (s SubscriberList) allDistributorsOnline() bool {
// 	list := s.filtedSubscribers(subscriber_type_distributor)
// 	return len(list.onlineSubscribers()) >= len(list)
// }
// func (s SubscriberList) onlineCountGreaterThan(count int) bool {
// 	return len(s.onlineSubscribers()) > count
// }

//过滤出特定类型的成员
// func (s SubscriberList) filtedSubscribers(subscriber_type int) (list SubscriberList) {
// 	for _, sub := range s {
// 		if sub.myType == subscriber_type {
// 			list = append(list, sub)
// 		}
// 	}
// 	return
// }
// func (s SubscriberList) onlineSubscribers() (list SubscriberList) {
// 	for _, sub := range s {
// 		if sub.IsOnline == true {
// 			list = append(list, sub)
// 		}
// 	}
// 	return
// }
// func (s SubscriberList) findByID(id string) *Subscriber {
// 	for _, sub := range s {
// 		if sub.ID == id {
// 			return sub
// 		}
// 	}
// 	return nil
// }
// func (s SubscriberList) remove(id string) (list SubscriberList) {
// 	for _, sub := range s {
// 		if sub.ID != id {
// 			list = append(list, sub)
// 		}
// 	}
// 	return
// }
