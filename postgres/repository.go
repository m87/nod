package postgres

import (
	"log/slog"

	"github.com/m87/nod"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewRepository(dsn string, log *slog.Logger, mappers *nod.MapperRegistry) (*nod.Repository, error) {
	db, err := initDB(log, dsn)
	if err != nil {
		return nil, err
	}

	return &nod.Repository{
		Db:      db,
		Log:     log,
		Mappers: mappers,
	}, nil
}

func initDB(log *slog.Logger, dsn string) (*gorm.DB, error) {
	log.Debug(">> open database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Debug("<< database opened")

	log.Debug(">> migrate database")
	if err := db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{}); err != nil {
		return nil, err
	}
	log.Debug("<< database migrated")

	return db, nil
}
