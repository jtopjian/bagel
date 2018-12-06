package util

import (
	"fmt"

	"github.com/jtopjian/bagel/lib/utils"
)

func LogIfError(msg, err string) {
	if err != "" && err != "nil" {
		logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
			"resource": "log.Error",
		})
		logger.Error(fmt.Sprintf("%s: %s", msg, err))
	}
}

func StopIfError(msg, err string) {
	if err != "" && err != "nil" {
		logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
			"resource": "log.Error",
		})
		logger.Fatal(fmt.Sprintf("%s: %s", msg, err))
	}
}
