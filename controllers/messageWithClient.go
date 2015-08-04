package controllers

import (
	// "time"
	// "encoding/json"
	// "strings"
	// "errors"
	"fmt"
)

func clientMessageTypeCodeCheck() {
	l := getClientMessageTypeCodeList()
	for _, code := range l {
		DebugInfoF("客户端事件定义：%3d : %s", code, code.name())
	}
}
func getClientMessageTypeCodeList() (l []ClientMessageTypeCode) {
	for i := 1; i < int(pro_max); i++ {
		l = append(l, ClientMessageTypeCode(i))
	}
	return
}

// func (c ClientMessageTypeCode) mapToSysEventCode() (code sysEventCode, err error) {
// 	switch c {
// 	case pro_order_select_response:
// 		code = sys_event_order_select_response

// 	case pro_distributor_prepared_select_order:
// 		code = sys_event_select_order_prepared

// 	case pro_reset_destination_request:
// 		code = sys_event_reset_destination_request

// 	case pro_change_state_request:
// 		code = sys_event_change_state_request
// 	default:
// 		err = errors.New("没有对应的系统事件")
// 	}
// 	return
// }

type ClientMessageTypeCode int

//后端与前端交互的消息类型的定义
var (
	pro_message_broadcast                 ClientMessageTypeCode = 1  //向配送员广播消息
	pro_on_line                           ClientMessageTypeCode = 2  //配送员上线
	pro_off_line                          ClientMessageTypeCode = 3  //配送员下线
	pro_prepared_for_select_order         ClientMessageTypeCode = 4  //配送员准备好抢订单了,当所有配送员都准备好之后，就可以分发订单了
	pro_begin_to_select_order             ClientMessageTypeCode = 5  //通知客户端可以开始抢订单了
	pro_order_distribution_proposal       ClientMessageTypeCode = 6  //向配送员分发订单
	pro_order_distribution_proposal_first ClientMessageTypeCode = 7  //第一次向配送员分发订单，自带倒数
	pro_order_select_response             ClientMessageTypeCode = 8  //单个配送员对订单的选择，向服务端提交
	pro_order_select_result               ClientMessageTypeCode = 9  //订单最终归属结果，向客户端推送
	pro_distribution_prepared             ClientMessageTypeCode = 10 //向配送员发送，可以开始配送，附带的数据包括配送员的所有信息
	pro_distribution_start_request        ClientMessageTypeCode = 11 //配送员向服务端请求开始订单的配送
	pro_reset_destination                 ClientMessageTypeCode = 12 //服务端通知配送员重置了目标点
	pro_reset_destination_request         ClientMessageTypeCode = 13 //配送员向服务端请求重置目标点
	pro_change_state_request              ClientMessageTypeCode = 14 //配送员向服务端请求改变运行状态，0 停止  1 运行
	pro_change_state                      ClientMessageTypeCode = 15 //服务端通知配送员改变运行状态，0 停止  1 运行
	pro_timer_count_down                  ClientMessageTypeCode = 16 //倒计时
	pro_max                               ClientMessageTypeCode = 17
)

func (c ClientMessageTypeCode) name() (s string) {
	switch c {
	case pro_order_distribution_proposal:
		s = "pro_order_distribution_proposal"
	case pro_order_select_response:
		s = "pro_order_select_response"
	case pro_timer_count_down:
		s = "pro_timer_count_down"
	case pro_begin_to_select_order:
		s = "pro_begin_to_select_order"
	case pro_message_broadcast:
		s = "pro_message_broadcast"
	case pro_order_select_result:
		s = "pro_order_select_result"
	case pro_on_line:
		s = "pro_on_line"
	case pro_off_line:
		s = "pro_off_line"
	case pro_prepared_for_select_order:
		s = "pro_prepared_for_select_order"
	case pro_distribution_prepared:
		s = "pro_distribution_prepared"
	case pro_distribution_start_request:
		s = "pro_distribution_start_request"
	case pro_reset_destination:
		s = "pro_reset_destination"
	case pro_reset_destination_request:
		s = "pro_reset_destination_request"
	case pro_change_state_request:
		s = "pro_change_state_request"
	case pro_change_state:
		s = "pro_change_state"
	case pro_order_distribution_proposal_first:
		s = "pro_order_distribution_proposal_first"
	default:
		if (c) < pro_max {
			panic(fmt.Sprintf("客户端事件(%3d)定义描述不完全", int(c)))
		}
	}
	return
}
func NewMessageWithClient(code ClientMessageTypeCode, data interface{}) *MessageWithClient {
	return &MessageWithClient{
		MessageType: code,
		Data:        data,
	}
}

type MessageWithClient struct {
	MessageType ClientMessageTypeCode
	Data        interface{}
	// Data        string
}

func (m *MessageWithClient) String() string {
	return fmt.Sprintf("type: %-2d  data: %s", m.MessageType, m.Data)
}

// type MessageFromClient struct {
// 	MessageType int
// 	Data        string
// }
