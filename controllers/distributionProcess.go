package controllers

import (
// "errors"
// "fmt"
// "math"
// "time"
)

/*
经纬度变化1秒/帧
路径两点的经纬度和长度是确定的，配送员的速度确定，那么运行时间就可以确定，根据系统时间与现实时间的比例，确定系统时间每秒行走的距离和改变的经纬度

*/

var (
	realityToSystemTimeRatio float64 = 5  //现实时间与系统时间的比例
	defaultSpeed             float64 = 20 //默认配送员行驶速度 20km/h
	timePerFrame             float64 = 1
)

func orderDistributionProcess(event *SysEvent) {
	DebugTraceF("执行事件：%d: %s", event.eventCode, event.eventCode.name())
	switch event.eventCode {

	case sys_event_distribution_prepared:
		distributor := event.data.(*Distributor)
		distributor.CheckPoint = checkpoint_flag_order_distribute
		// g_room_distributor.broadcastMsgToSubscribers(pro_distribution_prepared, distributor.ID)
		DebugInfoF("%s 配送准备完成", distributor.Name)
		triggerSysEvent(NewSysEvent(sys_event_start_order_distribution, distributor))

	case sys_event_start_order_distribution:
		distributor := event.data.(*Distributor)
		if distributor.CurrentPos == nil { //没有保存的位置信息，设置仓库为默认的出发点
			warehouses := g_mapData.Points.filter(createPositionFilter(POSITION_TYPE_WAREHOUSE))
			if len(warehouses) > 0 {
				distributor.CurrentPos = warehouses[0].copy()
			} else {
				DebugSysF("无法设置出发点")
			}
		}
		if distributor.Speed <= 0 {
			distributor.Speed = defaultSpeed
		}
		if distributor.StartPos == nil {
			distributor.StartPos = distributor.CurrentPos.copy()
		}
		triggerSysEvent(NewSysEvent(sys_event_order_distribute_additional_msg, distributor.ID))
		//启动经纬度值的计算
		// go startRunning(distributor)

	case sys_event_order_distribute_additional_msg:
		id := event.data.(string)
		distributor := g_distributors.find(id)
		type d struct {
			Distributor *Distributor
			MapData     *MapData
		}
		g_room_distributor.sendMsgToSpecialSubscriber(id, pro_distribution_prepared, &d{distributor, g_mapData})
		// g_room_distributor.sendMsgToSpecialSubscriber(id, pro_distribution_prepared, g_mapData)

	case sys_event_reset_destination_request:
		m := event.data.(map[string]interface{})
		if list, err := mappedValue(m).Getter("PositionID", "DistributorID"); err == nil {
			pos := g_mapData.Points.findByID(int(list[0].(float64)))
			if pos == nil {
				DebugMustF("重置目标点出错，不存在编号为 %d 的节点", int(list[0].(float64)))
				return
			}
			distributor := g_distributors.find(list[1].(string))
			if distributor == nil {
				DebugMustF("重置目标点出错，不存在配送员[%s]", distributor.ID)
				return
			}
			if pos.equals(distributor.CurrentPos) {
				DebugMustF("重置目标点出错，不能和当前所在点相同")
				return
			}
			//两点需要在同一条线上
			line := g_mapData.Lines.find(pos, distributor.CurrentPos)
			if line == nil {
				DebugMustF("重置目标点出错，两点不属于同一条路径")
				return
			}
			distributor.DestPos = pos.copy()
			distributor.Distance = line.Distance
			g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_reset_destination, distributor)
		} else {
			DebugMustF("客户端数据格式错误: %s", err)
		}
	case sys_event_change_state_request:
		if list, err := mappedValue(event.data.(map[string]interface{})).Getter("DistributorID", "State"); err == nil {
			distributor := g_distributors.find(list[0].(string))
			if distributor == nil {
				DebugMustF("重置目标点出错，不存在配送员[%s]", distributor.ID)
				return
			}
			state := int(list[1].(float64))
			if state == 0 {
				distributor.CurrentSpeed = 0
			} else {
				distributor.CurrentSpeed = distributor.Speed
			}
			g_room_distributor.sendMsgToSpecialSubscriber(distributor.ID, pro_change_state, distributor)
		} else {
			DebugMustF("客户端数据格式错误: %s", err)
		}
	case sys_event_order_distribute_result:

	}
}

/*
func processorCreator(distributor *Distributor) func(sysEventCode, *Distributor) {
	f := func(code sysEventCode, d *Distributor) {
		var timePerFrame float64 = 1
		timer := time.Tick(time.Duration(timePerFrame) * time.Second)
		for {
			select {
			case <-timer: //计算新位置
				if distributor.StartPos != nil && distributor.DestPos != nil {
					if distributor.StartPos.equals(distributor.DestPos) {
						continue
					}
					DebugInfoF("配送员 %s 运行路线 %s => %s", distributor.ID, distributor.StartPos.SimpleString(), distributor.DestPos.SimpleString())
					totalTime := distributor.Distance * 60 * 60 / (defaultSpeed * 1000) / realityToSystemTimeRatio //系统中运行路程所花费的时间
					totalFrames := totalTime / timePerFrame                                                        //一共大约这么多帧就可以走完
					//使用绝对值差距大的作为分片的标准
					totalLng := distributor.DestPos.Lng - distributor.StartPos.Lng
					totalLat := distributor.DestPos.Lat - distributor.StartPos.Lat
					lngPerFrame := totalLng / totalFrames
					latPerFrame := totalLat / totalFrames
					DebugTraceF("pos change per frame lng %f  lat %f", lngPerFrame, latPerFrame)
					lng, lat := distributor.DestPos.minus(distributor.CurrentPos) //是否已经足够接近目标点
					DebugTraceF("pos gap lng %f  lat %f", lng, lat)
					if math.Abs(lng) < math.Abs(lngPerFrame) || math.Abs(lat) < math.Abs(latPerFrame) {
						distributor.CurrentPos.addLngLat(lng, lat)
						DebugInfoF("配送员已经行驶到目标点 %s", distributor)
					} else {
						distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
					}
					DebugTraceF("配送员实时位置：%s", distributor.PosString())
				}
			}
		}
	}
	return f
}
func startRunning(distributor *Distributor) {
	var timePerFrame float64 = 1
	timer := time.Tick(time.Duration(timePerFrame) * time.Second)
	for {
		select {
		case <-timer: //计算新位置
			if distributor.StartPos != nil && distributor.DestPos != nil {
				if distributor.StartPos.equals(distributor.DestPos) {
					continue
				}
				DebugInfoF("配送员 %s 运行路线 %s => %s", distributor.ID, distributor.StartPos.SimpleString(), distributor.DestPos.SimpleString())
				totalTime := distributor.Distance * 60 * 60 / (defaultSpeed * 1000) / realityToSystemTimeRatio //系统中运行路程所花费的时间
				totalFrames := totalTime / timePerFrame                                                        //一共大约这么多帧就可以走完
				//使用绝对值差距大的作为分片的标准
				totalLng := distributor.DestPos.Lng - distributor.StartPos.Lng
				totalLat := distributor.DestPos.Lat - distributor.StartPos.Lat
				lngPerFrame := totalLng / totalFrames
				latPerFrame := totalLat / totalFrames
				DebugTraceF("pos change per frame lng %f  lat %f", lngPerFrame, latPerFrame)
				lng, lat := distributor.DestPos.minus(distributor.CurrentPos) //是否已经足够接近目标点
				DebugTraceF("pos gap lng %f  lat %f", lng, lat)
				if math.Abs(lng) < math.Abs(lngPerFrame) || math.Abs(lat) < math.Abs(latPerFrame) {
					distributor.CurrentPos.addLngLat(lng, lat)
					DebugInfoF("配送员已经行驶到目标点 %s", distributor)
				} else {
					distributor.CurrentPos.addLngLat(lngPerFrame, latPerFrame)
				}
				DebugTraceF("配送员实时位置：%s", distributor.PosString())
				// if math.Abs(totalLng) >= math.Abs(totalLat) {
				// } else {

				// }
			}
		}
	}
}
*/
