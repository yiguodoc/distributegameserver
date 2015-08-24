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

// type eventProcessor func(*SysEvent)

//对系统事件的订阅
type SysEventSubscribe struct {
	code             int
	eventSubscribers []eventProcessor
	// eventSubscribers map[string]eventProcessor
}

func (s *SysEventSubscribe) addProcessor(process eventProcessor) {
	for _, f := range s.eventSubscribers {
		if &f == &process {
			return
		}
	}
	s.eventSubscribers = append(s.eventSubscribers, process)
}
func (s *SysEventSubscribe) removeSubscriber(process eventProcessor) {
	l := []eventProcessor{}
	for _, f := range s.eventSubscribers {
		if &f != &process {
			l = append(l, f)
		}
	}
	s.eventSubscribers = l
}
func NewSysEventSubscribe(code interface{}) *SysEventSubscribe {
	return &SysEventSubscribe{
		code:             code.(int),
		eventSubscribers: []eventProcessor{},
		// eventSubscribers: make(map[string]func(*SysEvent)),
	}
}

type SysEventSubscribeList map[int]*SysEventSubscribe

// type SysEventSubscribeList []*SysEventSubscribe

func (s SysEventSubscribeList) contains(code interface{}) bool {
	_, ok := s[code.(int)]
	return ok
}

func (s SysEventSubscribeList) notifyEventSubscribers(code, para interface{}) {
	// func (s SysEventSubscribeList) notifyEventSubscribers(sysEvent *SysEvent) {
	// sub := s.find(code)
	// sub := s.find(sysEvent.eventCode)
	if s.contains(code.(int)) {
		sub := s[code.(int)]
		// DebugTraceF("notifyEventSubscribers %d => %d  ", code, len(sub.eventSubscribers))
		for _, f := range sub.eventSubscribers {
			go f(code, para)
		}
	} else {
		DebugSysF("can NOT find code %d", code)
	}

}
func (s SysEventSubscribeList) removeEventSubscriber(code interface{}, process eventProcessor) {
	// func (s SysEventSubscribeList) removeEventSubscriber(code sysEventCode, flag string) {
	if s.contains(code.(int)) {
		s[code.(int)].removeSubscriber(process)
	}

}

// func (s SysEventSubscribeList) addEventSubscriber(code sysEventCode, f eventProcessor) {
// func (s SysEventSubscribeList) addEventSubscriber(flag string, f eventProcessor, codes ...sysEventCode) {
func (s SysEventSubscribeList) addEventSubscriber(f eventProcessor, codes ...interface{}) {
	for _, code := range codes {
		// sub := s.find(code.(int))
		if s.contains(code) {
			s[code.(int)].addProcessor(f)
			// sub.eventSubscribers[flag] = f
		} else {
			sub := NewSysEventSubscribe(code)
			sub.addProcessor(f)
			s[code.(int)] = sub
		}
	}
}
