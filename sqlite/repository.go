package sqlite

import (
	"log/slog"

	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func NewRepository(path string, log *slog.Logger, mappers *nod.MapperRegistry) (*nod.Repository, error) {
	db, err := initDB(log, path)
	if err != nil {
		return nil, err
	}
	return &nod.Repository{
		Db:      db,
		Log:     log,
		Mappers: mappers,
	}, nil
}

func initDB(log *slog.Logger, path string) (*gorm.DB, error) {
	log.Debug(">> open database", slog.String("path", path))
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        path,
		DriverName: "sqlite",
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Debug("<< database opened")
	log.Debug(">> enable foreign keys")
	db.Exec("PRAGMA foreign_keys = ON;")
	log.Debug("<< foreign keys enabled")

	log.Debug(">> migrate database")
	err = db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
	if err != nil {
		return nil, err
	}

	log.Debug("<< database migrated")

	return db, nil
}
