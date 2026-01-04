package sqlite

import (
	"log/slog"

	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewRepository[T nod.NodeModel](path string, log *slog.Logger) *nod.Repository[T] {
	db := InitDB(log, path)
	return NewRepiositoryWithDB[T](db, log)
}

func NewRepiositoryWithDB[T nod.NodeModel](db *gorm.DB, log *slog.Logger) *nod.Repository[T] {
	return &nod.Repository[T]{
		Db:   db,
		Log:  log,
	}
}

func InitDB(log *slog.Logger, path string) *gorm.DB {
	log.Debug(">> open database", slog.String("path", path))
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Debug("<< database opened")
	log.Debug(">> enable foreign keys")
	db.Exec("PRAGMA foreign_keys = ON;")
	log.Debug("<< foreign keys enabled")

	log.Debug(">> migrate database")
	err = db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
	if err != nil {
		panic("failed to migrate database")
	}

	log.Debug("<< database migrated")

	return db
}
