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

func NewRepository(db *gorm.DB, log *slog.Logger) *Repository {
	return &Repository{
		db:       db,
		log:      log,
		adapters: NewAdapterRegistry(),
	}
}

func NewRepositoryWithAdapters(db *gorm.DB, log *slog.Logger, adapters *AdapterRegistry) *Repository {
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

// Transaction executes fn in a database transaction. The repository passed to
// fn uses the transactional database handle and preserves the logger and
// adapter registry of the parent repository.
func (r *Repository) Transaction(fn func(txRepository *Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{
			db:       tx,
			log:      r.log,
			adapters: r.adapters,
		})
	})
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
