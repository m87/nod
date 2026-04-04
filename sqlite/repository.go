package sqlite

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

const sharedMemoryDSN = "file::memory:?cache=shared"

// NewRepository creates a new nod Repository backed by SQLite at the given path.
// Use ":memory:" for an in-memory database.
func NewRepository(path string, log *slog.Logger, mappers *nod.MapperRegistry) (*nod.Repository, error) {
	db, err := initDB(log, path)
	if err != nil {
		return nil, err
	}
	return nod.NewRepository(db, log, mappers), nil
}

func initDB(log *slog.Logger, path string) (*gorm.DB, error) {
	path = normalizeSQLitePath(path)

	log.Debug(">> open database", slog.String("path", path))
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        path,
		DriverName: "sqlite",
	}), &gorm.Config{})
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
	err = db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
	if err != nil {
		return nil, err
	}

	log.Debug("<< database migrated")

	return db, nil
}

func normalizeSQLitePath(path string) string {
	if path == ":memory:" {
		return sharedMemoryDSN
	}
	return path
}

func configureConnectionPool(db *gorm.DB, path string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if strings.HasPrefix(path, "file::memory:") {
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
		return errors.New("sqlite foreign_keys pragma is disabled")
	}

	return nil
}
