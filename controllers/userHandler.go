package controllers

import (
	"encoding/json"
	"errors"
	// "fmt"
	// "github.com/BurntSushi/toml"
	"github.com/ungerik/go-dry"
	// "os"
	// "path"
	// "strings"
	// "log"
)

type groupInfo struct {
	Name string
	List []string
}

func (m *MainController) LeaveTeam() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		list := unmashUserIdList(m.Ctx.Input.CopyBody())
		if len(list) <= 0 {
			return nil, errors.New("没有传入成员需要离开团队")
		}
		//确保要组队的人当前不在任何团队中
		predictorUsersToGroup := func(u *User) bool {
			return dry.StringListContains(list, u.ID) == true
		}
		g_var.userdb.forEach(func(u *User) {
			u.leaveTeam()
		}, predictorUsersToGroup)
		return nil, nil
	})

}

//组建或者加入团队
func (m *MainController) AddTeam() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		var group groupInfo
		err := json.Unmarshal(m.Ctx.Input.CopyBody(), &group)
		if err != nil {
			DebugSysF("AddGroup: 解析数据出错：%s", err)
		}
		if len(group.List) <= 0 {
			return nil, errors.New("没有成员可以被加入到团队中")
		}
		//确保要组队的人当前不在任何团队中
		predictorUsersToGroup := func(u *User) bool {
			return dry.StringListContains(group.List, u.ID) == true
		}
		if g_var.userdb.find(predictorUsersToGroup).every(func(u *User) bool {
			return u.inTeam("") || u.inTeam(group.Name)
		}) == false {
			return nil, errors.New("有成员已经加入别的团队了")
		}
		g_var.userdb.forEach(func(u *User) {
			u.joinTeam(group.Name)
		}, predictorUsersToGroup)
		return nil, nil
	})
}

func unmashUserIdList(data []byte) []string {
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		DebugSysF("解析用户名列表出错：%s", err)
		DebugSysF("<- (%s)", string(data))
	}
	return list
}
func (m *MainController) DeleteUser() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		dry.StringEachMust(func(id string) {
			if g_var.userdb.Has(id) == false {
				DebugSysF("用户 %s 不存在", id)
				return
			}
			g_var.userdb.Delete(id)
		}, unmashUserIdList(m.Ctx.Input.CopyBody()))

		return nil, nil
	})
}

func (m *MainController) Resetpwd() {
	responseHandler(m, func(m *MainController) (interface{}, error) {

		dry.StringEachMust(func(id string) {
			DebugTraceF("%s 重置密码", id)
			pFindUser := func(u *User) bool { return u.ID == id }
			g_var.userdb.forOne(func(u *User) {
				u.setNewPwd(func(user *User) bool { return user.password == u.password }, default_password)
			}, pFindUser)
		}, unmashUserIdList(m.Ctx.Input.CopyBody()))
		// log.Println(m.Input())

		return nil, nil
		// id := m.GetString("id")
		// if g_var.userdb.Has(id) == false {
		// 	return nil, errors.New("用户ID错误")
		// }
		// pFindUser := func(u *User) bool { return u.ID == id }
		// _, err := g_var.userdb.forOne(func(u *User) {
		// 	u.setNewPwd(func(user *User) bool { return user.password == u.password }, default_password)
		// }, pFindUser)
		// if err != nil {
		// 	return nil, err
		// } else {
		// 	return nil, nil
		// }
	})
}
func (m *MainController) AddUser() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("id")
		name := m.GetString("name")
		if len(id) <= 0 || len(name) <= 0 {
			return nil, errors.New("注册的用户名错误")
		}

		if g_var.userdb.Has(id) == true {
			return nil, errors.New("用户名已被注册")
		}
		u := NewUser(id, name, "", g_var.userdb)
		g_var.userdb.Put(u.ID, u)
		return nil, nil
	})
}
func (m *MainController) UserList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		return g_var.userdb.users(), nil
	})
}

func (m *MainController) UserDetailIndex() {
	m.TplNames = "userDetailIndex.tpl"
}
