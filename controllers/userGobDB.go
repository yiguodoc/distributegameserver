package controllers

import (
	"github.com/ssor/GobDB"
	// "GobDB"
	"errors"
	// "log"
)

type UserGobDB struct {
	DB *GobDB.DB
}

func NewUserGobDB() *UserGobDB {
	return &UserGobDB{
		DB: GobDB.NewDB("users", func() interface{} { var user User; return &user }),
	}
}
func (db *UserGobDB) init() error {
	_, err := db.DB.Init()
	if err != nil {
		return err
	} else {
		db.every(func(u *User) {
			u.dbLink = db
		})
		return nil
	}
}
func (db *UserGobDB) count() int {
	return len(db.DB.ObjectsMap)
}
func (db *UserGobDB) AddUser(user *User) error {
	return db.Put(user.ID, user)
}
func (db *UserGobDB) updateUser(user *User) error {
	if db.Has(user.ID) {
		if err := db.Delete(user.ID); err == nil {
			db.AddUser(user)
		} else {
			return err
		}
	}
	return nil
}
func (db *UserGobDB) users() UserList {
	list := UserList{}
	for _, v := range db.DB.ObjectsMap {
		list = append(list, v.(*User))
	}
	return list
}
func (db *UserGobDB) Put(id string, user *User) error {
	return db.DB.Put(id, user)
}
func (db *UserGobDB) Delete(id string) error {
	return db.DB.Delete(id)
}
func (db *UserGobDB) find(p userPredictor) UserList {
	list := UserList{}
	for _, v := range db.DB.ObjectsMap {
		if p(v.(*User)) {
			list = append(list, v.(*User))
		}
	}
	return list
}
func (db *UserGobDB) every(f func(*User)) {
	db.forEach(f, func(*User) bool { return true })
}
func (db *UserGobDB) forEach(f func(*User), p userPredictor) {
	for _, v := range db.DB.ObjectsMap {
		if p(v.(*User)) {
			f(v.(*User))
		}
	}
}
func (db *UserGobDB) findOne(p userPredictor) *User {
	for _, v := range db.DB.ObjectsMap {
		if p(v.(*User)) {
			return v.(*User)
		}
	}
	return nil
}
func (db *UserGobDB) forOne(f func(*User), p userPredictor) (*User, error) {
	if u := db.findOne(p); u != nil {
		f(u)
		if err := db.Delete(u.ID); err != nil {
			return nil, err
		}
		if err := db.Put(u.ID, u); err != nil {
			return nil, err
		}
		return u, nil
	} else {
		return nil, errors.New("Not Found")
	}
}
func (db *UserGobDB) Has(id string) bool {
	return db.DB.Has(id)
}
