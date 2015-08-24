package controllers

import (
	// "errors"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego"
	"github.com/ungerik/go-dry"
	"os"
)

var viewerCount = 1

func getViewerID() string {
	viewerCount++
	return fmt.Sprintf("viewerndlfejqwrjlfjfww953957392%d", viewerCount)
}

type ResponseMsg struct {
	Code    int
	Message string
}

func NewResponseMsg(code int, msg ...string) *ResponseMsg {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return &ResponseMsg{
		Code:    code,
		Message: message,
	}
}

type MainController struct {
	beego.Controller
}

type MapData struct {
	Points PositionList
	Lines  LineList
}

func (m *MainController) Index() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_UnitCenter.distributors.findOne(func(o interface{}) bool { return o.(DataWithID).GetID() == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)

	m.TplNames = "index.tpl"
}
func (m *MainController) Login() {
	m.TplNames = "login.tpl"
}

//配送页面
func (m *MainController) DistributionIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_UnitCenter.distributors.findOne(func(o interface{}) bool { return o.(DataWithID).GetID() == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)
	m.TplNames = "distribution.tpl"
}

//载入地图数据
func loadMapData() *MapData {
	var mapData MapData
	file := "mapdata/data.toml"
	if dry.FileExists(file) == false {
		DebugInfoF("地图文件 %s 不存在", file)
		return nil
	}
	_, err := toml.DecodeFile(file, &mapData)
	if err != nil {
		DebugMustF("载入地图数据时出错：%s", err)
		return nil
	} else {
		DebugInfoF("载入地图数据成功，统计：%d 个关键点  %d 条路径", len(mapData.Points), len(mapData.Lines))
		DebugPrintList_Info(mapData.Points)
		DebugPrintList_Info(mapData.Lines)
	}
	return &mapData
}

//上传编辑后的地图数据
func (m *MainController) UploadMapData() {
	response := NewResponseMsg(0)
	defer func() {
		m.Data["json"] = response
		m.ServeJson()
	}()
	values := m.Input()
	if value, ok := values["data"]; !ok {
		m.Data["json"] = NewResponseMsg(1, "地图数据格式异常")
		DebugMust("地图数据格式异常")
	} else {
		if len(value) <= 0 {
			DebugMust("没有地图数据上传")
			m.Data["json"] = NewResponseMsg(1, "没有地图数据上传")
		} else {
			rawData := values["data"][0]
			fmt.Println(rawData)
			var mapData MapData
			err := json.Unmarshal([]byte(rawData), &mapData)
			if err != nil {
				DebugMustF("解析上传地图数据时出错：%s", err)
				m.Data["json"] = NewResponseMsg(1, "解析上传地图数据时出错")
			} else {
				// fmt.Println(mapData)
				DebugInfoF("接收到上传的地图数据，统计：%d 个关键点  %d 条路径", len(mapData.Points), len(mapData.Lines))
				DebugPrintList_Info(mapData.Points)
				DebugPrintList_Info(mapData.Lines)
				fileMapData, err := os.Create("./mapdata/data.toml")
				if err != nil {
					DebugMustF("创建地图文件出错：%s", err)
				} else {
					defer fileMapData.Close()
					err = toml.NewEncoder(fileMapData).Encode(mapData)
					if err != nil {
						DebugMustF("保存地图数据到文件时出错：%s", err)
					}
				}
			}
		}
	}
}

//查询输出地图数据
func (m *MainController) MapData() {
	m.Data["json"] = g_UnitCenter.mapData
	m.ServeJson()
}

//地图编辑页面
func (m *MainController) AddressEditIndex() {
	m.TplNames = "addressEdit.tpl"
}
func (m *MainController) OrderDistributeIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_UnitCenter.distributors.findOne(func(o interface{}) bool { return o.(DataWithID).GetID() == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)
	m.TplNames = "orderDistribute.tpl"
}
func setProData(m *MainController) {
	codes := getClientMessageTypeCodeList()
	for _, code := range codes {
		m.Data[code.name()] = code
	}
}
func (m *MainController) Orders() {
	id := m.GetString("id")
	if len(id) <= 0 {
		m.Data["json"] = g_UnitCenter.orders
	} else {
		d := g_UnitCenter.orders.findOne(func(o interface{}) bool { return o.(*Order).ID == id })
		// d := g_UnitCenter.orders.findByID(id)
		if d == nil {
			m.Data["json"] = OrderList{}
		} else {
			m.Data["json"] = OrderList{d}
		}
	}
	m.ServeJson()

}
func (m *MainController) Distributors() {
	id := m.GetString("id")
	if len(id) <= 0 {
		m.Data["json"] = g_UnitCenter.distributors
	} else {
		d := g_UnitCenter.distributors.findOne(func(o interface{}) bool { return o.(DataWithID).GetID() == id })
		// d := g_UnitCenter.distributors.find(id)
		if d == nil {
			m.Data["json"] = DistributorList{}
		} else {
			m.Data["json"] = DistributorList{d}
		}
	}
	m.ServeJson()
}
func (m *MainController) UserListIndex() {
	m.TplNames = "userListIndex.tpl"
}

func (m *MainController) ViewerIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	m.Data["ID"] = getViewerID()
	m.TplNames = "ViewerIndex.tpl"
}

/*
//启动分发控制器后，控制器开始分发订单，并将分发结果通过消息发送给维护在线配送员终端的列表
//配送终端接收到消息后，开始订单的接受或者拒绝
//正常情况，订单被某配送员抢中，该选择作为消息，发送到分发中心 （极端情况：所有的配送员选择拒绝，订单随机分配）
//分发中心接受两种消息反馈：超时和配送员客户端的消息。超时表示所有分发的订单没有别任何配送员接收（之前需要确认不是因为客户端掉线的原因）
// 接收到客户端的反馈可以开始下次分发，直至分发完毕，自行停止运作



//订单分发的控制器（和路由无关）
//控制订单分发的节奏
type OrderDistributeController struct {
	chanRunning         chan bool
	chanStopRunning     chan bool
	OrdersDistributed   OrderList       //已经分配的订单
	OrdersUndistributed OrderList       //尚未分配的订单
	Distributors        DistributorList //所有配送员
	// chanDistributionResponse chan *OrderDistribution //接收反馈消息

}

// //接收订单分配的反馈
// //如果反馈被接受，则成为分配结果，相应地订单也会发生变化
// func (o *OrderDistributeController) AcceptDistributionResponse(od *OrderDistribution) {
// 	o.chanDistributionResponse <- od
// }

//向配送员发送要分配的订单信息
func (o *OrderDistributeController) distributeProposal() {
	list, err := o.createDistributionProposal()
	if err != nil {

	}
	chanDistributionResponse := make(chan *OrderDistribution) //接收反馈消息
	go func() {                                               //开始监听配送员的反馈，只接收一个接收订单的信息
		for {
			od := <-chanDistributionResponse
			distributor := o.Distributors.find(od.DistributorID) //首先确保配送员满足订单分配条件，当前条件是已分配的订单未达到最大可接收数量
			if distributor == nil {
				// return nil, ERR_NO_SUCH_DISTRIBUTOR
				continue
			}
			if distributor.full() {
				// return nil, ERR_DISTRIBUTOR_FULL
				continue
			}
			//其次订单满足分配条件，当前的条件是尚未分配
			order := o.OrdersUndistributed.findByID(od.OrderID)
			if order == nil {
				// return nil, ERR_CANNOT_FIND_ORDER_FOR_DISTRIBUTION
				continue
			}
			//确定结果
			o.OrdersUndistributed.remove(order)
			o.OrdersDistributed = append(o.OrdersDistributed, order)
			distributor.acceptOrder(order)
			break
		}
	}()
	notifyDistributorsOrder(list)
}

//只是生成一个分配建议，不是最终的分配结果
func (o *OrderDistributeController) createDistributionProposal() (list OrderDistributionList, err error) {
	if len(o.OrdersUndistributed) <= 0 {
		err = ERR_NO_ENOUGH_ORDER_TO_DISTRIBUTE
		return
	}
	distributorsNotFull := o.Distributors.notFull()
	if len(o.OrdersUndistributed) < len(distributorsNotFull) {
		DebugMustF("There is %d orders and %d distributors", len(o.OrdersUndistributed), len(distributorsNotFull))
		err = ERR_NO_ENOUGH_ORDER
		return
	}
	order := o.OrdersUndistributed[0]
	for _, distributor := range o.Distributors {
		list = list.add(NewOrderDistribution(order.ID, distributor.ID))
	}
	return
}

// //停止监听运行
// func (o *OrderDistributeController) Stop() {
// 	o.chanStopRunning <- true
// }

// //新一轮分发
// func (o *OrderDistributeController) Newloop() {
// 	o.chanRunning <- true
// }

// //启动控制器，开始分发订单
// func (o *OrderDistributeController) Start() {
// 	go func() {
// 		for {
// 			select {
// 			case <-o.chanRunning: //下一轮的订单分发
// 				// if usersOnline() {
// 				o.distributeProposal()
// 				// }
// 			case result := <-o.chanStopRunning: //true,then stop
// 				if result == true {
// 					break
// 				}
// 			}

// 		}
// 	}()
// }
*/
