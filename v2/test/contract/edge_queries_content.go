package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryContent(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	t.Run("matches both the content key and value", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Content("body").Equals("alpha body")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha")
	})

	t.Run("does not confuse content with the same value under another key", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Content("summary").Equals("alpha body")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "beta")
	})

	t.Run("returns an empty result for a missing content key", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Content("missing").Equals("alpha body")).
			FindAll()

		require.NoError(t, err)
		require.Empty(t, edges)
	})
}
