package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/server"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
)

func main() {
	logger.InitLogger(os.Stderr, log.DebugLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("AUTH", "failed to load configuration: %v", err)
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("AUTH", "failed to connect to database: %v", err)
	} else {
		logger.Info("AUTH", "🎉 connected to database: %s", cfg.Database.Name)
	}

	entities := []any{
		&domain.ExternalLogin{},
	}

	if err := database.Migrate(db, entities); err != nil {
		logger.Fatal("USER", "failed to migrate database: %v", err)
	}

	httpServer := server.NewAuthHttpServer(cfg, db)
	httpServer.Run()
}
