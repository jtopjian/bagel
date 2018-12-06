package util

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
)

var Register = base.Register{
	LuaName: "util",
	Resources: map[string]lua.LGFunction{
		"LogIfError":  NewLuaErrorWrapper(LogIfError),
		"StopIfError": NewLuaErrorWrapper(StopIfError),
	},
}
