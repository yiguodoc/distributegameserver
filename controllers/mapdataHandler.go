package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ungerik/go-dry"
	"os"
	"path"
	"strings"
)

//地图编辑页面
func (m *MainController) AddressEditIndex() {
	m.TplNames = "addressEdit.tpl"
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

func (m *MainController) MapNameList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		return getMapList(), nil
	})
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
