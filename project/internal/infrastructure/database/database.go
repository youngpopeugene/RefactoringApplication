package database

import (
	"app/internal/config"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.DatabaseCfg, logger *zap.SugaredLogger) *gorm.DB {
	credentials := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pwd)
	conn, err := gorm.Open(postgres.Open(credentials), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to init database", "err", err)
	}
	return conn
}
