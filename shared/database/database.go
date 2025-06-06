package database

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: %v", err)
	}

	return db
}

func Migrate(db *gorm.DB, models []any) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal("failed to migrate database: %v", err)
	}
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
