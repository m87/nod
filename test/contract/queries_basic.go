package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

	t.Run("finds the first matching node", func(t *testing.T) {
		node, err := nod.NewNodeQuery(repo).
			Where(nod.NodeFields.Name.Equals("beta")).
			FindFirst()

		require.NoError(t, err)
		require.Equal(t, "beta", node.Core.Name)
	})

	t.Run("returns record not found when no first node matches", func(t *testing.T) {
		node, err := nod.NewNodeQuery(repo).
			Where(nod.NodeFields.Name.Equals("missing")).
			FindFirst()

		require.Nil(t, node)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("deletes matching nodes", func(t *testing.T) {
		err := nod.NewNodeQuery(repo).
			Where(nod.NodeFields.Name.Equals("beta")).
			DeleteAll()
		require.NoError(t, err)

		_, err = repo.Nodes().GetNode(queryNodeBetaID)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("rejects an unfiltered delete", func(t *testing.T) {
		err := nod.NewNodeQuery(repo).DeleteAll()

		require.ErrorIs(t, err, gorm.ErrMissingWhereClause)
	})
}
