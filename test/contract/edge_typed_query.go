package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testTypedEdgeQuery(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("Codec", func(t *testing.T) {
		repo := factory(t)
		defer func() { require.NoError(t, repo.Close()) }()

		sourceID, targetID := createEdgeEndpoints(t, repo)
		original := &CustomEdgeModel{
			SourceId: sourceID,
			TargetId: targetID,
			Name:     "codec-query-edge",
			Key:      "codec-key",
		}
		_, err := nod.Edges[CustomEdgeModel](repo).SaveEdge(original)
		require.NoError(t, err)

		models, err := nod.NewTypedEdgeQuery[CustomEdgeModel](repo).
			WithKV().
			WithContent().
			WithTags().
			Where(nod.EdgeFields.Name.Equals(original.Name)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, models, 1)
		require.Equal(t, original.SourceId, models[0].SourceId)
		require.Equal(t, original.TargetId, models[0].TargetId)
		require.Equal(t, original.Name, models[0].Name)
		require.Equal(t, original.Key, models[0].Key)
	})

	t.Run("Adapter", func(t *testing.T) {
		repo := factory(t)
		defer func() { require.NoError(t, repo.Close()) }()
		require.NoError(t, nod.RegisterEdgeAdapter(repo.Adapters(), &CustomEdgeAdapter{}))

		sourceID, targetID := createEdgeEndpoints(t, repo)
		original := &CustomEdgeModelWithAdapter{
			SourceId: sourceID,
			TargetId: targetID,
			Name:     "adapter-query-edge",
			Key:      "adapter-key",
		}
		_, err := nod.Edges[CustomEdgeModelWithAdapter](repo).SaveEdge(original)
		require.NoError(t, err)

		models, err := nod.NewTypedEdgeQuery[CustomEdgeModelWithAdapter](repo).
			WithKV().
			WithContent().
			WithTags().
			Where(nod.EdgeFields.Name.Equals(original.Name)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, models, 1)
		require.Equal(t, original.SourceId, models[0].SourceId)
		require.Equal(t, original.TargetId, models[0].TargetId)
		require.Equal(t, original.Name, models[0].Name)
		require.Equal(t, original.Key, models[0].Key)
	})
}
