# A websocket lib for GopherLua

A websocket client for the [GopherLua](https://github.com/yuin/gopher-lua) VM, based on [gorilla/websocket](https://github.com/gorilla/websocket)

## Installation
```bash
go get github.com/Root-lee/gluawebsocket
```

## Using

### Loading Modules

```go
import (
	"github.com/Root-lee/gluawebsocket"
)

// Bring up a GopherLua VM
L := lua.NewState()
defer L.Close()

// Preload websocket modules
gluawebsocket.Preload(L)
```

### Usage In lua <a name="lua-demo-anchor"></a>

```lua
local ws = require("websocket")

local c, _, err = ws.dial("ws://localhost:8080/echo", {})
if err then
    print(err)
    return
end

err = c:write_text("hello")
if err then
    print(err)
    return
end

local msg_type, msg
msg_type, msg, err = c:read()
if err then
    print(err)
    return
end

print(string.format("receive a message, msg_type: %d, msg: %s", msg_type, msg))

err = c:ping("ping")
if err then
    print(err)
    return
end

err = c:close(1000, "")
if err then
    print(err)
    return
end

```

## Testing

### Unit Test
```bash
$ go test github.com/Root-lee/gluawebsocket...
PASS
coverage: 89.4% of statements
ok  	gluawebsocket	0.389s
```

### Manual Test
You can use this [Websocket Echo Server](https://github.com/gorilla/websocket/tree/release-1.5/examples/echo) to start a websocket server

Then you can use this demo to test your lua file

```go
package main

import (
	"github.com/Root-lee/gluawebsocket"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	gluawebsocket.Preload(L)
	defer L.Close()

	if err := L.DoFile("test.lua"); err != nil {
		panic(err)
	}
}
```
You can refer to this [lua script](#lua-demo-anchor) to write your own lua script

## License

MIT
