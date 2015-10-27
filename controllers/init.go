package controllers

import (
// "github.com/ssor/GobDB"
// "github.com/ssor/fauxgaux"
// "github.com/gorilla/websocket"
// "encoding/json"
// "time"
// "strings"
// "fmt"
// "reflect"
// "math/rand"
)

var (
	default_time_of_one_loop = 5 * 60
	default_password         = "111"

// g_UnitCenter       *GameUnit
// g_gameUnits        GameUnitList = GameUnitList{}
// g_distributorStore = DistributorList{ //配送员列表
// 	NewDistributor("d01", "张军", color_orange),
// 	NewDistributor("d02", "刘晓莉", color_red),
// 	NewDistributor("d03", "桑鸿庆", color_purple),
// }
)
var g_var *global_var = &global_var{}

type global_var struct {
	distributors DistributorList
	gameUnits    GameUnitList
	userdb       *UserGobDB
}

func (g *global_var) init() {
	g.gameUnits = GameUnitList{}
	g.distributors = DistributorList{} //配送员列表
	g.userdb = NewUserGobDB()

	g.userdb.init()
	DebugInfoF("load %d user", g_var.userdb.count())

	g.userdb.every(func(u *User) {
		g.distributors = append(g.distributors, NewDistributor(u))
	})
	// g.distributors = DistributorList{ //配送员列表
	// NewDistributor(NewUser("d01", "张军", color_orange, g.distributors)),
	// NewDistributor(NewUser("d02", "刘晓莉", color_red, g.distributors)),
	// NewDistributor(NewUser("d03", "桑鸿庆", color_purple, g.distributors)),
	// }

}

func init() {

	if err := clientMessageTypeCodeCheck(); err != nil {
		DebugSysF(err.Error())
	}
	g_var.init()
	// maps := getMapList()
	// fmt.Println(maps)
	// game := NewGame([]string{"d01", "d02", "d03"}[:1], "", default_time_of_one_loop)
	// startNewGame(game)
	//--------------------------------------------------------------------------
}
func startNewGame(game *Game) error {
	// if g_UnitCenter != nil {
	// 	DebugInfoF("游戏重新启动...")
	// 	g_UnitCenter.broadcastMsgToSubscribers(pro_2c_restart_game, nil)
	// 	g_UnitCenter.stop()
	// }

	// distributors := g_distributorStore.filter(func(d *Distributor) bool { return stringInArray(d.ID, game.distributorIDList) })
	// DebugSysF("%d", len(distributors))

	gameUnit := NewGameUnit(game.distributorIDList, game.mapName, game.game_time_loop)
	if gameUnit != nil {
		gameUnit.start()
		g_var.gameUnits = append(g_var.gameUnits, gameUnit)

		DebugInfoF("创建游戏完成: ")
		DebugInfoF("%s", game)
	}
	return nil
}

// //字符串数组中是否含有指定字符串
// func stringInArray(str string, a []string) bool {
// 	if len(a) <= 0 {
// 		return false
// 	}
// 	if a[0] == str {
// 		return true
// 	} else {
// 		return stringInArray(str, a[1:])
// 	}
// }
