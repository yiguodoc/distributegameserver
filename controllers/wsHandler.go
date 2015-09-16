package controllers

import (
	// "fmt"
	// "github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
	// "strings"
	// "net/url"
)

// Join method handles WebSocket requests for WebSocketController.
func (this *MainController) ServerWSOrderDistribution() {
	requestURI := this.Ctx.Request.RequestURI
	DebugTraceF(requestURI)
	userID := this.GetString("id")
	if len(userID) <= 0 {
		DebugInfoF("no user ID: [%s]", userID)
		http.Error(this.Ctx.ResponseWriter, "no user ID", 404)
		return
	}
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		DebugMustF("Cannot setup WebSocket connection: %s", err)
		return
	}
	// DebugInfoF("websocket => %s", userID)
	// beego.Debug(requestURI)
	// beego.Trace(ws.LocalAddr())
	// Join chat room.
	// distributor := g_UnitCenter.distributors.findOne(func(d *Distributor) bool { return d.ID == userID })
	// g_UnitCenter.wsRoom.join(Subscriber(distributor), ws)
	if distributor := g_UnitCenter.containsDistributor(userID); distributor != nil {
		g_UnitCenter.distributorOnLine(distributor, ws)
		// g_UnitCenter.distributorOnLine(userID, ws)
		// g_room_distributor.join(userID, subscriber_type_distributor, ws)
		defer g_UnitCenter.distributorOffLine(distributor)
		// defer g_UnitCenter.distributorOffLine(userID)
		// defer g_UnitCenter.wsRoom.leave(userID)
		// Message receive loop.
		for {
			_, p, err := ws.ReadMessage()
			if err != nil { //EOF
				break
			}
			// chanPublish <- newEvent(EVENT_MESSAGE, requestURI, string(p))
			// g_UnitCenter.wsRoom.newMessage(userID, (p))
			g_UnitCenter.distributorMessageIn(distributor, p)
			// g_UnitCenter.distributorMessageIn(userID, p)
		}
	}

	// this.TplNames = ""
	this.ServeJson()
}

/*
//观察者视角
func (this *MainController) ServerWSViewer() {
	requestURI := this.Ctx.Request.RequestURI
	DebugTraceF(requestURI)
	userID := this.GetString("id")
	if len(userID) <= 0 {
		DebugInfoF("no user ID: [%s]", userID)
		http.Error(this.Ctx.ResponseWriter, "no user ID", 404)
		return
	}
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		DebugMustF("Cannot setup WebSocket connection: %s", err)
		return
	}
	g_room_viewer.join(Subscriber(NewViewer(userID)), ws)
	// g_room_viewer.join(userID, ws)
	defer g_room_viewer.leave(userID)
	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil { //EOF
			break
		}
		g_room_viewer.newMessage(userID, (p))
	}
	this.ServeJson()
}
*/
