package nod

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	schemaVersionPropertyKey = "version"

	// CurrentSchemaVersion is the schema version managed by nod migrations.
	CurrentSchemaVersion = 2
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

// Migrate migrates nod tables and records the applied nod schema version.
func Migrate(db *gorm.DB) error {
	version, hasVersion, err := readSchemaVersion(db)
	if err != nil {
		return err
	}

	migrateErr := migrateSchemav1(db, version)
	if migrateErr != nil {
		return migrateErr
	}

	if err := db.AutoMigrate(nodSchemaModels()...); err != nil {
		return err
	}

	if !hasVersion || version < CurrentSchemaVersion {
		return writeSchemaVersion(db, CurrentSchemaVersion)
	}
	return nil
}

func migrateSchemav1(db *gorm.DB, schemaVersion int) error {
	if schemaVersion < 2 {
		if db.Migrator().HasTable("kv") {
			if err := db.Migrator().RenameTable("kv", "node_kv"); err != nil {
				return err
			}
		}
		if db.Migrator().HasTable("content") {
			if err := db.Migrator().RenameTable("content", "node_content"); err != nil {
				return err
			}
		}
	}
	return nil
}

func nodSchemaModels() []any {
	return []any{
		&NodeCore{},
		&Tag{},
		&NodeTag{},
		&NodeKV{},
		&NodeContent{},
		&Property{},
		&EdgeCore{},
		&EdgeKV{},
	}
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
		return 0, false, NewInvalidSchemaVersionError(err)
	}
	return version, true, nil
}

func writeSchemaVersion(db *gorm.DB, version int) error {
	return db.Save(&Property{
		Key:   schemaVersionPropertyKey,
		Value: strconv.Itoa(version),
	}).Error
}
