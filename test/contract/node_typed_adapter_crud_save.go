package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type CustomModelWithAdatper struct {
	Name        string
	Active      bool
	Description string
	Labels      []string
	Key         string
}

type CustomAdapter struct{}

func (a *CustomAdapter) FromNode(node *nod.Node) (*CustomModelWithAdatper, error) {
	model := &CustomModelWithAdatper{
		Name:   node.Core.Name,
		Active: node.Core.Status == "active",
	}

	model.Description = *node.Content["description"].Value
	model.Labels = []string{}
	for _, tag := range node.Tags {
		model.Labels = append(model.Labels, tag.Name)
	}

	model.Key = *node.KV["key"].ValueText

	return model, nil
}

func (a *CustomAdapter) ToNode(model *CustomModelWithAdatper) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Name: model.Name,
			Status: func() string {
				if model.Active {
					return "active"
				}
				return "inactive"
			}(),
		},
	}

	node.Content = map[string]*nod.NodeContent{
		"description": {
			Key:   "description",
			Value: &model.Description,
		},
	}

	node.Tags = []*nod.Tag{}
	for _, label := range model.Labels {
		node.Tags = append(node.Tags, &nod.Tag{
			Name: label,
		})
	}

	node.KV = map[string]*nod.NodeKV{
		"key": {
			Key:       "key",
			ValueText: &model.Key,
		},
	}
	return node, nil
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
		Name:        "Test Model",
		Active:      true,
		Description: "This is a test model with adapter.",
		Labels:      []string{"label1", "label2"},
		Key:         "value",
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
	require.Equal(t, original.Description, retrievedModel.Description)
	require.ElementsMatch(t, original.Labels, retrievedModel.Labels)
	require.Equal(t, original.Key, retrievedModel.Key)
}
