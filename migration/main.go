package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

func main() {
	// load .env file manually only for this file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	os.Setenv("POSTGRES_HOST", "localhost")

	app := bootstrap.NewApp()
	defer app.Close()

	if err := app.DB.AutoMigrate(
		&domain.User{},
	); err != nil {
		log.Fatal("failed to migrate database: %v", err)
	}

	log.Info("🎉 database migrated successfully")
}
