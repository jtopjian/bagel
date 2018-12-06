package cron

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
)

var Register = base.Register{
	LuaName: "cron",
	Resources: map[string]lua.LGFunction{
		"Entry": base.NewLuaBasicWrapper(Entry),
	},
}
