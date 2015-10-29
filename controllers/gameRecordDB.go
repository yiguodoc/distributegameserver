package controllers

import (
	// "fmt"
	"github.com/ssor/GobDB"
	// "time"
)

type GameRecordDB struct {
	DB *GobDB.DB
}

func NewGameRecordDB() *GameRecordDB {
	return &GameRecordDB{
		DB: GobDB.NewDB("records", func() interface{} { var record GameScoreRecord; return &record }),
	}
}
func (db *GameRecordDB) init() error {
	_, err := db.DB.Init()
	if err != nil {
		return err
	} else {
		// db.every(func(u *User) {
		// 	u.dbLink = db
		// })
		return nil
	}
}
func (db *GameRecordDB) add(record *GameScoreRecord) {
	db.DB.Put(record.UserID+"_"+record.RecordTime, record)
}
