package controllers

import (
// "github.com/ssor/fauxgaux"
// "github.com/gorilla/websocket"
// "encoding/json"
// "time"
// "strings"
// "fmt"
// "reflect"
// "math/rand"
)

var default_time_of_one_loop = 5 * 60

var (
	g_UnitCenter       *DistributorProcessUnitCenter
	g_distributorStore = DistributorList{ //配送员列表
		NewDistributor("d01", "张军", color_orange),
		NewDistributor("d02", "刘晓莉", color_red),
		NewDistributor("d03", "桑鸿庆", color_purple),
	}
	// g_regions = RegionList{
	// 	NewRegion("1", "255,128,128", 39.928935, 39.944789, 116.614041, 116.618676),
	// 	NewRegion("2", "255,179,128", 39.928935, 39.944789, 116.618676, 116.625898),
	// 	NewRegion("3", "255,255,128", 39.928935, 39.944789, 116.625898, 116.639373),
	// }
	// g_room_viewer      *WsRoom            //= NewRoom(eventReceiver)
)

func init() {

	clientMessageTypeCodeCheck()
	restartGame()
	//--------------------------------------------------------------------------
}
func restartGame() {
	if g_UnitCenter != nil {
		// g_UnitCenter.stop()
		g_UnitCenter.restart()
	} else {
		// //加载地图数据
		// mapData := loadMapData()

		// orders := mapData.Points.filter(func(pos *Position) bool {
		// 	return pos.PointType == POSITION_TYPE_ORDER
		// }).Map(OrderList{}, func(pos *Position, list interface{}) interface{} {
		// 	o := NewOrder(generateOrderID(), pos)
		// 	return append(list.(OrderList), o)
		// }).(OrderList).random(rand.New(rand.NewSource(time.Now().UnixNano())), OrderList{})
		// DebugPrintList_Info(orders)

		// l := []string{"d01", "d02", "d03"}
		// filter := func(d *Distributor) bool { return stringInArray(d.ID, l[:]) }
		g_UnitCenter = NewDistributorProcessUnitCenter([]string{"d01", "d02", "d03"}[:1], "", default_time_of_one_loop)
		// g_UnitCenter = NewDistributorProcessUnitCenter(g_distributorStore.clone(filter), orders, mapData, default_time_of_one_loop)
		if g_UnitCenter != nil {
			g_UnitCenter.start()
		}
	}

}

//字符串数组中是否含有指定字符串
func stringInArray(str string, a []string) bool {
	if len(a) <= 0 {
		return false
	}
	if a[0] == str {
		return true
	} else {
		return stringInArray(str, a[1:])
	}
}
