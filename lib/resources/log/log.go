package log

import (
	"github.com/yuin/gopher-lua"

	"github.com/jtopjian/bagel/lib/resources/base"
	"github.com/jtopjian/bagel/lib/utils"
)

var Register = base.Register{
	LuaName: "log",
	Resources: map[string]lua.LGFunction{
		"Info":  NewLuaLogWrapper(Info),
		"Warn":  NewLuaLogWrapper(Warn),
		"Error": NewLuaLogWrapper(Error),
		"Fatal": NewLuaLogWrapper(Fatal),
	},
}

func Info(v interface{}) {
	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": "log.Info",
	})
	logger.Info(v)
}

func Warn(v interface{}) {
	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": "log.Warn",
	})
	logger.Warn(v)
}

func Error(v interface{}) {
	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": "log.Error",
	})
	logger.Error(v)
}

func Fatal(v interface{}) {
	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": "log.Fatal",
	})
	logger.Fatal(v)
}
