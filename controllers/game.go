package controllers

import (
	"fmt"
)

var game_index = 0

func getGameUniqueID() string {
	game_index++
	return fmt.Sprintf("game_%d", game_index)
}

type GamePreditor func(*Game) bool
type Game struct {
	ID                string
	DistributorIDList []string
	MapName           string
	Game_time_loop    int
	Mode              string //dual or team
}

func (g *Game) String() string {
	return fmt.Sprintf("地图名称: %s  时长: %d   参与者: %s", g.MapName, g.Game_time_loop, g.DistributorIDList)
}
func NewGame(list []string, mapName string, loop int, mode string) *Game {
	return &Game{
		ID:                getGameUniqueID(),
		DistributorIDList: list,
		MapName:           mapName,
		Game_time_loop:    loop,
		Mode:              mode,
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
