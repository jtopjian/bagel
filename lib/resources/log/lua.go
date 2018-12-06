package log

import (
	"github.com/yuin/gopher-lua"
)

type LogResource func(interface{})

func NewLuaLogWrapper(r LogResource) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		r(L.ToStringMeta(L.Get(1)).String())

		return 0
	}
}
