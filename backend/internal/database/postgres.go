package database

import (
	"log"
	"time"

	"backend/internal/config"
	"backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgres opens a GORM connection and applies migrations.
func InitPostgres(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	log.Println("Postgres connected")
	return db, nil
}
