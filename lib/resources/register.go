package resources

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/apt"
	"github.com/jtopjian/bagel/lib/resources/cron"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/resources/file"
	"github.com/jtopjian/bagel/lib/resources/log"
	"github.com/jtopjian/bagel/lib/resources/util"
)

func Register(L *lua.LState) {
	// Register Apt
	mt := L.NewTypeMetatable(apt.Register.LuaName)
	L.SetGlobal(apt.Register.LuaName, mt)
	for k, v := range apt.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

	// Register Cron
	mt = L.NewTypeMetatable(cron.Register.LuaName)
	L.SetGlobal(cron.Register.LuaName, mt)
	for k, v := range cron.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

	// Register Exec
	mt = L.NewTypeMetatable(exec.Register.LuaName)
	L.SetGlobal(exec.Register.LuaName, mt)
	for k, v := range exec.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

	// Register File
	mt = L.NewTypeMetatable(file.Register.LuaName)
	L.SetGlobal(file.Register.LuaName, mt)
	for k, v := range file.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

	// Register Log
	mt = L.NewTypeMetatable(log.Register.LuaName)
	L.SetGlobal(log.Register.LuaName, mt)
	for k, v := range log.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

	// Register Util
	mt = L.NewTypeMetatable(util.Register.LuaName)
	L.SetGlobal(util.Register.LuaName, mt)
	for k, v := range util.Register.Resources {
		L.SetField(mt, k, L.NewFunction(v))
	}

}
