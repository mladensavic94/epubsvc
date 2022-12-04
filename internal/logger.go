package epubsvc

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger(level string) {
	switch level {
	case "local":
		Logger, _ = zap.NewDevelopment()
	default:
		Logger, _ = zap.NewProduction()
	}
}
