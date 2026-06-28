package postgres

import (
	"log/slog"

	"github.com/m87/nod"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewRepository creates a new nod Repository backed by PostgreSQL at the given DSN.
func NewRepository(dsn string, log *slog.Logger, mappers *nod.MapperRegistry, options ...nod.MigrationOption) (*nod.Repository, error) {
	db, err := initDB(log, dsn, options...)
	if err != nil {
		return nil, err
	}

	return nod.NewRepository(db, log, mappers), nil
}

func initDB(log *slog.Logger, dsn string, options ...nod.MigrationOption) (*gorm.DB, error) {
	log.Debug(">> open database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Debug("<< database opened")

	log.Debug(">> migrate database")
	if err := nod.Migrate(db, options...); err != nil {
		return nil, err
	}
	log.Debug("<< database migrated")

	return db, nil
}
