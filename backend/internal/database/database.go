package database

import (
	"fmt"
	"os"
	"time"

	"geofencing-alerts/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	const maxAttempts = 10
	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				pingErr = sqlDB.Ping()
			}
			if pingErr == nil {
				break
			}
			err = fmt.Errorf("connect database: %w", pingErr)
		}

		if attempt < maxAttempts {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("connect database after %d attempts: %w", maxAttempts, err)
	}

	if err := db.AutoMigrate(
		&models.Geofence{},
		&models.Vehicle{},
		&models.VehicleLocation{},
		&models.AlertRule{},
		&models.Violation{},
		&models.VehicleGeofenceState{},
		&models.AlertEvent{},
	); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}
