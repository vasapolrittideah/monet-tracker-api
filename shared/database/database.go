package database

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbConfig *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB, models []any) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db.AutoMigrate(models...)
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from GORM: %v", err)
	}

	if err = sqlDB.Close(); err != nil {
		log.Fatal("failed to close database: %v", err)
	}

	log.Info("🎉 connection to database closed")
}
