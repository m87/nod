package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomModelWithAdatper struct {
	Name   string
	Active bool
}

type CustomAdapter struct{}

func (a *CustomAdapter) FromNode(node *nod.Node) (*CustomModelWithAdatper, error) {
	return &CustomModelWithAdatper{
		Name:   node.Core.Name,
		Active: node.Core.Status == "active",
	}, nil
}

func (a *CustomAdapter) ToNode(model *CustomModelWithAdatper) (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{
			Name: model.Name,
			Status: func() string {
				if model.Active {
					return "active"
				}
				return "inactive"
			}(),
		},
	}, nil
}

func (a *CustomAdapter) IsApplicable(node *nod.Node) bool {
	return true
}

func testAdapterSave(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer func() { require.NoError(t, repo.Close()) }()

	adapter := &CustomAdapter{}
	nod.RegisterNodeAdapter(repo.Adapters(), adapter)

	original := &CustomModelWithAdatper{
		Name:   "Test Model",
		Active: true,
	}

	nodeScope := nod.Nodes[CustomModelWithAdatper](repo)

	nodeId, err := nodeScope.SaveNode(original)
	require.NoError(t, err)
	require.NotNil(t, nodeId)
	require.NotEmpty(t, nodeId)

	retrievedModel, err := nodeScope.GetNode(nodeId)
	require.NoError(t, err)
	require.NotNil(t, retrievedModel)
	require.Equal(t, original.Name, retrievedModel.Name)
	require.Equal(t, original.Active, retrievedModel.Active)
}
