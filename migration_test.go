package nod

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type legacyNodeKV struct {
	NodeId    string  `gorm:"type:varchar(36);primaryKey"`
	Key       string  `gorm:"type:text;primaryKey"`
	ValueText *string `gorm:"type:text"`
}

func (legacyNodeKV) TableName() string { return "kvs" }

type legacyNodeContent struct {
	NodeId    string    `gorm:"type:varchar(36);primaryKey"`
	Key       string    `gorm:"type:text;primaryKey"`
	Value     *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (legacyNodeContent) TableName() string { return "contents" }

func TestMigratePreservesLegacyNodeData(t *testing.T) {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        ":memory:",
		DriverName: "sqlite",
	}), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON").Error)

	require.NoError(t, db.AutoMigrate(
		&NodeCore{},
		&legacyNodeKV{},
		&legacyNodeContent{},
		&Property{},
	))
	require.NoError(t, db.Create(&NodeCore{Id: "legacy-node", Name: "legacy", Kind: "test"}).Error)
	require.NoError(t, db.Create(&legacyNodeKV{
		NodeId:    "legacy-node",
		Key:       "legacy-key",
		ValueText: Ptr("legacy-value"),
	}).Error)
	require.NoError(t, db.Create(&legacyNodeContent{
		NodeId: "legacy-node",
		Key:    "body",
		Value:  Ptr("legacy body"),
	}).Error)
	require.NoError(t, db.Save(&Property{Key: schemaVersionPropertyKey, Value: "1"}).Error)

	require.NoError(t, Migrate(db))
	require.False(t, db.Migrator().HasTable("kvs"))
	require.False(t, db.Migrator().HasTable("contents"))
	require.True(t, db.Migrator().HasTable("node_kvs"))
	require.True(t, db.Migrator().HasTable("node_contents"))

	var kv NodeKV
	require.NoError(t, db.First(&kv, "node_id = ? AND key = ?", "legacy-node", "legacy-key").Error)
	require.Equal(t, "legacy-value", *kv.ValueText)

	var content NodeContent
	require.NoError(t, db.First(&content, "node_id = ? AND key = ?", "legacy-node", "body").Error)
	require.Equal(t, "legacy body", *content.Value)

	version, found, err := readSchemaVersion(db)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, CurrentSchemaVersion, version)
}

func TestMigrateRepairsVersionTwoNodeTables(t *testing.T) {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        ":memory:",
		DriverName: "sqlite",
	}), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON").Error)

	require.NoError(t, db.AutoMigrate(
		&NodeCore{},
		&legacyNodeKV{},
		&legacyNodeContent{},
		&NodeKV{},
		&NodeContent{},
		&Property{},
	))
	require.NoError(t, db.Create(&NodeCore{Id: "mixed-node", Name: "mixed", Kind: "test"}).Error)
	require.NoError(t, db.Create(&legacyNodeKV{
		NodeId:    "mixed-node",
		Key:       "legacy-key",
		ValueText: Ptr("legacy-value"),
	}).Error)
	require.NoError(t, db.Create(&NodeKV{
		NodeId:    "mixed-node",
		Key:       "current-key",
		ValueText: Ptr("current-value"),
	}).Error)
	require.NoError(t, db.Save(&Property{Key: schemaVersionPropertyKey, Value: "2"}).Error)

	require.NoError(t, Migrate(db))
	require.False(t, db.Migrator().HasTable("kvs"))

	var kvs []NodeKV
	require.NoError(t, db.Order("key").Find(&kvs, "node_id = ?", "mixed-node").Error)
	require.Len(t, kvs, 2)
	require.Equal(t, "current-key", kvs[0].Key)
	require.Equal(t, "legacy-key", kvs[1].Key)
}
