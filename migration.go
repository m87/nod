package nod

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	schemaVersionPropertyKey = "version"

	// CurrentSchemaVersion is the schema version managed by nod migrations.
	CurrentSchemaVersion = 3
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

	if err := migrateLegacyNodeTables(db); err != nil {
		return err
	}

	if err := db.AutoMigrate(nodSchemaModels()...); err != nil {
		return err
	}

	if !hasVersion || version < CurrentSchemaVersion {
		return writeSchemaVersion(db, CurrentSchemaVersion)
	}
	return nil
}

func migrateLegacyNodeTables(db *gorm.DB) error {
	for _, source := range []string{"kvs", "kv", "node_kv"} {
		if err := migrateLegacyNodeKVTable(db, source); err != nil {
			return err
		}
	}
	for _, source := range []string{"contents", "content", "node_content"} {
		if err := migrateLegacyNodeContentTable(db, source); err != nil {
			return err
		}
	}
	return nil
}

func migrateLegacyNodeKVTable(db *gorm.DB, source string) error {
	const target = "node_kvs"
	if !db.Migrator().HasTable(source) {
		return nil
	}
	if !db.Migrator().HasTable(target) {
		return db.Migrator().RenameTable(source, target)
	}

	var values []*NodeKV
	if err := db.Table(source).Find(&values).Error; err != nil {
		return err
	}
	if len(values) > 0 {
		if err := db.Table(target).Clauses(clause.OnConflict{DoNothing: true}).Create(&values).Error; err != nil {
			return err
		}
	}
	return db.Migrator().DropTable(source)
}

func migrateLegacyNodeContentTable(db *gorm.DB, source string) error {
	const target = "node_contents"
	if !db.Migrator().HasTable(source) {
		return nil
	}
	if !db.Migrator().HasTable(target) {
		return db.Migrator().RenameTable(source, target)
	}

	var values []*NodeContent
	if err := db.Table(source).Find(&values).Error; err != nil {
		return err
	}
	if len(values) > 0 {
		if err := db.Table(target).Clauses(clause.OnConflict{DoNothing: true}).Create(&values).Error; err != nil {
			return err
		}
	}
	return db.Migrator().DropTable(source)
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
		&EdgeTag{},
		&EdgeContent{},
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
