package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testDeleteEdge(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeID, err := repo.Edges().SaveEdge(&nod.Edge{Core: nod.EdgeCore{
		SourceId: sourceID,
		TargetId: targetID,
		Name:     "ingredient",
		Kind:     "contains",
		Status:   "active",
	}})
	require.NoError(t, err)

	err = repo.Edges().DeleteEdge(&nod.Edge{Core: nod.EdgeCore{Id: edgeID}})
	require.NoError(t, err)

	edge, err := repo.Edges().GetEdge(edgeID)
	require.Error(t, err)
	require.Nil(t, edge)
}

func testDeleteEdgeIfSourceDeleted(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeID, err := repo.Edges().SaveEdge(&nod.Edge{Core: nod.EdgeCore{
		SourceId: sourceID,
		TargetId: targetID,
		Name:     "ingredient",
		Kind:     "contains",
		Status:   "active",
	}})
	require.NoError(t, err)

	err = repo.Nodes().DeleteNode(&nod.Node{Core: nod.NodeCore{Id: sourceID}})
	require.NoError(t, err)

	edge, err := repo.Edges().GetEdge(edgeID)
	require.Error(t, err)
	require.Nil(t, edge)
}

func testDeleteEdgeIfTargetDeleted(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeID, err := repo.Edges().SaveEdge(&nod.Edge{Core: nod.EdgeCore{
		SourceId: sourceID,
		TargetId: targetID,
		Name:     "ingredient",
		Kind:     "contains",
		Status:   "active",
	}})
	require.NoError(t, err)

	err = repo.Nodes().DeleteNode(&nod.Node{Core: nod.NodeCore{Id: targetID}})
	require.NoError(t, err)

	edge, err := repo.Edges().GetEdge(edgeID)
	require.Error(t, err)
	require.Nil(t, edge)
}
