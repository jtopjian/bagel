package file

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
)

var Register = base.Register{
	LuaName: "file",
	Resources: map[string]lua.LGFunction{
		"Delete": NewLuaFileWrapper(Delete),
		"Exists": NewLuaFileWrapper(Exists),
		"Pull":   NewLuaFileWrapper(Pull),
		"Push":   NewLuaFileWrapper(Push),
	},
}
