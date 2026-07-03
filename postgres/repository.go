package postgres

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/m87/nod"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var errMissingDSN = errors.New("postgres dsn is required")

// NewRepository creates a new nod Repository backed by PostgreSQL at the given DSN.
func NewRepository(dsn string, log *slog.Logger, mappers *nod.MapperRegistry) (*nod.Repository, error) {
	db, err := initDB(log, dsn)
	if err != nil {
		return nil, err
	}

	return nod.NewRepository(db, log, mappers), nil
}

func initDB(log *slog.Logger, dsn string) (*gorm.DB, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errMissingDSN
	}

	log.Debug(">> open database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Debug("<< database opened")

	log.Debug(">> migrate database")
	if err := nod.Migrate(db); err != nil {
		return nil, err
	}
	log.Debug("<< database migrated")

	return db, nil
}
