package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego"
	"github.com/ungerik/go-dry"
	"os"
	"path"
	"strings"
)

/*



 */

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

//by walking through the data file dir
func getMapList() []string {
	files, err := dry.ListDirFiles(default_map_data_dir)
	if err != nil {
		return []string{}
	} else {
		fmt.Println(files)
		return dry.StringMap(func(s string) string {
			return strings.Replace(s, path.Ext(s), "", 1)
		}, files)
		// return files
	}
}

type MapData struct {
	Points PositionList
	Lines  LineList
}

func (m *MainController) MapNameList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		return getMapList(), nil
	})
}
func (m *MainController) GameList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("gameID")
		if len(id) > 0 {
			return g_gameUnits.find(func(u *GameUnit) bool { return u.ID == id }), nil
		} else {
			return g_gameUnits, nil
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
		}
		var para GamePara
		if err := json.Unmarshal([]byte(body), &para); err != nil {
			return nil, err
		}
		fmt.Println(para)
		game := NewGame(para.ID, para.MapID, default_time_of_one_loop)
		return nil, startNewGame(game)

		// return nil, nil
	})
}
func (m *MainController) RestartGame() {
	game := NewGame([]string{"d01", "d02", "d03"}[:1], "", default_time_of_one_loop)
	startNewGame(game)
	// restartGame()
	m.ServeJson()
}

func (m *MainController) GameListIndex() {
	m.TplNames = "gameListIndex.tpl"
}
func (m *MainController) RankIndex() {
	m.Data["gameID"] = m.GetString("gameID")
	m.TplNames = "rankIndex.tpl"
}
func (m *MainController) Index() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_distributorStore.findOne(func(d *Distributor) bool { return d.ID == distributorID })
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
	d := g_distributorStore.findOne(func(d *Distributor) bool { return d.ID == distributorID })
	// d := g_UnitCenter.distributors.find(distributorID)
	if d == nil {
		panic("没有配送员 " + distributorID)
	}
	m.Data["distributor"] = d
	setProData(m)
	m.TplNames = "distribution.tpl"
}

//载入地图数据
func loadMapData(mapName string) *MapData {
	var mapData MapData
	if len(mapName) <= 0 {
		mapName = "data"
	}
	mapFilePath := default_map_data_dir + mapName + ".toml"
	if dry.FileExists(mapFilePath) == false {
		DebugInfoF("地图文件 %s 不存在", mapFilePath)
		return nil
	}
	_, err := toml.DecodeFile(mapFilePath, &mapData)
	if err != nil {
		DebugMustF("载入地图数据时出错：%s", err)
		return nil
	} else {
		bornPoints := mapData.Points.filter(func(pos *Position) bool { return pos.IsBornPoint })
		if len(bornPoints) <= 0 {
			DebugSysF("地图不符合要求，至少设置一个出生点")
		}
		DebugInfoF("地图数据载入统计：%d 个出生点 %d 个路径节点  %d 条路径", len(bornPoints), len(mapData.Points), len(mapData.Lines))
		// DebugPrintList_Info(mapData.Points)
		// DebugPrintList_Info(mapData.Lines)
	}
	return &mapData
}

//上传编辑后的地图数据
func (m *MainController) UploadMapData() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		mapID := m.GetString("id")
		if len(mapID) <= 0 {
			return nil, errors.New("地图名称没有指定")
		}
		values := m.Input()
		value, ok := values["data"]
		if !ok {
			DebugMust("地图数据格式异常")
			return nil, errors.New("地图数据格式异常")
		}
		if len(value) <= 0 {
			DebugMust("没有地图数据上传")
			return nil, errors.New("没有地图数据上传")
		}
		rawData := values["data"][0]
		// fmt.Println(rawData)
		var mapData MapData
		err := json.Unmarshal([]byte(rawData), &mapData)
		if err != nil {
			DebugMustF("解析上传地图数据时出错：%s", err)
			return nil, errors.New("解析上传地图数据时出错")
		}
		// fmt.Println(mapData)
		bornPoints := mapData.Points.filter(func(pos *Position) bool { return pos.IsBornPoint })
		if len(bornPoints) <= 0 {
			return nil, errors.New("地图不符合要求，至少设置一个出生点")
		}
		DebugInfoF("接收到上传的地图数据，统计：%d 个出生点 %d 个路径节点  %d 条路径", len(bornPoints), len(mapData.Points), len(mapData.Lines))
		DebugPrintList_Info(mapData.Points)
		DebugPrintList_Info(mapData.Lines)

		mapFilePath := fmt.Sprintf(default_map_data_dir+"%s.toml", mapID)
		if dry.FileExists(mapFilePath) {
			if e := os.Remove(mapFilePath); e != nil {
				return nil, e
			}
		}
		fileMapData, err := os.Create(mapFilePath)
		if err != nil {
			DebugMustF("创建地图文件出错：%s", err)
			return nil, errors.New("创建地图文件出错")
		}
		defer fileMapData.Close()
		err = toml.NewEncoder(fileMapData).Encode(mapData)
		if err != nil {
			DebugMustF("保存地图数据到文件时出错：%s", err)
			return nil, errors.New("保存地图数据到文件时出错")
		}
		return nil, nil
	})

}

//查询输出地图数据
func (m *MainController) MapData() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		mapName := m.GetString("id")
		return loadMapData(mapName), nil
	})
}

//地图编辑页面
func (m *MainController) AddressEditIndex() {
	m.TplNames = "addressEdit.tpl"
}
func (m *MainController) OrderDistributeIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	distributorID := m.GetString("id")
	d := g_distributorStore.findOne(func(d *Distributor) bool { return d.ID == distributorID })
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

	checkPointMap := getCheckPointMap()
	for key, value := range checkPointMap {
		m.Data[key] = value
	}
}

func (m *MainController) Distributors() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("id")
		if len(id) > 0 {
			return g_distributorStore.filter(func(d *Distributor) bool { return d.ID == id }), nil
		}
		gameID := m.GetString("gameID")
		if len(gameID) > 0 {
			return g_distributorStore.filter(func(d *Distributor) bool { return d.GameID == gameID }), nil
		}

		atgame := m.GetString("atgame")
		if len(atgame) > 0 {
			if atgame == "0" {
				return g_distributorStore.filter(func(d *Distributor) bool { return len(d.GameID) <= 0 }), nil
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

		return g_distributorStore, nil
	})
}
func (m *MainController) NewGameIndex() {
	m.TplNames = "newGameIndex.tpl"
}
func (m *MainController) UserListIndex() {
	m.TplNames = "userListIndex.tpl"
}

func (m *MainController) ViewerIndex() {
	m.Data["HOST"] = fmt.Sprintf("%s:%d", m.Ctx.Input.Host(), m.Ctx.Input.Port())
	m.Data["ID"] = getViewerID()
	m.TplNames = "ViewerIndex.tpl"
}
