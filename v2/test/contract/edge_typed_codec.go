package contract

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomEdgeModel struct {
	Name string
	Key  string
}

func (c *CustomEdgeModel) ToEdge() (*nod.Edge, error) {
	return &nod.Edge{
		Core: nod.EdgeCore{
			Name: c.Name,
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
	c.Name = edge.Core.Name
	if kv, ok := edge.KV["key"]; ok && kv.ValueText != nil {
		c.Key = *kv.ValueText
	}
	return nil
}

func (c *CustomEdgeModel) IsApplicable(edge *nod.Edge) bool {
	return true
}

func testCustomEdgeCodec(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	createEdgeEndpoints(t, repo)
	edgeScope := nod.Edges[CustomEdgeModel](repo)
	id, err := edgeScope.SaveEdge(&CustomEdgeModel{Name: "custom", Key: "value"})
	require.NoError(t, err)
	require.NoError(t, uuid.Validate(id))

	edge, err := repo.Edges().GetEdge(id)
	require.NoError(t, err)

	customModel := &CustomEdgeModel{}
	err = customModel.FromEdge(edge)
	require.NoError(t, err)
	require.Equal(t, "custom", customModel.Name)
	require.Equal(t, "value", customModel.Key)

	newEdge, err := customModel.ToEdge()
	require.NoError(t, err)
	require.Equal(t, "custom", newEdge.Core.Name)
	require.Equal(t, "value", *newEdge.KV["key"].ValueText)
}
