package nod

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	schemaVersionPropertyKey = "version"

	// CurrentSchemaVersion is the schema version managed by nod migrations.
	CurrentSchemaVersion = 1
)

var (
	// ErrLegacySchemaMigrationRequired is returned when an existing nod schema
	// needs a one-time migration that may clean up orphaned rows before adding
	// database constraints.
	ErrLegacySchemaMigrationRequired = errors.New("nod: legacy schema migration required")
)

// Property stores nod's internal schema properties.
type Property struct {
	Key       string    `gorm:"type:text;primaryKey"`
	Value     string    `gorm:"type:text;not null"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the table used for nod internal properties.
func (Property) TableName() string {
	return "nod_properties"
}

// MigrationOptions controls schema migration behavior.
type MigrationOptions struct {
	LegacySchemaMigration bool
}

// MigrationOption configures schema migration behavior.
type MigrationOption func(*MigrationOptions)

// WithLegacySchemaMigration controls the one-time migration path for existing
// nod schemas that were created before nod tracked schema versions. It is
// enabled by default for now and can be disabled to fail fast on legacy schemas.
func WithLegacySchemaMigration(enabled bool) MigrationOption {
	return func(opts *MigrationOptions) {
		opts.LegacySchemaMigration = enabled
	}
}

// Migrate migrates nod tables and records the applied nod schema version.
func Migrate(db *gorm.DB, options ...MigrationOption) error {
	opts := defaultMigrationOptions()
	for _, option := range options {
		option(&opts)
	}

	hasTables := hasAnyNodTable(db)
	version, hasVersion, err := readSchemaVersion(db)
	if err != nil {
		return err
	}

	needsLegacyMigration := hasTables && (!hasVersion || version < CurrentSchemaVersion)
	if needsLegacyMigration {
		if !opts.LegacySchemaMigration {
			return fmt.Errorf("%w: call nod.WithLegacySchemaMigration(true) to migrate existing nod tables to schema version %d", ErrLegacySchemaMigrationRequired, CurrentSchemaVersion)
		}
		if err := cleanupLegacyRows(db); err != nil {
			return err
		}
	}

	migrateSchema := func() error {
		return db.AutoMigrate(nodSchemaModels()...)
	}
	if needsLegacyMigration && db.Dialector.Name() == "sqlite" {
		if err := withSQLiteForeignKeysDisabled(db, migrateSchema); err != nil {
			return err
		}
	} else {
		if err := migrateSchema(); err != nil {
			return err
		}
	}

	if !hasVersion || version < CurrentSchemaVersion {
		return writeSchemaVersion(db, CurrentSchemaVersion)
	}
	return nil
}

func defaultMigrationOptions() MigrationOptions {
	return MigrationOptions{
		LegacySchemaMigration: true,
	}
}

func nodSchemaModels() []any {
	return []any{
		&NodeCore{},
		&Tag{},
		&NodeTag{},
		&KV{},
		&Content{},
		&Property{},
	}
}

func hasAnyNodTable(db *gorm.DB) bool {
	migrator := db.Migrator()
	return migrator.HasTable(&NodeCore{}) ||
		migrator.HasTable(&Tag{}) ||
		migrator.HasTable(&NodeTag{}) ||
		migrator.HasTable(&KV{}) ||
		migrator.HasTable(&Content{})
}

func readSchemaVersion(db *gorm.DB) (int, bool, error) {
	if !db.Migrator().HasTable(&Property{}) {
		return 0, false, nil
	}

	var property Property
	err := db.First(&property, "key = ?", schemaVersionPropertyKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	version, err := strconv.Atoi(property.Value)
	if err != nil {
		return 0, false, fmt.Errorf("nod: invalid schema version property %q: %w", property.Value, err)
	}
	return version, true, nil
}

func writeSchemaVersion(db *gorm.DB, version int) error {
	return db.Save(&Property{
		Key:   schemaVersionPropertyKey,
		Value: strconv.Itoa(version),
	}).Error
}

func withSQLiteForeignKeysDisabled(db *gorm.DB, fn func() error) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	previousMaxOpenConns := sqlDB.Stats().MaxOpenConnections
	sqlDB.SetMaxOpenConns(1)
	defer sqlDB.SetMaxOpenConns(previousMaxOpenConns)

	if err := setSQLiteForeignKeys(db, false); err != nil {
		return err
	}

	migrationErr := fn()
	enableErr := setSQLiteForeignKeys(db, true)
	var checkErr error
	if enableErr == nil {
		checkErr = checkSQLiteForeignKeys(db)
	}
	return errors.Join(migrationErr, enableErr, checkErr)
}

func setSQLiteForeignKeys(db *gorm.DB, enabled bool) error {
	value := "OFF"
	expected := 0
	if enabled {
		value = "ON"
		expected = 1
	}

	if err := db.Exec("PRAGMA foreign_keys = " + value).Error; err != nil {
		return err
	}

	var actual int
	if err := db.Raw("PRAGMA foreign_keys;").Scan(&actual).Error; err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("nod: sqlite foreign_keys pragma is %d, expected %d", actual, expected)
	}
	return nil
}

func checkSQLiteForeignKeys(db *gorm.DB) error {
	type foreignKeyFailure struct {
		Table  string        `gorm:"column:table"`
		RowID  sql.NullInt64 `gorm:"column:rowid"`
		Parent string        `gorm:"column:parent"`
		FKID   int           `gorm:"column:fkid"`
	}

	var failures []foreignKeyFailure
	if err := db.Raw("PRAGMA foreign_key_check;").Scan(&failures).Error; err != nil {
		return err
	}
	if len(failures) == 0 {
		return nil
	}

	first := failures[0]
	return fmt.Errorf("nod: sqlite foreign key check failed after migration: table=%s rowid=%v parent=%s fkid=%d failures=%d", first.Table, first.RowID, first.Parent, first.FKID, len(failures))
}

func cleanupLegacyRows(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if tx.Migrator().HasTable(&NodeCore{}) {
			if err := tx.Exec(`
	UPDATE node_cores
	SET parent_id = NULL
	WHERE parent_id IS NOT NULL
	  AND (
	    parent_id = ''
	    OR parent_id NOT IN (SELECT id FROM node_cores)
	  )
	`).Error; err != nil {
				return err
			}
		}

		if tx.Migrator().HasTable(&KV{}) && tx.Migrator().HasTable(&NodeCore{}) {
			if err := tx.Exec(`
DELETE FROM kvs
WHERE node_id NOT IN (SELECT id FROM node_cores)
`).Error; err != nil {
				return err
			}
		}

		if tx.Migrator().HasTable(&Content{}) && tx.Migrator().HasTable(&NodeCore{}) {
			if err := tx.Exec(`
DELETE FROM contents
WHERE node_id NOT IN (SELECT id FROM node_cores)
`).Error; err != nil {
				return err
			}
		}

		if tx.Migrator().HasTable(&NodeTag{}) {
			if tx.Migrator().HasTable(&NodeCore{}) {
				if err := tx.Exec(`
DELETE FROM node_tags
WHERE node_id NOT IN (SELECT id FROM node_cores)
`).Error; err != nil {
					return err
				}
			}
			if tx.Migrator().HasTable(&Tag{}) {
				if err := tx.Exec(`
DELETE FROM node_tags
WHERE tag_id NOT IN (SELECT id FROM tags)
`).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
