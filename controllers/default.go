package controllers

import (
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/BurntSushi/toml"
	"github.com/astaxie/beego"
	// "github.com/ungerik/go-dry"
	// "os"
	// "path"
	// "strings"
)

var viewerCount = 1
var default_map_data_dir = "./mapdata/"

func getViewerID() string {
	viewerCount++
	return fmt.Sprintf("953957392%d", viewerCount)
}

type ResponseMsg struct {
	Code    int
	Message string
	Data    interface{}
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

type logicHandler func(m *MainController) (interface{}, error)

func responseHandler(m *MainController, handler logicHandler) {
	response := NewResponseMsg(0)
	defer func() {
		m.Data["json"] = response
		m.ServeJson()
	}()

	if value, err := handler(m); err != nil {
		DebugMustF("controller error: %s", err.Error())
		response = NewResponseMsg(1, err.Error())
	} else {
		response.Data = value
	}
}

type MainController struct {
	beego.Controller
}

//客户端游戏主页
func (m *MainController) Index() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_var.distributors.findOne(func(d *Distributor) bool { return d.UserInfo.ID == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)

	m.TplNames = "index.tpl"
}

//客户端游戏登录页面
func (m *MainController) Login() {
	m.TplNames = "login.tpl"
}

//向模板注入定义的数据
func setProData(m *MainController) {
	codes := getClientMessageTypeCodeList()
	for _, code := range codes {
		m.Data[code.name()] = code
	}

	checkPointMap := getCheckPointMap()
	for key, value := range checkPointMap {
		m.Data[key] = value
	}
}

func (m *MainController) Distributors() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("id")
		if len(id) > 0 {
			return g_var.distributors.filter(func(d *Distributor) bool { return d.UserInfo.ID == id }), nil
		}
		gameID := m.GetString("gameID")
		if len(gameID) > 0 {
			return g_var.distributors.filter(func(d *Distributor) bool { return d.GameID == gameID }), nil
		}

		atgame := m.GetString("atgame")
		if len(atgame) > 0 {
			if atgame == "0" {
				return g_var.distributors.filter(func(d *Distributor) bool { return len(d.GameID) <= 0 }), nil
			}
		}
		// 	gameID := m.GetString("gameID")
		// 	if len(gameID) > 0 {
		// 		unit := g_gameUnits.findOne(func(gu *GameUnit) bool { return gu.ID == gameID })
		// 		if unit == nil {
		// 			return nil, errors.New("no such game")
		// 		}
		// 		return unit.Distributors
		// 	}

		// }

		return g_var.distributors, nil
	})
}

func (m *MainController) UserAdminIndex() {
	m.TplNames = "userAdminIndex.tpl"
}

func (m *MainController) ViewerIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	m.Data["ID"] = getViewerID()
	m.TplNames = "ViewerIndex.tpl"
}

//废弃
func (m *MainController) DistributionIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_var.distributors.findOne(func(d *Distributor) bool { return d.UserInfo.ID == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)
	m.TplNames = "distribution.tpl"
}

//废弃
func (m *MainController) OrderDistributeIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_var.distributors.findOne(func(d *Distributor) bool { return d.UserInfo.ID == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)
	m.TplNames = "orderDistribute.tpl"
}
