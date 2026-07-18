package contract

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomEdgeModel struct {
	SourceId string
	TargetId string
	Name     string
	Key      string
}

func (c *CustomEdgeModel) ToEdge() (*nod.Edge, error) {
	return &nod.Edge{
		Core: nod.EdgeCore{
			SourceId: c.SourceId,
			TargetId: c.TargetId,
			Name:     c.Name,
			Kind:     "custom-codec",
		},
		KV: map[string]*nod.EdgeKV{
			"key": {Key: "key", ValueText: &c.Key},
		},
	}, nil
}

func (c *CustomEdgeModel) FromEdge(edge *nod.Edge) error {
	if edge == nil {
		return nod.NewEdgeIsNilError()
	}
	c.SourceId = edge.Core.SourceId
	c.TargetId = edge.Core.TargetId
	c.Name = edge.Core.Name
	if kv, ok := edge.KV["key"]; ok && kv.ValueText != nil {
		c.Key = *kv.ValueText
	}
	return nil
}

func (c *CustomEdgeModel) IsApplicable(edge *nod.Edge) bool {
	return edge != nil && edge.Core.Kind == "custom-codec"
}

func testCustomEdgeCodec(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeScope := nod.Edges[CustomEdgeModel](repo)
	id, err := edgeScope.SaveEdge(&CustomEdgeModel{
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
