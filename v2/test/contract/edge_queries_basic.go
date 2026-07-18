package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryBasic(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	t.Run("finds all edges without a filter", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha", "beta", "gamma", "delta")
	})

	t.Run("finds one edge with a simple filter", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.EdgeFields.Name.Equals("beta")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "beta")
	})

	t.Run("returns an empty result when nothing matches", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.EdgeFields.Name.Equals("missing")).
			FindAll()

		require.NoError(t, err)
		require.Empty(t, edges)
	})
}
