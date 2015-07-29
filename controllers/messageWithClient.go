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

//后端与前端交互的消息类型的定义
const (
	pro_order_distribution_proposal = 1  //配送员可选订单分发
	pro_order_select_response       = 2  //配送员抢订单结果
	pro_timer_count_down            = 3  //倒计时
	pro_begin_select_order          = 4  //通知客户端可以开始抢订单了
	pro_message_broadcast           = 5  //向配送员广播消息
	pro_order_distribute_result     = 6  //订单最终归属结果，向客户端推送
	pro_distributor_on_line         = 7  //配送员上线
	pro_distributor_off_line        = 8  //配送员下线
	pro_distributor_prepared        = 9  //配送员准备好抢订单了
	pro_distribution_prepared       = 10 //配送员订单满载，可以开始配送了
)

type MessageWithClient struct {
	MessageType int
	Data        interface{}
}

func (m *MessageWithClient) String() string {
	return fmt.Sprintf("type: %-2d  data: %s", m.MessageType, m.Data)
}
