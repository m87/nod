package contract

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testBasicEdgeSave(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	id, err := repo.Edges().SaveEdge(&nod.Edge{Core: nod.EdgeCore{
		SourceId: sourceID,
		TargetId: targetID,
		Name:     "ingredient",
		Kind:     "contains",
		Status:   "active",
	},
		KV: map[string]*nod.EdgeKV{
			"quantity": {Key: "quantity", ValueText: nod.Ptr("2")},
			"unit":     {Key: "unit", ValueText: nod.Ptr("cups")},
		},
	})
	require.NoError(t, err)
	require.NoError(t, uuid.Validate(id))

	edge, err := repo.Edges().GetEdge(id)
	require.NoError(t, err)
	require.Equal(t, sourceID, edge.Core.SourceId)
	require.Equal(t, targetID, edge.Core.TargetId)
	require.Equal(t, "ingredient", edge.Core.Name)
	require.Equal(t, "contains", edge.Core.Kind)
	require.Equal(t, "active", edge.Core.Status)
	require.False(t, edge.Core.CreatedAt.IsZero())
	require.False(t, edge.Core.UpdatedAt.IsZero())
	require.Equal(t, "2", *edge.KV["quantity"].ValueText)
	require.Equal(t, "cups", *edge.KV["unit"].ValueText)
}

func testFullEdgeSave(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	namespaceID := "recipes"
	id, err := repo.Edges().SaveEdge(&nod.Edge{
		Core: nod.EdgeCore{
			NamespaceId: &namespaceID,
			SourceId:    sourceID,
			TargetId:    targetID,
			Name:        "ingredient",
			Kind:        "contains",
			Status:      "active",
		},
		Tags: []*nod.Tag{
			{Name: "required"},
			{Name: "food"},
		},
		KV: map[string]*nod.EdgeKV{
			"quantity": {Key: "quantity", ValueText: nod.Ptr("2")},
			"unit":     {Key: "unit", ValueText: nod.Ptr("cups")},
		},
		Content: map[string]*nod.EdgeContent{
			"note":        {Key: "note", Value: nod.Ptr("sifted")},
			"description": {Key: "description", Value: nod.Ptr("all-purpose flour")},
		},
	})
	require.NoError(t, err)

	edge, err := repo.Edges().GetEdge(id)
	require.NoError(t, err)
	require.Equal(t, sourceID, edge.Core.SourceId)
	require.Equal(t, targetID, edge.Core.TargetId)
	require.Equal(t, namespaceID, *edge.Core.NamespaceId)
	require.Equal(t, "ingredient", edge.Core.Name)
	require.Equal(t, "contains", edge.Core.Kind)
	require.Equal(t, "active", edge.Core.Status)
	require.False(t, edge.Core.CreatedAt.IsZero())
	require.False(t, edge.Core.UpdatedAt.IsZero())
	require.NoError(t, uuid.Validate(edge.Core.Id))

	require.Len(t, edge.Tags, 2)
	tagNames := []string{edge.Tags[0].Name, edge.Tags[1].Name}
	require.Contains(t, tagNames, "required")
	require.Contains(t, tagNames, "food")

	require.Len(t, edge.KV, 2)
	require.Equal(t, "2", *edge.KV["quantity"].ValueText)
	require.Equal(t, "cups", *edge.KV["unit"].ValueText)

	require.Len(t, edge.Content, 2)
	require.Equal(t, "sifted", *edge.Content["note"].Value)
	require.Equal(t, "all-purpose flour", *edge.Content["description"].Value)
}
