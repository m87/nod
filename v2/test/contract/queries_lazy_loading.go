package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryLazyLoading(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	findAlpha := func(t *testing.T, query *nod.NodeQuery) *nod.Node {
		t.Helper()

		nodes, err := query.
			Where(nod.NodeFields.Id.Equals(queryNodeAlphaID)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, nodes, 1)
		return nodes[0]
	}

	t.Run("does not load relations by default", func(t *testing.T) {
		node := findAlpha(t, nod.NewNodeQuery(repo))

		require.Nil(t, node.KV)
		require.Nil(t, node.Content)
		require.Nil(t, node.Tags)
	})

	t.Run("loads only KV when requested", func(t *testing.T) {
		node := findAlpha(t, nod.NewNodeQuery(repo).WithKV())

		require.Len(t, node.KV, 2)
		require.Equal(t, "red", requireString(t, node.KV["color"].ValueText))
		require.Equal(t, "pl", requireString(t, node.KV["language"].ValueText))
		require.Nil(t, node.Content)
		require.Nil(t, node.Tags)
	})

	t.Run("loads only content when requested", func(t *testing.T) {
		node := findAlpha(t, nod.NewNodeQuery(repo).WithContent())

		require.Len(t, node.Content, 2)
		require.Equal(t, "alpha body", requireString(t, node.Content["body"].Value))
		require.Equal(t, "alpha summary", requireString(t, node.Content["summary"].Value))
		require.Nil(t, node.KV)
		require.Nil(t, node.Tags)
	})

	t.Run("loads only tags when requested", func(t *testing.T) {
		node := findAlpha(t, nod.NewNodeQuery(repo).WithTags())

		require.ElementsMatch(t, []string{"news", "featured", "shared"}, tagNames(node.Tags))
		require.Nil(t, node.KV)
		require.Nil(t, node.Content)
	})

	t.Run("loads all requested relations", func(t *testing.T) {
		node := findAlpha(t, nod.NewNodeQuery(repo).
			WithKV().
			WithContent().
			WithTags())

		require.Len(t, node.KV, 2)
		require.Len(t, node.Content, 2)
		require.ElementsMatch(t, []string{"news", "featured", "shared"}, tagNames(node.Tags))
	})
}

func requireString(t *testing.T, value *string) string {
	t.Helper()
	require.NotNil(t, value)
	return *value
}

func tagNames(tags []*nod.Tag) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}
