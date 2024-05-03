package gluawebsocket

import (
	"net/http"

	"github.com/gorilla/websocket"

	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("websocket", Loader)
}

var exports = map[string]lua.LGFunction{
	"dial":         dial,
	"get_msg_type": getMsgType,
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	registerWebsocketConnType(L)

	return 1
}

func getMsgType(L *lua.LState) int {
	msgTypeString := L.CheckString(1)
	switch msgTypeString {
	case "text":
		L.Push(lua.LNumber(websocket.TextMessage))
	case "binary":
		L.Push(lua.LNumber(websocket.BinaryMessage))
	case "ping":
		L.Push(lua.LNumber(websocket.PingMessage))
	case "pong":
		L.Push(lua.LNumber(websocket.PongMessage))
	case "close":
		L.Push(lua.LNumber(websocket.CloseMessage))
	default:
		L.Push(lua.LNumber(-1))
	}
	return 1
}

func dial(L *lua.LState) int {
	addr := L.CheckString(1)
	headerInput := L.CheckTable(2)

	header := make(http.Header)
	headerInput.ForEach(func(key lua.LValue, value lua.LValue) {
		header.Set(key.String(), value.String())
	})

	wsConn, _, err := websocket.DefaultDialer.Dial(addr, header)
	ud := L.NewUserData()
	ud.Value = wsConn
	L.SetMetatable(ud, L.GetTypeMetatable(WEBSOCKET_CONN_TYPENAME))
	L.Push(ud)
	// TODO support http response
	L.Push(lua.LNil)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 3
}
