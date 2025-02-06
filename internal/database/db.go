package database

import (
	"fmt"
	"platform-service/internal/config"
	"platform-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	driver := config.GetDBDriver()
	connectionString := config.GetDBConnectionString()

	switch driver {
	case "postgres":
		DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return nil
}
