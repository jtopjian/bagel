package base

import (
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/connections"
)

type Register struct {
	LuaName   string
	Resources map[string]lua.LGFunction
}

type BasicResource func(map[string]interface{}, connections.Connection) (bool, error)

func NewLuaBasicWrapper(r BasicResource) lua.LGFunction {
	return func(L *lua.LState) int {
		var input map[string]interface{}

		tbl := L.CheckTable(1)
		err := gluamapper.Map(tbl, &input)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		ctx := L.Context()
		conn := ctx.Value("connection").(connections.Connection)

		changed, err := r(input, conn)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(changed))
		L.Push(lua.LNil)

		return 2
	}
}
