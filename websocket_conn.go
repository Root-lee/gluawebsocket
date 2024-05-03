package gluawebsocket

import (
	"time"

	"github.com/gorilla/websocket"
	lua "github.com/yuin/gopher-lua"
)

const (
	WRITE_TIMEOUT           = time.Second
	WEBSOCKET_CONN_TYPENAME = "websocket_conn_typename"
)

func registerWebsocketConnType(L *lua.LState) {
	mt := L.NewTypeMetatable(WEBSOCKET_CONN_TYPENAME)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), websocketConnMethods))
}

var websocketConnMethods = map[string]lua.LGFunction{
	"read":       read,
	"write_text": writeText,
	"ping":       ping,
	"close":      closeConn,
}

func checkWebsocketConn(L *lua.LState) *websocket.Conn {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*websocket.Conn); ok {
		return v
	}
	L.ArgError(1, "websocket connection expected")
	return nil
}

func read(L *lua.LState) int {
	conn := checkWebsocketConn(L)
	msgType, msg, err := conn.ReadMessage()
	L.Push(lua.LNumber(msgType))
	L.Push(lua.LString(msg))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 3
}

func writeText(L *lua.LState) int {
	wsConn := checkWebsocketConn(L)
	msg := L.CheckString(2)
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func ping(L *lua.LState) int {
	wsConn := checkWebsocketConn(L)
	msg := L.CheckString(2)
	err := wsConn.WriteControl(websocket.PingMessage, []byte(msg), time.Now().Add(WRITE_TIMEOUT))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func closeConn(L *lua.LState) int {
	wsConn := checkWebsocketConn(L)
	closeCode := L.CheckInt(2)
	msg := L.CheckString(3)
	err := wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, msg), time.Now().Add(WRITE_TIMEOUT))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}
