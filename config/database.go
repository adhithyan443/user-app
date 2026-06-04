package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"os"
	"time"
	models "user-app/internal/domain"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {

	dsn := getDSN()

	var err error

	gormLogger := logger.New(
		slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		slog.Error("Failed to get sql.DB", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * 60)

	slog.Info("Successfully connected to PostgreSQL")

	return DB
}



func getDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}

func AutoMigrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		slog.Error("Auto migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migration completed successfully")
}
