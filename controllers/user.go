package controllers

import (
// "errors"
// "fmt"
)

type User struct {
	ID     string
	Name   string
	Color  string //地图上marker颜色
	dbLink *UserGobDB
}
type userPredictor func(*User) bool
type UserList []*User

func NewUser(id, name, color string, dbLink *UserGobDB) *User {
	return &User{
		ID:     id,
		Name:   name,
		Color:  color,
		dbLink: dbLink,
	}
}

func (u *User) copy() *User {
	return &User{
		ID:     u.ID,
		Name:   u.Name,
		Color:  u.Color,
		dbLink: u.dbLink,
	}
}
