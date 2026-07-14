package sqlite

import (
	"log/slog"
	"strings"

	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

const sharedMemoryDSN = "file::memory:?mode=memory&cache=shared"

// NewRepository creates a new nod Repository backed by SQLite at the given path.
// Use ":memory:" for an in-memory database.
func NewRepository(path string, log *slog.Logger, mappers *nod.AdapterRegistry) (*nod.Repository, error) {
	db, err := initDB(log, path)
	if err != nil {
		return nil, err
	}
	return nod.NewRepository(db, log, mappers), nil
}

// NewRepositoryInMemory creates a new nod Repository backed by an in-memory SQLite database.
func NewRepositoryInMemory(log *slog.Logger, mappers *nod.AdapterRegistry) (*nod.Repository, error) {
	return NewRepository(sharedMemoryDSN, log, mappers)
}

func initDB(log *slog.Logger, path string) (*gorm.DB, error) {
	log.Debug(">> open database", slog.String("path", path))
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        path,
		DriverName: "sqlite",
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Debug("<< database opened")

	if err := configureConnectionPool(db, path); err != nil {
		return nil, err
	}

	log.Debug(">> enable foreign keys")
	if err := enableForeignKeys(db); err != nil {
		return nil, err
	}
	log.Debug("<< foreign keys enabled")

	log.Debug(">> migrate database")
	if err := nod.Migrate(db); err != nil {
		return nil, err
	}

	log.Debug("<< database migrated")

	return db, nil
}

func configureConnectionPool(db *gorm.DB, path string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if path == ":memory:" || strings.HasPrefix(path, "file::memory:") {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	}

	return nil
}

func enableForeignKeys(db *gorm.DB) error {
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return err
	}

	var enabled int
	if err := db.Raw("PRAGMA foreign_keys;").Scan(&enabled).Error; err != nil {
		return err
	}

	if enabled != 1 {
		return NewForeignKeysDisabledError()
	}

	return nil
}
