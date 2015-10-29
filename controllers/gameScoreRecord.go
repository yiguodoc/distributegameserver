package controllers

import (
	// "fmt"
	// "github.com/ssor/GobDB"
	"time"
)

type GameScoreRecord struct {
	GameID     string
	UserID     string
	MapID      string
	Score      int
	TimeElapse int
	RecordTime string
	Mode       string //dual or team
}

// func (r *GameScoreRecord) equals(record *GameScoreRecord)bool {
// 	return r.UserID==record.UserID &&

// }

func NewGameScoreRecord(gameID, userID, mapID, mode string, score, timeElapse int) *GameScoreRecord {
	addedTime := time.Now().Format("2006-01-02 15:04:05")
	return &GameScoreRecord{
		GameID:     gameID,
		UserID:     userID,
		MapID:      mapID,
		Score:      score,
		TimeElapse: timeElapse,
		RecordTime: addedTime,
		Mode:       mode,
	}
}
