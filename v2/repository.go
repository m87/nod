package nod

import (
	"log/slog"

	"gorm.io/gorm"
)

type Repository struct {
	db      *gorm.DB
	log     *slog.Logger
	mappers *AdapterRegistry
}

func NewRepository(db *gorm.DB, log *slog.Logger, mappers *AdapterRegistry) *Repository {
	if mappers == nil {
		mappers = NewAdapterRegistry()
	}
	return &Repository{
		db:      db,
		log:     log,
		mappers: mappers,
	}
}

// DB returns the underlying GORM database connection.
func (r *Repository) DB() *gorm.DB { return r.db }

// Log returns the repository's logger.
func (r *Repository) Log() *slog.Logger { return r.log }

// Mappers returns the mapper registry used by this repository.
func (r *Repository) Mappers() *AdapterRegistry { return r.mappers }

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
