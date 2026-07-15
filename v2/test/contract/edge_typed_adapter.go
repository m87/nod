package contract

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomEdgeModelWithAdapter struct {
	SourceId string
	TargetId string
	Name     string
	Key      string
}

type CustomEdgeAdapter struct{}

func (a *CustomEdgeAdapter) ToEdge(model *CustomEdgeModelWithAdapter) (*nod.Edge, error) {
	return &nod.Edge{
		Core: nod.EdgeCore{
			SourceId: model.SourceId,
			TargetId: model.TargetId,
			Name:     model.Name,
			Kind:     "custom-adapter",
		},
		KV: map[string]*nod.EdgeKV{
			"key": {Key: "key", ValueText: &model.Key},
		},
	}, nil
}

func (a *CustomEdgeAdapter) FromEdge(edge *nod.Edge) (*CustomEdgeModelWithAdapter, error) {
	model := &CustomEdgeModelWithAdapter{
		SourceId: edge.Core.SourceId,
		TargetId: edge.Core.TargetId,
		Name:     edge.Core.Name,
	}
	if kv, ok := edge.KV["key"]; ok && kv.ValueText != nil {
		model.Key = *kv.ValueText
	}
	return model, nil
}

func (a *CustomEdgeAdapter) IsApplicable(edge *nod.Edge) bool {
	return edge != nil && edge.Core.Kind == "custom-adapter"
}

func testCustomEdgeAdapter(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	adapter := &CustomEdgeAdapter{}
	require.NoError(t, nod.RegisterEdgeAdapter(repo.Adapters(), adapter))

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeScope := nod.Edges[CustomEdgeModelWithAdapter](repo)
	id, err := edgeScope.SaveEdge(&CustomEdgeModelWithAdapter{
		SourceId: sourceID,
		TargetId: targetID,
		Name:     "custom",
		Key:      "value",
	})
	require.NoError(t, err)
	require.NoError(t, uuid.Validate(id))

	customModel, err := edgeScope.GetEdge(id)
	require.NoError(t, err)
	require.Equal(t, sourceID, customModel.SourceId)
	require.Equal(t, targetID, customModel.TargetId)
	require.Equal(t, "custom", customModel.Name)
	require.Equal(t, "value", customModel.Key)
}
