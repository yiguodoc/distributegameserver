package controllers

import (
	// "github.com/ssor/fauxgaux"
	// "github.com/gorilla/websocket"
	// "encoding/json"
	// "time"
	// "strings"
	"fmt"
	// "reflect"
	// "math/rand"
)

var default_time_of_one_loop = 5 * 60

var (
	// g_UnitCenter       *GameUnit
	g_gameUnits        GameUnitList = GameUnitList{}
	g_distributorStore              = DistributorList{ //配送员列表
		NewDistributor("d01", "张军", color_orange),
		NewDistributor("d02", "刘晓莉", color_red),
		NewDistributor("d03", "桑鸿庆", color_purple),
	}
)

type GamePreditor func(*Game) bool
type Game struct {
	distributorIDList []string
	mapName           string
	game_time_loop    int
	mode              string //dual or team
}

func (g *Game) String() string {
	return fmt.Sprintf("地图名称: %s  时长: %d   参与者: %s", g.mapName, g.game_time_loop, g.distributorIDList)
}
func NewGame(list []string, mapName string, loop int, mode string) *Game {
	return &Game{
		distributorIDList: list,
		mapName:           mapName,
		game_time_loop:    loop,
		mode:              mode,
	}
}

type GameList []*Game

func (gl GameList) findOne(p GamePreditor) *Game {
	if len(gl) <= 0 {
		return nil
	}
	if p(gl[0]) {
		return gl[0]
	} else {
		return gl[1:].findOne(p)
	}
}
func (gl GameList) find(p GamePreditor) GameList {
	return gl.findRecursive(p, GameList{})
}
func (gl GameList) findRecursive(p GamePreditor, l GameList) GameList {
	if len(gl) <= 0 {
		return l
	}
	if p(gl[0]) {
		l = append(l, gl[0])
	}
	return gl[1:].findRecursive(p, l)
}

func init() {

	if err := clientMessageTypeCodeCheck(); err != nil {
		DebugSysF(err.Error())
	}
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
		g_gameUnits = append(g_gameUnits, gameUnit)

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
