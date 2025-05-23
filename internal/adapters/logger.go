package adapters

import (
	"go.uber.org/zap"
)

// Global logger instance shared across the application.
var logger zap.Logger

// Initialize the development logger on startup.
func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = *log
}

// GetLogger returns a sugared version of the global logger for simplified usage.
func GetLogger() *zap.SugaredLogger {
	return logger.Sugar()
}

// SyncLogger synchronizes log buffers to ensure complete flush before shutdown.
func SyncLogger() {
	err := logger.Sync()
	if err != nil {
		return
	}
}
