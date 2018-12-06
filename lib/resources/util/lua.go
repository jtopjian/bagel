package util

import (
	"github.com/yuin/gopher-lua"
)

type ErrorResource func(string, string)

func NewLuaErrorWrapper(r ErrorResource) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		r(L.ToStringMeta(L.Get(1)).String(), L.ToStringMeta(L.Get(2)).String())

		return 0
	}
}
