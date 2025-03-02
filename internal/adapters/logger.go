package adapters

import (
	"go.uber.org/zap"
)

var logger zap.Logger

func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = *log
}

func GetLogger() *zap.SugaredLogger {
	return logger.Sugar()
}

func SyncLogger() {
	logger.Sync()
}
