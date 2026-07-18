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

func testDeleteEdgeRelatedData(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	sourceID, targetID := createEdgeEndpoints(t, repo)
	edgeID, err := repo.Edges().SaveEdge(&nod.Edge{
		Core: nod.EdgeCore{
			SourceId: sourceID,
			TargetId: targetID,
			Name:     "ingredient",
			Kind:     "contains",
			Status:   "active",
		},
		Tags: []*nod.Tag{{Name: "required"}},
		KV: map[string]*nod.EdgeKV{
			"quantity": {Key: "quantity", ValueText: nod.Ptr("2")},
		},
		Content: map[string]*nod.EdgeContent{
			"note": {Key: "note", Value: nod.Ptr("sifted")},
		},
	})
	require.NoError(t, err)

	require.NoError(t, repo.Edges().DeleteEdge(&nod.Edge{Core: nod.EdgeCore{Id: edgeID}}))

	var kvCount, contentCount, tagBindingCount int64
	require.NoError(t, repo.DB().Model(&nod.EdgeKV{}).Where("edge_id = ?", edgeID).Count(&kvCount).Error)
	require.NoError(t, repo.DB().Model(&nod.EdgeContent{}).Where("edge_id = ?", edgeID).Count(&contentCount).Error)
	require.NoError(t, repo.DB().Model(&nod.EdgeTag{}).Where("edge_id = ?", edgeID).Count(&tagBindingCount).Error)
	require.Zero(t, kvCount)
	require.Zero(t, contentCount)
	require.Zero(t, tagBindingCount)
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
