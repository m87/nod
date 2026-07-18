package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryTags(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	t.Run("finds every edge with the tag", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Tags().Has("news")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha", "gamma")
	})

	t.Run("finds an edge with a unique tag", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Tags().Has("ops")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "delta")
	})

	t.Run("returns an empty result for a missing tag", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Tags().Has("missing")).
			FindAll()

		require.NoError(t, err)
		require.Empty(t, edges)
	})
}
