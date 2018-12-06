package apt

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
)

var Register = base.Register{
	LuaName: "apt",
	Resources: map[string]lua.LGFunction{
		"Key":     base.NewLuaBasicWrapper(Key),
		"Package": base.NewLuaBasicWrapper(Package),
		"PPA":     base.NewLuaBasicWrapper(PPA),
		"Source":  base.NewLuaBasicWrapper(Source),
	},
}
