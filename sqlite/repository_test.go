package sqlite

import (
	"log/slog"
	"testing"

	"github.com/m87/nod"
	"github.com/m87/nod/internal/testsuite/contract"
	"github.com/stretchr/testify/require"
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
	err = repo.Db.Raw("PRAGMA foreign_keys;").Scan(&enabled).Error
	require.NoError(t, err)
	require.Equal(t, 1, enabled)
}

func TestInit_ConfiguresSingleConnectionForInMemory(t *testing.T) {
	repo, err := NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())
	require.NoError(t, err)
	defer func() { require.NoError(t, repo.Close()) }()

	sqlDB, err := repo.Db.DB()
	require.NoError(t, err)

	stats := sqlDB.Stats()
	require.Equal(t, 1, stats.MaxOpenConnections)
}
