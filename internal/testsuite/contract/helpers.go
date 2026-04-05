package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type contractModel struct {
	ID   string
	Name string
	Note string
	Tag  string
}

type contractModelMapper struct{}

func (m contractModelMapper) ToNode(model *contractModel) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:     model.ID,
			Name:   model.Name,
			Kind:   "contract-kind",
			Status: "active",
		},
		Tags: []*nod.Tag{{Name: model.Tag}},
		KV: map[string]*nod.KV{
			"note": {
				Key:       "note",
				ValueText: ptr(model.Note),
			},
		},
		Content: map[string]*nod.Content{
			"note": {
				Key:   "note",
				Value: ptr(model.Note),
			},
		},
	}
	return node, nil
}

func (m contractModelMapper) FromNode(node *nod.Node) (*contractModel, error) {
	model := &contractModel{
		ID:   node.Core.Id,
		Name: node.Core.Name,
	}

	if kv, ok := node.KV["note"]; ok && kv.ValueText != nil {
		model.Note = *kv.ValueText
	}
	if len(node.Tags) > 0 {
		model.Tag = node.Tags[0].Name
	}

	return model, nil
}

func (m contractModelMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "contract-kind"
}

func closeRepo(t *testing.T, repo *nod.Repository) {
	t.Helper()
	if repo == nil {
		return
	}
	require.NoError(t, repo.Close())
}

func ptr[T any](value T) *T {
	return &value
}
