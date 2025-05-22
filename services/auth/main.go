package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/server"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
)

func main() {
	logger.InitLogger(os.Stderr, log.DebugLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("AUTH", "failed to load configuration: %v", err)
	}

	httpServer := server.NewAuthHttpServer(cfg)
	httpServer.Run()
}
