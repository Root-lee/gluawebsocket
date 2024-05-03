package gluawebsocket

import (
	"strings"
	"testing"

	"gluawebsocket/tools"

	lua "github.com/yuin/gopher-lua"
)

func TestDialFail(t *testing.T) {
	if err := evalLua(t, `
        local ws = require("websocket")
        local conn, _, err = ws.dial("ws://localhost:8080/echo", {})
        assert_not_equal(nil, err)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestDialSuccess(t *testing.T) {
	go tools.StartEchoServer("localhost:8080")
	if err := evalLua(t, `
        local ws = require("websocket")
        local c, _, err = ws.dial("ws://localhost:8080/echo", {})
        assert_equal(err, nil)
        assert_not_equal(c, nil)

        local err = c:write_text("hello")
        assert_equal(err, nil)
        local msg_type, msg, err = c:read()
        assert_equal(err, nil)
        assert_equal(msg_type, ws.get_msg_type("text"))
        assert_equal(msg, "hello")

        local err = c:ping("")
        assert_equal(err, nil)
        local err = c:close(1000, "")
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestGetMsgType(t *testing.T) {
	if err := evalLua(t, `
        local ws = require("websocket")

        assert_equal(1, ws.get_msg_type("text"))
        assert_equal(2, ws.get_msg_type("binary"))
        assert_equal(8, ws.get_msg_type("close"))
        assert_equal(9, ws.get_msg_type("ping"))
        assert_equal(10, ws.get_msg_type("pong"))
        assert_equal(-1, ws.get_msg_type("invalid_msg_type"))

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func evalLua(t *testing.T, script string) error {
	L := lua.NewState()
	defer L.Close()

	Preload(L)

	L.SetGlobal("assert_equal", L.NewFunction(func(L *lua.LState) int {
		expected := L.Get(1)
		actual := L.Get(2)

		if expected.Type() != actual.Type() || expected.String() != actual.String() {
			t.Errorf("Expected %s %q, got %s %q", expected.Type(), expected, actual.Type(), actual)
		}

		return 0
	}))

	L.SetGlobal("assert_not_equal", L.NewFunction(func(L *lua.LState) int {
		expected := L.Get(1)
		actual := L.Get(2)

		if expected.Type() == actual.Type() && expected.String() == actual.String() {
			t.Errorf("not expected %s %q, got %s %q", expected.Type(), expected, actual.Type(), actual)
		}

		return 0
	}))

	L.SetGlobal("assert_contains", L.NewFunction(func(L *lua.LState) int {
		contains := L.Get(1)
		actual := L.Get(2)

		if !strings.Contains(actual.String(), contains.String()) {
			t.Errorf("Expected %s %q contains %s %q", actual.Type(), actual, contains.Type(), contains)
		}

		return 0
	}))

	return L.DoString(script)
}
