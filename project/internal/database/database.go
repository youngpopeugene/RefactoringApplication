package database

import (
	"app/internal/config"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.PostgresCfg) (*gorm.DB, error) {
	credentials := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pwd)
	conn, err := gorm.Open(postgres.Open(credentials), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}
