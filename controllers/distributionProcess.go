package controllers

import (
	// "errors"
	"fmt"
	"time"
)

type ProcessNode struct {
	chanEvent chan *SysEvent
}

func NewProcessNode() *ProcessNode {
	return &ProcessNode{
		chanEvent: make(chan *SysEvent, 128),
	}
}
func (po *ProcessNode) init() {
	po.initProcessListening()
}
func (po *ProcessNode) acceptEvent(event *SysEvent) {
	po.chanEvent <- event
}
func (po *ProcessNode) initProcessListening() {
	go func() {
		for {
			event := <-po.chanEvent
			switch event.eventCode {

			}
		}
	}()

}
