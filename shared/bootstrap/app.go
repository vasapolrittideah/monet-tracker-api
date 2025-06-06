package bootstrap

import (
	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"gorm.io/gorm"
)

type Application struct {
	Config *config.Config
	DB     *gorm.DB
}

func App() Application {
	app := &Application{}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration: %v", err)
	}
	app.Config = cfg

	app.DB, err = database.Connect(&app.Config.Database)
	if err != nil {
		log.Fatal("failed to connect to database: %v", err)
	}

	return *app
}

func (app *Application) Close() {
	database.Close(app.DB)
}
