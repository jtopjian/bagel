package exec

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
)

var Register = base.Register{
	LuaName: "exec",
	Resources: map[string]lua.LGFunction{
		"Run": NewLuaExecWrapper(Run),
	},
}
