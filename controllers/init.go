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
)

type game struct {
	distributorIDList []string
	mapName           string
	game_time_loop    int
}

func NewGame(list []string, mapName string, loop int) *game {
	return &game{
		distributorIDList: list,
		mapName:           mapName,
		game_time_loop:    loop,
	}
}
func init() {

	if err := clientMessageTypeCodeCheck(); err != nil {
		DebugSysF(err.Error())
	}
	restartGame()
	//--------------------------------------------------------------------------
}
func restartGame() {
	// dpc.stop()
	// dpc.start()
	game := NewGame([]string{"d01", "d02", "d03"}[:1], "", default_time_of_one_loop)
	if g_UnitCenter != nil {
		DebugInfoF("游戏重新启动...")
		g_UnitCenter.broadcastMsgToSubscribers(pro_2c_restart_game, nil)
		g_UnitCenter.stop()
		// g_UnitCenter.restart()
	}

	g_UnitCenter = NewDistributorProcessUnitCenter(game.distributorIDList, game.mapName, game.game_time_loop)
	// g_UnitCenter = NewDistributorProcessUnitCenter(g_distributorStore.clone(filter), orders, mapData, default_time_of_one_loop)
	if g_UnitCenter != nil {
		g_UnitCenter.start()
		DebugInfoF("游戏启动完成")
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
