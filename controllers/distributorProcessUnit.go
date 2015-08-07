package controllers

import (
	// "errors"
	"fmt"
	// "encoding/json"
	"math"
	"time"
)

type DistributorProcessUnitList map[string]*DistributorProcessUnit

type DistributorProcessUnitCenter struct {
	units     DistributorProcessUnitList
	chanEvent (chan *MessageWithClient)
	// chanResult chan bool //返回执行的结果
}

func NewDistributorProcessUnitCenter() *DistributorProcessUnitCenter {
	return &DistributorProcessUnitCenter{
		units:     DistributorProcessUnitList{},
		chanEvent: make(chan *MessageWithClient),
	}
}

func (dpc *DistributorProcessUnitCenter) start() *DistributorProcessUnitCenter {
	go func() {
		for {
			select {
			case msg := <-dpc.chanEvent:
				switch msg.MessageType {
				case pro_order_distribution_proposal_first: //全体倒计时
					g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, "配送员全部准备完毕")
					g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, "一大波订单即将到来")
					//倒计时
					timer := time.Tick(1 * time.Second)
					count := 3
					// DebugInfo("start timer...")
					for {
						<-timer
						DebugTraceF("timer count : %d", count)
						if count <= 0 {
							break
						}
						g_room_distributor.broadcastMsgToSubscribers(pro_timer_count_down, count)
						count--
					}
					if len(g_distributors.notFull()) > 0 {
						broadOrderSelectProposal()
					}

				case pro_order_select_response:
					m := msg.Data.(map[string]interface{})
					if list, err := mappedValue(m).Getter("OrderID", "DistributorID"); err == nil {
						orderID := list[0].(string)
						distributorID := list[1].(string)
						if err := disposeOrderSelectResponse(orderID, distributorID); err != nil {
							DebugInfoF("处理订单分配时的信息提示：%s", err)
						} else {
							//将分配结果通知到各方，包括获得订单的客户端、群通知，并引发分配结果事件，使得观察者也可以得到通知
							distributionResult := NewOrderDistribution(orderID, distributorID)
							distributor := g_distributors.find(distributorID)
							g_room_distributor.sendMsgToSpecialSubscriber(distributorID, pro_order_select_result, distributionResult)

							msg := fmt.Sprintf("订单[%s]已经由配送员[%s]选定", distributionResult.OrderID, distributor.Name)
							g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, msg)
							DebugInfo(msg)

							if distributor.full() == true { //配送员的订单满载了
								g_room_distributor.broadcastMsgToSubscribers(pro_message_broadcast, fmt.Sprintf("配送员 %s 订单满载", distributor.Name))
								g_room_distributor.broadcastMsgToSubscribers(pro_distribution_prepared, distributor.ID)
								distributor.setCheckPoint(checkpoint_flag_order_distribute)
								// triggerSysEvent(NewSysEvent(sys_event_distribution_prepared, distributor))
							}
							if len(g_distributors.notFull()) > 0 {
								broadOrderSelectProposal()
							}
						}
					} else {
						DebugMustF("客户端数据格式错误: %s", err)
					}
				default:
				}
			}
		}
	}()
	return dpc
}
func (dpc *DistributorProcessUnitCenter) allUnitsProcess(msg *MessageWithClient) {
	// func (dpc *DistributorProcessUnitCenter) allUnitsProcess(data interface{}) {
	dpc.chanEvent <- msg
}

//id : Distributor ID
func (dpc *DistributorProcessUnitCenter) singleUnitprocess(id string, msg *MessageWithClient) {
	go func() {
		// var msg MessageWithClient
		// err := json.Unmarshal(data.([]byte), &msg)
		// if err != nil {
		// 	DebugSysF("解析数据出错：%s", err)
		// 	continue
		// }
		if unit, ok := dpc.units[id]; ok {
			unit.process(msg)
		} else {
			DebugSysF("没有找到可以处理的配送单元，系统异常")
		}
	}()

}
func (dpc *DistributorProcessUnitCenter) newUnit(distributor *Distributor, processors ...(func(*MessageWithClient, *DistributorProcessUnit))) {
	// func (dpc *DistributorProcessUnitCenter) newUnit(distributor *Distributor, processors ...(func(ClientMessageTypeCode, interface{}, *DistributorProcessUnit))) {
	if _, ok := dpc.units[distributor.ID]; ok {
		DebugInfoF("配送处理单元 %s 重复添加", distributor.ID)
	} else {
		unit := &DistributorProcessUnit{
			center:      dpc,
			distributor: distributor,
			chanStop:    make(chan bool),
			processors:  processors,
		}
		dpc.units[distributor.ID] = unit
		unit.start()
	}
}
func (dpc *DistributorProcessUnitCenter) removeUnit(id string) {
	if u, ok := dpc.units[id]; ok {
		u.stop()
		delete(dpc.units, id)
	} else {
		DebugSysF("配送处理单元 %s 不存在，无法移除", id)
	}
}

type DistributorProcessUnit struct {
	center      *DistributorProcessUnitCenter
	processors  []func(*MessageWithClient, *DistributorProcessUnit)
	chanEvent   (chan *MessageWithClient)
	distributor *Distributor
	chanStop    chan bool
	// chanEvent   (chan []byte)
}

func (u *DistributorProcessUnit) process(data *MessageWithClient) {
	u.chanEvent <- data
}
func (u *DistributorProcessUnit) stop() {
	u.chanStop <- true
}
func (u *DistributorProcessUnit) start() {
	u.chanEvent = make(chan *MessageWithClient)
	f := func() {
		timer := time.Tick(1 * time.Second) //计时器功能
		for {
			select {
			case <-u.chanStop:
				break
			case msg := <-u.chanEvent:
				switch msg.MessageType {
				default:
					for _, processor := range u.processors {
						processor(msg, u)
					}
				}
			case <-timer:
				distributor := u.distributor
				//----------------------------------------------------------------------------
				//计算行走的坐标位置
				if u.distributor.CurrentSpeed > 0 {
					if u.distributor.StartPos != nil && u.distributor.DestPos != nil {
						if u.distributor.StartPos.equals(u.distributor.DestPos) == false {
							DebugInfoF("配送员 %s 运行路线 %s => %s", u.distributor.Name, u.distributor.StartPos.SimpleString(), u.distributor.DestPos.SimpleString())
							totalTime := u.distributor.Distance * 60 * 60 / (u.distributor.CurrentSpeed * 1000) / realityToSystemTimeRatio //系统中运行路程所花费的时间
							totalFrames := totalTime / timePerFrame                                                                        //一共大约这么多帧就可以走完
							//使用绝对值差距大的作为分片的标准
							totalLng := u.distributor.DestPos.Lng - u.distributor.StartPos.Lng
							totalLat := u.distributor.DestPos.Lat - u.distributor.StartPos.Lat
							lngPerFrame := totalLng / totalFrames
							latPerFrame := totalLat / totalFrames
							DebugTraceF("pos change per frame lng %f  lat %f", lngPerFrame, latPerFrame)
							lng, lat := u.distributor.DestPos.minus(u.distributor.CurrentPos) //是否已经足够接近目标点
							DebugTraceF("pos gap lng %f  lat %f", lng, lat)
							if math.Abs(lng) < math.Abs(lngPerFrame) || math.Abs(lat) < math.Abs(latPerFrame) {
								u.distributor.CurrentPos.addLngLat(lng, lat)
								//已经到达目标点，运动停止
								// u.distributor.CurrentSpeed = 0
								// u.distributor.StartPos.setLngLat(u.distributor.DestPos.Lng, u.distributor.DestPos.Lat) //
								u.distributor.StartPos = g_mapData.Points.findLngLat(u.distributor.DestPos.Lng, distributor.DestPos.Lat)
								u.distributor.DestPos = nil
								g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_move_to_new_position, u.distributor) //通知客户端移动到新坐标
								g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_change_state, u.distributor)         //通知客户端移动到新坐标
								DebugInfoF("配送员已经行驶到目标点 %s", u.distributor)
								DebugTraceF("配送员实时位置：%s", u.distributor.PosString())
							} else {
								u.distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
								g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_move_to_new_position, u.distributor) //通知客户端移动到新坐标
								DebugTraceF("配送员实时位置：%s", u.distributor.PosString())
							}
						}
					}
				}
				//----------------------------------------------------------------------------

			}
		}
	}
	go f()
}

func ddd() {
	// if u.distributor.StartPos.equals(u.distributor.DestPos) == false {
	// 	DebugInfoF("配送员 %s 运行路线 %s => %s", u.distributor.Name, u.distributor.StartPos.SimpleString(), u.distributor.DestPos.SimpleString())
	// 	totalTime := u.distributor.Distance * 60 * 60 / (u.distributor.CurrentSpeed * 1000) / realityToSystemTimeRatio //系统中运行路程所花费的时间
	// 	totalFrames := totalTime / timePerFrame                                                                        //一共大约这么多帧就可以走完
	// 	//使用绝对值差距大的作为分片的标准
	// 	totalLng := u.distributor.DestPos.Lng - u.distributor.StartPos.Lng
	// 	totalLat := u.distributor.DestPos.Lat - u.distributor.StartPos.Lat
	// 	lngPerFrame := totalLng / totalFrames
	// 	latPerFrame := totalLat / totalFrames
	// 	DebugTraceF("pos change per frame lng %f  lat %f", lngPerFrame, latPerFrame)
	// 	lng, lat := u.distributor.DestPos.minus(u.distributor.CurrentPos) //是否已经足够接近目标点
	// 	DebugTraceF("pos gap lng %f  lat %f", lng, lat)
	// 	if math.Abs(lng) < math.Abs(lngPerFrame) || math.Abs(lat) < math.Abs(latPerFrame) {
	// 		u.distributor.CurrentPos.addLngLat(lng, lat)
	// 		//已经到达目标点，运动停止
	// 		u.distributor.CurrentSpeed = 0
	// 		u.distributor.StartPos.setLngLat(u.distributor.DestPos.Lng, u.distributor.DestPos.Lat) //
	// 		u.distributor.DestPos = nil
	// 		g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_move_to_new_position, u.distributor) //通知客户端移动到新坐标
	// 		g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_change_state, u.distributor)         //通知客户端移动到新坐标
	// 		DebugInfoF("配送员已经行驶到目标点 %s", u.distributor)
	// 		DebugTraceF("配送员实时位置：%s", u.distributor.PosString())
	// 	} else {
	// 		u.distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
	// 		g_room_distributor.sendMsgToSpecialSubscriber(u.distributor.ID, pro_move_to_new_position, u.distributor) //通知客户端移动到新坐标
	// 		DebugTraceF("配送员实时位置：%s", u.distributor.PosString())
	// 	}
	// }

}
