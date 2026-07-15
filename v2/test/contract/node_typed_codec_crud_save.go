package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomModelWithNodeCodec struct {
	Name        string
	Active      bool
	Description string
	Labels      []string
	Key         string
}

func (m *CustomModelWithNodeCodec) ToNode() (*nod.Node, error) {
	node := &nod.Node{
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

	node.Content = map[string]*nod.NodeContent{
		"description": {
			Key:   "description",
			Value: &m.Description,
		},
	}

	node.Tags = []*nod.Tag{}
	for _, label := range m.Labels {
		node.Tags = append(node.Tags, &nod.Tag{
			Name: label,
		})
	}

	node.KV = map[string]*nod.NodeKV{
		"key": {
			Key:       "key",
			ValueText: &m.Key,
		},
	}

	return node, nil
}

func (m *CustomModelWithNodeCodec) FromNode(node *nod.Node) error {
	m.Name = node.Core.Name
	m.Active = node.Core.Status == "active"
	m.Description = *node.Content["description"].Value
	m.Labels = []string{}
	for _, tag := range node.Tags {
		m.Labels = append(m.Labels, tag.Name)
	}
	m.Key = *node.KV["key"].ValueText
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
		Name:        "Test Model",
		Active:      true,
		Description: "This is a test model with NodeCodec.",
		Labels:      []string{"label1", "label2"},
		Key:         "value1",
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
	require.Equal(t, original.Description, retrievedModel.Description)
	require.ElementsMatch(t, original.Labels, retrievedModel.Labels)
	require.Equal(t, original.Key, retrievedModel.Key)
}
