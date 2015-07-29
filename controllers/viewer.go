package controllers

import (
	// "errors"
	"fmt"
	"github.com/gorilla/websocket"
	// "math/rand"
	// "time"
)

type Viewer struct {
	ID   string
	Conn *websocket.Conn `json:"-"` // Only for WebSocket users; otherwise nil.

}

func NewViewer(id string) *Viewer {
	return &Viewer{
		ID: id,
	}
}
func (v *Viewer) SetConn(conn *websocket.Conn) {
	v.Conn = conn
}
func (v *Viewer) SendBinaryMessage(msg []byte) error {
	if v.Conn != nil {
		return v.Conn.WriteMessage(websocket.TextMessage, msg)
	}
	return error_no_websocket_connection
}
func (v *Viewer) IdEqals(id string) bool {
	return v.ID == id
}
func (v *Viewer) SetOffline() error {
	if v.Conn != nil {
		if err := v.Conn.Close(); err == nil {
			v.Conn = nil
			DebugInfoF("[%s] OffLine WebSocket closed", v.ID)
		} else {
			DebugMustF("[%s] OffLine,But close websocket err: %s", v.ID, err)
			return err
		}
	}
	return nil
}
func (v *Viewer) String() string {
	return fmt.Sprintf("ID: %s", v.ID)
}
func (v *Viewer) GetID() string {
	return v.ID
}
