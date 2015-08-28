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
		DebugTraceF("客户端事件定义：%3d : %s", code, code.name())
	}
}
func getClientMessageTypeCodeList() (l []ClientMessageTypeCode) {
	f := func(start, stop int) (l []ClientMessageTypeCode) {
		for i := start; i <= stop; i++ {
			l = append(l, ClientMessageTypeCode(i))
		}
		return
	}
	l = append(append(l, f(pro_min+1, int(pro_max)-1)...), f(pro_2c_min+1, int(pro_2c_max)-1)...)
	// for i := pro_min + 1; i < int(pro_max); i++ {
	// 	l = append(l, ClientMessageTypeCode(i))
	// }
	// for i := pro_2c_min + 1; i < int(pro_2c_max); i++ {
	// 	l = append(l, ClientMessageTypeCode(i))
	// }
	return
}

type ClientMessageTypeCode int

type MessageWithClientHandler func(*MessageWithClient)
type MessageWithClientHandlerGenerator (func(interface{}) MessageWithClientHandler)

// type MessageWithClientHandlerGenerator (func(*DistributorProcessUnit) MessageWithClientHandler)
type ProHandlerGeneratorMap map[ClientMessageTypeCode]MessageWithClientHandlerGenerator
type ProHandlerMap map[ClientMessageTypeCode]MessageWithClientHandler

func (p ProHandlerGeneratorMap) generateHandlerMap(codes []ClientMessageTypeCode, obj interface{}) ProHandlerMap {
	m := make(ProHandlerMap)
	for _, code := range codes {
		if generator, ok := p[code]; ok {
			m[code] = generator(obj)
		} else {
			DebugSysF("未定义事件 %s 的处理函数", code.name())
		}
	}
	return m
}

//*******************************************************************************************************************
var handler_map = ProHandlerGeneratorMap{
	pro_on_line:                   pro_on_line_handlerGenerator,
	pro_off_line:                  pro_off_line_handlerGenerator,
	pro_prepared_for_select_order: pro_prepared_for_select_order_handlerGenerator,
	pro_game_start:                pro_game_start_handlerGenerator,
	pro_order_select_response:     pro_order_select_response_handlerGenerator,
	pro_reset_destination_request: pro_reset_destination_request_handlerGenerator,
	pro_game_time_tick:            pro_game_time_tick_handlerGenerator,
	pro_change_state_request:      pro_change_state_request_handlerGenerator,
	pro_sign_order_request:        pro_sign_order_request_handlerGenerator,
	pro_move_from_node_to_route:   pro_move_from_node_to_route_handlerGenerator,
	pro_move_from_route_to_node:   pro_move_from_route_to_node_handlerGenerator,
	pro_distributor_info_request:  pro_distributor_info_request_handlerGenerator,
}

//后端与前端交互的消息类型的定义
var (
	pro_min                                              = 0
	pro_on_line                    ClientMessageTypeCode = 1  //配送员上线
	pro_off_line                   ClientMessageTypeCode = 2  //配送员下线
	pro_prepared_for_select_order  ClientMessageTypeCode = 3  //配送员准备好抢订单了,当所有配送员都准备好之后，就可以分发订单了
	pro_game_start                 ClientMessageTypeCode = 4  //第一次向配送员分发订单，自带倒数
	pro_order_select_response      ClientMessageTypeCode = 5  //单个配送员对订单的选择，向服务端提交
	pro_distribution_start_request ClientMessageTypeCode = 6  //配送员向服务端请求开始订单的配送
	pro_reset_destination_request  ClientMessageTypeCode = 7  //配送员向服务端请求重置目标点
	pro_change_state_request       ClientMessageTypeCode = 8  //配送员向服务端请求改变运行状态，0 停止  1 运行
	pro_timer_count_down           ClientMessageTypeCode = 9  //倒计时
	pro_sign_order_request         ClientMessageTypeCode = 10 //
	pro_move_from_node_to_route    ClientMessageTypeCode = 11 //配送员从节点上路
	pro_move_from_route_to_node    ClientMessageTypeCode = 12 //
	pro_game_time_tick             ClientMessageTypeCode = 13 //系统时间流逝出发
	pro_distributor_info_request   ClientMessageTypeCode = 14 //系统时间流逝出发
	pro_max                                              = 15

	pro_2c_min                                                       = 400
	pro_2c_message_broadcast                   ClientMessageTypeCode = 401 //向配送员广播消息
	pro_2c_order_distribution_proposal         ClientMessageTypeCode = 402 //向配送员分发订单
	pro_2c_order_select_result                 ClientMessageTypeCode = 403 //订单最终归属结果，向客户端推送
	pro_2c_map_data                            ClientMessageTypeCode = 404 //向配送员发送地图信息
	pro_2c_reset_destination                   ClientMessageTypeCode = 405 //服务端通知配送员重置了目标点
	pro_2c_change_state                        ClientMessageTypeCode = 406 //服务端通知配送员改变运行状态，0 停止  1 运行
	pro_2c_move_to_new_position                ClientMessageTypeCode = 407 //通知客户端新位置
	pro_2c_reach_route_node                    ClientMessageTypeCode = 408 //到达一个路径节点
	pro_2c_sign_order                          ClientMessageTypeCode = 409 //
	pro_2c_distributor_info                    ClientMessageTypeCode = 410 //
	pro_2c_order_full                          ClientMessageTypeCode = 411 //订单满载
	pro_2c_all_prepared_4_select_order         ClientMessageTypeCode = 412
	pro_2c_message_broadcast_before_game_start ClientMessageTypeCode = 413
	pro_2c_move_from_node                      ClientMessageTypeCode = 414
	pro_2c_max                                 ClientMessageTypeCode = 415
)

func (c ClientMessageTypeCode) name() (s string) {
	switch c {
	case pro_2c_move_from_node:
		s = "pro_2c_move_from_node"
	case pro_2c_message_broadcast_before_game_start:
		s = "pro_2c_message_broadcast_before_game_start"
	case pro_2c_all_prepared_4_select_order:
		s = "pro_2c_all_prepared_4_select_order"
	case pro_2c_order_full:
		s = "pro_2c_order_full"
	case pro_2c_distributor_info:
		s = "pro_2c_distributor_info"
	case pro_game_time_tick:
		s = "pro_game_time_tick"
	case pro_2c_order_distribution_proposal:
		s = "pro_2c_order_distribution_proposal"
	case pro_order_select_response:
		s = "pro_order_select_response"
	case pro_timer_count_down:
		s = "pro_timer_count_down"
	case pro_distributor_info_request:
		s = "pro_distributor_info_request"
	case pro_2c_message_broadcast:
		s = "pro_2c_message_broadcast"
	case pro_2c_order_select_result:
		s = "pro_2c_order_select_result"
	case pro_on_line:
		s = "pro_on_line"
	case pro_off_line:
		s = "pro_off_line"
	case pro_prepared_for_select_order:
		s = "pro_prepared_for_select_order"
	case pro_2c_map_data:
		s = "pro_2c_map_data"
	case pro_distribution_start_request:
		s = "pro_distribution_start_request"
	case pro_2c_reset_destination:
		s = "pro_2c_reset_destination"
	case pro_reset_destination_request:
		s = "pro_reset_destination_request"
	case pro_change_state_request:
		s = "pro_change_state_request"
	case pro_2c_change_state:
		s = "pro_2c_change_state"
	case pro_game_start:
		s = "pro_game_start"
	case pro_2c_move_to_new_position:
		s = "pro_2c_move_to_new_position"
	case pro_2c_reach_route_node:
		s = "pro_2c_reach_route_node"
	case pro_sign_order_request:
		s = "pro_sign_order_request"
	case pro_2c_sign_order:
		s = "pro_2c_sign_order"
	case pro_move_from_node_to_route:
		s = "pro_move_from_node_to_route"
	case pro_move_from_route_to_node:
		s = "pro_move_from_route_to_node"
	default:
		DebugSysF("客户端事件(%3d)定义描述不完全", int(c))
		// if (c) < pro_max {
		// 	panic(fmt.Sprintf("客户端事件(%3d)定义描述不完全", int(c)))
		// }
	}
	return
}
func NewMessageWithClient(code ClientMessageTypeCode, targetID string, data interface{}) *MessageWithClient {
	return &MessageWithClient{
		MessageType: code,
		TargetID:    targetID,
		Data:        data,
	}
}

type MessageWithClient struct {
	MessageType ClientMessageTypeCode
	TargetID    string
	Data        interface{}
	// Data        string
}

func (m *MessageWithClient) String() string {
	return fmt.Sprintf("type: %s(%d) TargetID: %s data: %s", m.MessageType.name(), m.MessageType, m.TargetID, m.Data)
}

// type MessageFromClient struct {
// 	MessageType int
// 	Data        string
// }
