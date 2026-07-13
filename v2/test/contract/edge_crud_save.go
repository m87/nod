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
	}})
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
}
