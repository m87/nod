package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testTypedNodeQuery(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("Codec", func(t *testing.T) {
		repo := factory(t)
		defer func() { require.NoError(t, repo.Close()) }()

		original := &CustomModelWithNodeCodec{
			Name:        "codec-query-model",
			Active:      true,
			Description: "codec description",
			Labels:      []string{"codec", "typed"},
			Key:         "codec-key",
		}
		_, err := nod.Nodes[CustomModelWithNodeCodec](repo).SaveNode(original)
		require.NoError(t, err)

		models, err := nod.NewTypedNodeQuery[CustomModelWithNodeCodec](repo).
			WithKV().
			WithContent().
			WithTags().
			Where(nod.NodeFields.Name.Equals(original.Name)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, models, 1)
		require.Equal(t, original.Name, models[0].Name)
		require.Equal(t, original.Active, models[0].Active)
		require.Equal(t, original.Description, models[0].Description)
		require.ElementsMatch(t, original.Labels, models[0].Labels)
		require.Equal(t, original.Key, models[0].Key)
	})

	t.Run("Adapter", func(t *testing.T) {
		repo := factory(t)
		defer func() { require.NoError(t, repo.Close()) }()
		require.NoError(t, nod.RegisterNodeAdapter(repo.Adapters(), &CustomAdapter{}))

		original := &CustomModelWithAdatper{
			Name:        "adapter-query-model",
			Active:      true,
			Description: "adapter description",
			Labels:      []string{"adapter", "typed"},
			Key:         "adapter-key",
		}
		_, err := nod.Nodes[CustomModelWithAdatper](repo).SaveNode(original)
		require.NoError(t, err)

		models, err := nod.NewTypedNodeQuery[CustomModelWithAdatper](repo).
			WithKV().
			WithContent().
			WithTags().
			Where(nod.NodeFields.Name.Equals(original.Name)).
			FindAll()

		require.NoError(t, err)
		require.Len(t, models, 1)
		require.Equal(t, original.Name, models[0].Name)
		require.Equal(t, original.Active, models[0].Active)
		require.Equal(t, original.Description, models[0].Description)
		require.ElementsMatch(t, original.Labels, models[0].Labels)
		require.Equal(t, original.Key, models[0].Key)
	})
}
