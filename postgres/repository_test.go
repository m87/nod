package postgres

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/m87/nod"
	"github.com/m87/nod/internal/testsuite/contract"
	"github.com/stretchr/testify/require"
)

func TestRepositoryContractSuite(t *testing.T) {
	dsn := postgresDSN(t)

	contract.RunRepositoryContractTests(t, func(t *testing.T) *nod.Repository {
		t.Helper()

		repo, err := NewRepository(dsn, slog.Default(), nod.NewMapperRegistry())
		require.NoError(t, err)

		cleanupPostgresData(t, repo)
		return repo
	})
}

func postgresDSN(t *testing.T) string {
	t.Helper()

	dsn := os.Getenv("NOD_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set NOD_TEST_POSTGRES_DSN to run postgres integration tests")
	}

	return dsn
}

func cleanupPostgresData(t *testing.T, repo *nod.Repository) {
	t.Helper()

	err := repo.Db.Exec("TRUNCATE TABLE node_tags, kvs, contents, tags, node_cores CASCADE").Error
	require.NoError(t, err)
}

func TestPostgresRepository_MissingDSNSkips(t *testing.T) {
	if os.Getenv("NOD_TEST_POSTGRES_DSN") != "" {
		t.Skip("skip only validated when NOD_TEST_POSTGRES_DSN is unset")
	}

	_, err := NewRepository("", slog.Default(), nod.NewMapperRegistry())
	require.Error(t, err)
	require.NotEmpty(t, fmt.Sprint(err))
}
