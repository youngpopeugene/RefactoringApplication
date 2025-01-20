package main

import (
	"app/internal/api"
	"app/internal/config"
	"app/internal/infrastructure/database"
)

func main() {
	cfg := config.NewConfig()
	logger := config.NewSugaredLogger()
	db := database.New(cfg.Database, logger)
	api.Run(db, logger)
}
