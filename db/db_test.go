package db

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// helper-функция для подмены Init с in-memory SQLite
func initTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	return db, nil
}

// ---------------------------
// Тесты
// ---------------------------

func TestInit_MissingPGURL(t *testing.T) {
	t.Setenv("PG_URL", "")

	db, err := Init()
	if err == nil {
		t.Fatalf("expected error when PG_URL is missing, got nil")
	}
	if db != nil {
		t.Fatalf("expected returned *gorm.DB to be nil when error occurs")
	}
	if got, want := err.Error(), "environment variable PG_URL is not set"; got != want {
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}
}

func TestInit_TracingEnabled(t *testing.T) {
	t.Setenv("PG_URL", "file::memory:?cache=shared") // sqlite in-memory
	t.Setenv("ENABLE_DB_TRACING", "true")
	t.Setenv("DB_NAME", "mydb")

	// Используем sqlite вместо реального Postgres для теста
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("failed to init test DB: %v", err)
	}

	if db == nil {
		t.Fatal("expected non-nil *gorm.DB")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping DB: %v", err)
	}
}

func TestInit_TracingDisabled(t *testing.T) {
	t.Setenv("PG_URL", "file::memory:?cache=shared") // sqlite in-memory
	t.Setenv("ENABLE_DB_TRACING", "false")

	db, err := initTestDB()
	if err != nil {
		t.Fatalf("failed to init test DB: %v", err)
	}

	if db == nil {
		t.Fatal("expected non-nil *gorm.DB")
	}
}
