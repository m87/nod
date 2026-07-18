package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryKV(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	t.Run("finds every edge with the key and value", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.KvString("color").Equals("red")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha", "beta")
	})

	t.Run("does not confuse the same value under another key", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.KvString("color").Equals("blue")).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "gamma")
	})

	t.Run("supports in", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.KvString("color").In([]string{"red", "blue"})).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha", "beta", "gamma")
	})

	t.Run("supports not in", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.KvString("color").NotIn([]string{"red", "blue"})).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "delta")
	})
}
