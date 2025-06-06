package bootstrap

import (
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"gorm.io/gorm"
)

type Application struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewApp() Application {
	app := &Application{}
	app.Config = config.Load()
	app.DB = database.Connect(&app.Config.Database)

	return *app
}

func (app *Application) Close() {
	database.Close(app.DB)
}
