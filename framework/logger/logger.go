package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Log struct {
	Logger *logrus.Logger
}

func InitLogger() *Log {
	var logger = logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}

	if os.Getenv("environment") == "development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return &Log{Logger: logger}
}
