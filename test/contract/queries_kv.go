package contract

import (
	"testing"
	"time"

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

	t.Run("compares time values", func(t *testing.T) {
		before := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
		cutoff := time.Date(2026, time.July, 18, 12, 0, 0, 0, time.UTC)
		after := time.Date(2026, time.July, 19, 12, 0, 0, 0, time.UTC)

		for _, node := range []*nod.Node{
			{
				Core: nod.NodeCore{Id: "query-time-before", Name: "before", Kind: "event"},
				KV: map[string]*nod.NodeKV{
					"start": {Key: "start", ValueTime: &before},
				},
			},
			{
				Core: nod.NodeCore{Id: "query-time-after", Name: "after", Kind: "event"},
				KV: map[string]*nod.NodeKV{
					"start": {Key: "start", ValueTime: &after},
				},
			},
		} {
			_, err := repo.Nodes().SaveNode(node)
			require.NoError(t, err)
		}

		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.KvTime("start").LessThanOrEqual(cutoff)).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "before")
	})
}
