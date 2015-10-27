package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	// "github.com/BurntSushi/toml"
	"github.com/ungerik/go-dry"
	// "os"
	// "path"
	// "strings"
)

func (m *MainController) RankIndex() {
	m.Data["gameID"] = m.GetString("gameID")
	m.TplNames = "rankIndex.tpl"
}
func (m *MainController) GameList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("gameID")
		if len(id) > 0 {
			return g_var.gameUnits.find(func(u *GameUnit) bool { return u.ID == id }), nil
		} else {
			return g_var.gameUnits, nil
		}
	})
	// m.Data["json"] = g_gameUnits
	// m.ServeJson()
}
func (m *MainController) NewGame() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		contentType := m.Ctx.Input.Header("Content-Type")
		fmt.Println("content-type : ", contentType)
		// fmt.Println(m.Ctx.Input.Request.Header)
		body := m.Ctx.Input.CopyBody()
		fmt.Println("body => " + string(body))
		type GamePara struct {
			ID    []string
			MapID string
			Mode  string
		}
		var para GamePara
		if err := json.Unmarshal([]byte(body), &para); err != nil {
			return nil, err
		}
		fmt.Println(para)
		if dry.StringListContains(getMapList(), para.MapID) == false {
			return nil, errors.New("地图不存在")
		}
		mapData := loadMapData(para.MapID)
		if mapData == nil {
			return nil, errors.New("载入地图时出错")
		}
		if mapData.TimeLength <= 0 {
			mapData.TimeLength = default_time_of_one_loop
		}
		game := NewGame(para.ID, para.MapID, mapData.TimeLength, para.Mode)
		return nil, startNewGame(game)

		// return nil, nil
	})
}
func (m *MainController) NewGameIndex() {
	m.TplNames = "newGameIndex.tpl"
}

func (m *MainController) GameListIndex() {
	m.TplNames = "gameListIndex.tpl"
	// m.LayoutSections = make(map[string]string)
	// m.LayoutSections["Sidebar"] = "sidebar.tpl"
}
