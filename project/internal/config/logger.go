package config

import (
	"go.uber.org/zap"
	"log"
)

func NewSugaredLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.StacktraceKey = ""
	logger, err := config.Build()
	if err != nil {
		log.Fatal("Failed to init logger", "err", err)
	}
	return logger.Sugar()
}
