package nod

import (
	"log/slog"

	"gorm.io/gorm"
)

type Repository struct {
	db       *gorm.DB
	log      *slog.Logger
	adapters *AdapterRegistry
}

func NewRepository(db *gorm.DB, log *slog.Logger, adapters *AdapterRegistry) *Repository {
	if adapters == nil {
		adapters = NewAdapterRegistry()
	}
	return &Repository{
		db:       db,
		log:      log,
		adapters: adapters,
	}
}

// DB returns the underlying GORM database connection.
func (r *Repository) DB() *gorm.DB { return r.db }

// Log returns the repository's logger.
func (r *Repository) Log() *slog.Logger { return r.log }

// Adapters returns the repository's adapter registry.
func (r *Repository) Adapters() *AdapterRegistry { return r.adapters }

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
