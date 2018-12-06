package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetLogger() *logrus.Logger {

	log := logrus.New()
	if v := viper.GetBool("debug"); v {
		log.Level = logrus.DebugLevel
	}

	return log
}

func SetLogFields(logger *logrus.Logger, fields map[string]interface{}) *logrus.Entry {
	v := logrus.Fields(fields)
	return logger.WithFields(v)
}
