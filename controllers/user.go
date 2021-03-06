package controllers

import (
	"errors"
	// "fmt"
)

type User struct {
	ID       string
	Name     string
	password string
	Color    string //地图上marker颜色
	Team     string
	Score    int
	dbLink   *UserGobDB
}
type userPredictor func(*User) bool
type UserList []*User

func (ul UserList) every(p userPredictor) bool {
	for _, u := range ul {
		if p(u) == false {
			return false
		}
	}
	return true
}

func NewUser(id, name, color string, dbLink *UserGobDB) *User {
	return &User{
		ID:     id,
		Name:   name,
		Color:  color,
		dbLink: dbLink,
	}
}
func (u *User) inTeam(name string) bool {
	return u.Team == name
}
func (u *User) leaveTeam() error {
	u.Team = ""
	return u.dbLink.updateUser(u)
}
func (u *User) joinTeam(name string) error {
	u.Team = name
	return u.dbLink.updateUser(u)
}
func (u *User) copy() *User {
	return &User{
		ID:     u.ID,
		Name:   u.Name,
		Color:  u.Color,
		dbLink: u.dbLink,
	}
}
func (u *User) setNewPwd(p userPredictor, pwdNew string) error {
	id := u.ID
	if p(u) {
		_, err := u.dbLink.forOne(func(u *User) {
			u.password = pwdNew
		}, func(_u *User) bool {
			return _u.ID == id
		})
		return err
	} else {
		return errors.New("当前密码错误！")
	}
}
