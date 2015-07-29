package controllers

import (
// "container/list"
// "github.com/astaxie/beego"
// "github.com/gorilla/websocket"
// "time"
// "encoding/json"
// "strings"
// "fmt"
)

//对系统事件的订阅
type SysEventSubscribe struct {
	eventType        sysEventCode
	eventSubscribers map[string]func(*SysEvent)
}

func NewSysEventSubscribe(eventType sysEventCode) *SysEventSubscribe {
	return &SysEventSubscribe{
		eventType:        eventType,
		eventSubscribers: make(map[string]func(*SysEvent)),
	}
}

type SysEventSubscribeList []*SysEventSubscribe

func (s SysEventSubscribeList) find(eventType sysEventCode) *SysEventSubscribe {
	for _, sub := range s {
		if sub.eventType == eventType {
			return sub
		}
	}
	return nil
}
func (s SysEventSubscribeList) notifyEventSubscribers(sysEvent *SysEvent) {
	sub := s.find(sysEvent.eventCode)
	if sub != nil {
		for _, f := range sub.eventSubscribers {
			go f(sysEvent)
		}
	}
}
func (s SysEventSubscribeList) removeEventSubscriber(eventType sysEventCode, flag string) {
	sub := s.find(eventType)
	if sub != nil {
		delete(sub.eventSubscribers, flag)
	}
}
func (s SysEventSubscribeList) addEventSubscriber(code sysEventCode, flag string, f func(*SysEvent)) {
	// func (s SysEventSubscribeList) addEventSubscriber(eventName string, flag string, f func(*SysEvent)) {
	// code := getSysEventDefValue(eventName)
	sub := s.find(code)
	if sub != nil {
		sub.eventSubscribers[flag] = f
	}
}
