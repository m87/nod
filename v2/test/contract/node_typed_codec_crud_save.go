package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomModelWithNodeCodec struct {
	Name   string
	Active bool
}

func (m *CustomModelWithNodeCodec) ToNode() *nod.Node {
	return &nod.Node{
		Core: nod.NodeCore{
			Name: m.Name,
			Status: func() string {
				if m.Active {
					return "active"
				}
				return "inactive"
			}(),
		},
	}
}

func (m *CustomModelWithNodeCodec) FromNode(node *nod.Node) error {
	m.Name = node.Core.Name
	m.Active = node.Core.Status == "active"
	return nil
}

func (m *CustomModelWithNodeCodec) IsApplicable(node *nod.Node) bool {
	return true
}

func testCodecSave(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer func() { require.NoError(t, repo.Close()) }()

	original := &CustomModelWithNodeCodec{
		Name:   "Test Model",
		Active: true,
	}

	nodeScope := nod.Nodes[CustomModelWithNodeCodec](repo)

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
