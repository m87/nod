package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryKV(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	t.Run("finds every node with the key and value", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.KvString("color").Equals("red")).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha", "beta")
	})

	t.Run("does not confuse the same value under another key", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.KvString("color").Equals("blue")).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "gamma")
	})

	t.Run("supports in", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.KvString("color").In([]string{"red", "blue"})).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha", "beta", "gamma")
	})

	t.Run("supports not in", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.KvString("color").NotIn([]string{"red", "blue"})).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "delta")
	})
}
