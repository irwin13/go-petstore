package logger

import "go.uber.org/zap"

var appLogger *zap.Logger

func InitAppLogger() *zap.Logger {

	appLogger, _ = zap.NewDevelopment()
	return appLogger
}

func GetAppLogger() *zap.Logger {
	return appLogger
}
