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
