package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testRepositoryClose(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	require.NoError(t, repo.Close())
}
