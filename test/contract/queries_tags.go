package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryTags(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	t.Run("finds every node with the tag", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.Tags().Has("news")).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha", "gamma")
	})

	t.Run("finds a node with a unique tag", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.Tags().Has("ops")).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "delta")
	})

	t.Run("returns an empty result for a missing tag", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.Tags().Has("missing")).
			FindAll()

		require.NoError(t, err)
		require.Empty(t, nodes)
	})
}
