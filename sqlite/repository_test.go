package sqlite

import (
	"errors"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/m87/nod"
	"github.com/m87/nod/internal/testsuite/contract"
	"github.com/stretchr/testify/require"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepositoryContractSuite(t *testing.T) {
	contract.RunRepositoryContractTests(t, func(t *testing.T) *nod.Repository {
		t.Helper()
		repo, err := NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())
		require.NoError(t, err)
		return repo
	})
}

func TestInit_EnablesForeignKeys(t *testing.T) {
	repo, err := NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())
	require.NoError(t, err)
	defer func() { require.NoError(t, repo.Close()) }()

	var enabled int
	err = repo.DB().Raw("PRAGMA foreign_keys;").Scan(&enabled).Error
	require.NoError(t, err)
	require.Equal(t, 1, enabled)
}

func TestInit_ConfiguresSingleConnectionForInMemory(t *testing.T) {
	repo, err := NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())
	require.NoError(t, err)
	defer func() { require.NoError(t, repo.Close()) }()

	sqlDB, err := repo.DB().DB()
	require.NoError(t, err)

	stats := sqlDB.Stats()
	require.Equal(t, 1, stats.MaxOpenConnections)
}

func TestMigrate_LegacySchemaMigratesByDefault(t *testing.T) {
	path := createLegacySQLiteDatabase(t, false)
	db := openSQLiteDB(t, path, false)
	defer closeGormDB(t, db)

	err := nod.Migrate(db)
	require.NoError(t, err)

	var property nod.Property
	err = db.First(&property, "key = ?", "version").Error
	require.NoError(t, err)
	require.Equal(t, "1", property.Value)
}

func TestMigrate_LegacySchemaCanBeDisabled(t *testing.T) {
	path := createLegacySQLiteDatabase(t, false)
	db := openSQLiteDB(t, path, false)
	defer closeGormDB(t, db)

	err := nod.Migrate(db, nod.WithLegacySchemaMigration(false))
	require.Error(t, err)
	require.True(t, errors.Is(err, nod.ErrLegacySchemaMigrationRequired))
}

func TestMigrate_LegacySchemaCleansRowsAndAddsConstraints(t *testing.T) {
	path := createLegacySQLiteDatabase(t, true)
	db := openSQLiteDB(t, path, false)
	defer closeGormDB(t, db)

	err := nod.Migrate(db)
	require.NoError(t, err)

	var property nod.Property
	err = db.First(&property, "key = ?", "version").Error
	require.NoError(t, err)
	require.Equal(t, "1", property.Value)

	var orphanKVCount int64
	err = db.Model(&nod.KV{}).Where("node_id = ?", "missing-node").Count(&orphanKVCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), orphanKVCount)

	var child nod.NodeCore
	err = db.First(&child, "id = ?", "child").Error
	require.NoError(t, err)
	require.Nil(t, child.ParentId)

	err = db.Create(&nod.KV{NodeId: "missing-node", Key: "blocked", ValueText: ptr("v")}).Error
	require.Error(t, err)

	err = db.Delete(&nod.NodeCore{}, "id = ?", "node").Error
	require.NoError(t, err)

	var contentCount int64
	err = db.Model(&nod.Content{}).Where("node_id = ?", "node").Count(&contentCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), contentCount)
}

func createLegacySQLiteDatabase(t *testing.T, seed bool) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "legacy.db")
	db := openSQLiteDB(t, path, true)
	createLegacySchema(t, db)
	if seed {
		seedLegacyRows(t, db)
	}
	closeGormDB(t, db)
	return path
}

func openSQLiteDB(t *testing.T, path string, disableFKMigration bool) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(gormsqlite.New(gormsqlite.Config{
		DSN:        path,
		DriverName: "sqlite",
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: disableFKMigration,
	})
	require.NoError(t, err)
	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON;").Error)
	return db
}

func closeGormDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
}

func createLegacySchema(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{}))
}

func seedLegacyRows(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`INSERT INTO node_cores (id, kind, status, name, created_at, updated_at)
		 VALUES ('node', 'kind', 'active', 'node', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		`INSERT INTO node_cores (id, parent_id, kind, status, name, created_at, updated_at)
		 VALUES ('child', 'missing-parent', 'kind', 'active', 'child', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		`INSERT INTO tags (id, name, created_at)
		 VALUES ('tag', 'tag', CURRENT_TIMESTAMP)`,
		`INSERT INTO node_tags (node_id, tag_id)
		 VALUES ('node', 'tag')`,
		`INSERT INTO node_tags (node_id, tag_id)
		 VALUES ('missing-node', 'tag')`,
		`INSERT INTO node_tags (node_id, tag_id)
		 VALUES ('node', 'missing-tag')`,
		`INSERT INTO kvs (node_id, key, value_text)
		 VALUES ('node', 'k', 'v')`,
		`INSERT INTO kvs (node_id, key, value_text)
		 VALUES ('missing-node', 'k', 'v')`,
		`INSERT INTO contents (node_id, key, value, created_at, updated_at)
		 VALUES ('node', 'c', 'v', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		`INSERT INTO contents (node_id, key, value, created_at, updated_at)
		 VALUES ('missing-node', 'c', 'v', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
	}

	for _, statement := range statements {
		require.NoError(t, db.Exec(statement).Error)
	}
}

func ptr[T any](value T) *T {
	return &value
}
