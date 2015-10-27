package controllers

import (
	// "encoding/json"
	"errors"
	// "fmt"
	// "github.com/BurntSushi/toml"
	// "github.com/ungerik/go-dry"
	// "os"
	// "path"
	// "strings"
)

func (m *MainController) DeleteUser() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("id")
		if g_var.userdb.Has(id) == false {
			return nil, errors.New("用户名错误")
		}
		g_var.userdb.Delete(id)
		return nil, nil
	})
}

func (m *MainController) Resetpwd() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		id := m.GetString("id")
		if g_var.userdb.Has(id) == false {
			return nil, errors.New("用户ID错误")
		}
		pFindUser := func(u *User) bool { return u.ID == id }
		_, err := g_var.userdb.forOne(func(u *User) {
			u.setNewPwd(func(user *User) bool { return user.Password == u.Password }, default_password)
		}, pFindUser)
		if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
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
		u := NewUser(id, default_password, name, g_var.userdb)
		g_var.userdb.Put(u.ID, u)
		return nil, nil
	})
}
func (m *MainController) UserList() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		return g_var.userdb.users(), nil
	})
}
