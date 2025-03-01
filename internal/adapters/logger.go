package adapters

import "go.uber.org/zap"

var sugar zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return &sugar
}
