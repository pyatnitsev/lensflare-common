package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"go.opentelemetry.io/otel"
)

// Init создает подключение к БД и возвращает *gorm.DB
func Init() (*gorm.DB, error) {
	dsn := os.Getenv("PG_URL")
	if dsn == "" {
		return nil, fmt.Errorf("environment variable PG_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Подключаем OpenTelemetry-плагин, если включен флаг
	if os.Getenv("ENABLE_DB_TRACING") == "true" {
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "postgres" // fallback, если переменная не задана
		}

		if err := db.Use(tracing.NewPlugin(
			tracing.WithTracerProvider(otel.GetTracerProvider()),
			tracing.WithDBSystem(dbName),
		)); err != nil {
			log.Printf("failed to init otel gorm plugin: %v", err)
		} else {
			log.Printf("[DB] OpenTelemetry tracing enabled for GORM (db=%s)", dbName)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Пул соединений
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
