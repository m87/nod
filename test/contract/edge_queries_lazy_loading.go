package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryLazyLoading(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	findAlpha := func(t *testing.T, query *nod.EdgeQuery) *nod.Edge {
		t.Helper()

		edges, err := query.
			Where(nod.EdgeFields.Id.Equals(queryEdgeAlphaID)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, edges, 1)
		return edges[0]
	}

	t.Run("does not load relations by default", func(t *testing.T) {
		edge := findAlpha(t, nod.NewEdgeQuery(repo))

		require.Nil(t, edge.KV)
		require.Nil(t, edge.Content)
		require.Nil(t, edge.Tags)
	})

	t.Run("loads only KV when requested", func(t *testing.T) {
		edge := findAlpha(t, nod.NewEdgeQuery(repo).WithKV())

		require.Len(t, edge.KV, 2)
		require.Equal(t, "red", requireString(t, edge.KV["color"].ValueText))
		require.Equal(t, "pl", requireString(t, edge.KV["language"].ValueText))
		require.Nil(t, edge.Content)
		require.Nil(t, edge.Tags)
	})

	t.Run("loads only content when requested", func(t *testing.T) {
		edge := findAlpha(t, nod.NewEdgeQuery(repo).WithContent())

		require.Len(t, edge.Content, 2)
		require.Equal(t, "alpha body", requireString(t, edge.Content["body"].Value))
		require.Equal(t, "alpha summary", requireString(t, edge.Content["summary"].Value))
		require.Nil(t, edge.KV)
		require.Nil(t, edge.Tags)
	})

	t.Run("loads only tags when requested", func(t *testing.T) {
		edge := findAlpha(t, nod.NewEdgeQuery(repo).WithTags())

		require.ElementsMatch(t, []string{"news", "featured", "shared"}, tagNames(edge.Tags))
		require.Nil(t, edge.KV)
		require.Nil(t, edge.Content)
	})

	t.Run("loads all requested relations", func(t *testing.T) {
		edge := findAlpha(t, nod.NewEdgeQuery(repo).
			WithKV().
			WithContent().
			WithTags())

		require.Len(t, edge.KV, 2)
		require.Len(t, edge.Content, 2)
		require.ElementsMatch(t, []string{"news", "featured", "shared"}, tagNames(edge.Tags))
	})
}
