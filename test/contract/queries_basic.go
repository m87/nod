package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryBasic(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	t.Run("finds all nodes without a filter", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha", "beta", "gamma", "delta")
	})

	t.Run("finds one node with a simple filter", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.NodeFields.Name.Equals("beta")).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "beta")
	})

	t.Run("returns an empty result when nothing matches", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.NodeFields.Name.Equals("missing")).
			FindAll()

		require.NoError(t, err)
		require.Empty(t, nodes)
	})
}
