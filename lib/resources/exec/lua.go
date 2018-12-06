package exec

import (
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/connections"
)

type ExecResource func(map[string]interface{}, connections.Connection) (*connections.RunResult, error)

func NewLuaExecWrapper(r ExecResource) func(L *lua.LState) int {
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

		result, err := r(input, conn)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(result.ToLTable(L))
		L.Push(lua.LNil)

		return 2
	}
}
