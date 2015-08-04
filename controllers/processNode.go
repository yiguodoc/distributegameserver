package controllers

import (
// "errors"
// "fmt"
// "time"
)

//实现了事件分发处理的功能，类似于路由

type ProcessNode struct {
	chanEvent chan *SysEvent
	processor eventProcessor
}

func NewProcessNode(processor eventProcessor) *ProcessNode {
	return &ProcessNode{
		chanEvent: make(chan *SysEvent),
		processor: processor,
	}
}
func (po *ProcessNode) NodeEventProcessor() eventProcessor {
	f := func(code, data interface{}) {
		// f := func(event *SysEvent) {
		// po.chanEvent <- event
	}
	return f
}
func (po *ProcessNode) init() *ProcessNode {
	po.initProcessListening()
	return po
}
func (po *ProcessNode) acceptEvent(event *SysEvent) *ProcessNode {
	po.chanEvent <- event
	return po
}
func (po *ProcessNode) initProcessListening() *ProcessNode {
	go func() {
		for {
			// event := <-po.chanEvent
			// if po.processor != nil {
			// 	po.processor(event)
			// }
		}
	}()
	return po
}
