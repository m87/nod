package contract

import (
	"testing"
	"time"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeAndEdgeKV(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("SaveAndReadEdgeWithKV", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		createEdgeNode(t, repo, "edge-source")
		createEdgeNode(t, repo, "edge-target")

		edge := &nod.EdgeCore{
			Id:          "edge",
			NamespaceId: ptr("namespace"),
			SourceId:    "edge-source",
			TargetId:    "edge-target",
			Name:        "ingredient",
			Kind:        "contains",
			Status:      "active",
		}
		require.NoError(t, repo.DB().Create(edge).Error)

		valueTime := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
		edgeKV := &nod.EdgeKV{
			EdgeId:      edge.Id,
			Key:         "amount",
			ValueText:   ptr("150 g"),
			ValueNumber: ptr(150.5),
			ValueInt:    ptr(150),
			ValueInt64:  ptr(int64(150)),
			ValueBool:   ptr(true),
			ValueTime:   ptr(valueTime),
		}
		require.NoError(t, repo.DB().Create(edgeKV).Error)

		var foundEdge nod.EdgeCore
		require.NoError(t, repo.DB().First(&foundEdge, "id = ?", edge.Id).Error)
		require.Equal(t, edge.NamespaceId, foundEdge.NamespaceId)
		require.Equal(t, edge.SourceId, foundEdge.SourceId)
		require.Equal(t, edge.TargetId, foundEdge.TargetId)
		require.Equal(t, edge.Name, foundEdge.Name)
		require.Equal(t, edge.Kind, foundEdge.Kind)
		require.Equal(t, edge.Status, foundEdge.Status)
		require.False(t, foundEdge.CreatedAt.IsZero())
		require.False(t, foundEdge.UpdatedAt.IsZero())

		var foundKV nod.EdgeKV
		require.NoError(t, repo.DB().First(&foundKV, "edge_id = ? AND key = ?", edge.Id, edgeKV.Key).Error)
		require.Equal(t, edgeKV.ValueText, foundKV.ValueText)
		require.Equal(t, edgeKV.ValueNumber, foundKV.ValueNumber)
		require.Equal(t, edgeKV.ValueInt, foundKV.ValueInt)
		require.Equal(t, edgeKV.ValueInt64, foundKV.ValueInt64)
		require.Equal(t, edgeKV.ValueBool, foundKV.ValueBool)
		require.Equal(t, edgeKV.ValueTime, foundKV.ValueTime)
	})

	t.Run("ParallelEdgesAndSelfLoop", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		createEdgeNode(t, repo, "parallel-source")
		createEdgeNode(t, repo, "parallel-target")

		edges := []*nod.EdgeCore{
			{Id: "parallel-a", SourceId: "parallel-source", TargetId: "parallel-target", Name: "step-a", Kind: "contains"},
			{Id: "parallel-b", SourceId: "parallel-source", TargetId: "parallel-target", Name: "step-b", Kind: "contains"},
			{Id: "self-loop", SourceId: "parallel-source", TargetId: "parallel-source", Name: "loop", Kind: "depends-on"},
		}
		for _, edge := range edges {
			require.NoError(t, repo.DB().Create(edge).Error)
		}

		var count int64
		require.NoError(t, repo.DB().Model(&nod.EdgeCore{}).Count(&count).Error)
		require.Equal(t, int64(len(edges)), count)
	})

	t.Run("PrimaryKeys", func(t *testing.T) {
		t.Run("EdgeCore", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "pk-source")
			createEdgeNode(t, repo, "pk-target")

			require.NoError(t, repo.DB().Create(&nod.EdgeCore{Id: "duplicate-edge", SourceId: "pk-source", TargetId: "pk-target", Name: "first"}).Error)
			require.Error(t, repo.DB().Create(&nod.EdgeCore{Id: "duplicate-edge", SourceId: "pk-source", TargetId: "pk-target", Name: "second"}).Error)
		})

		t.Run("EdgeKV", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "kv-source")
			createEdgeNode(t, repo, "kv-target")
			createEdgeCore(t, repo, "kv-edge", "kv-source", "kv-target")

			require.NoError(t, repo.DB().Create(&nod.EdgeKV{EdgeId: "kv-edge", Key: "amount", ValueNumber: ptr(100.0)}).Error)
			require.Error(t, repo.DB().Create(&nod.EdgeKV{EdgeId: "kv-edge", Key: "amount", ValueNumber: ptr(200.0)}).Error)
		})
	})

	t.Run("ForeignKeys", func(t *testing.T) {
		t.Run("Source", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "existing-target")
			err := repo.DB().Create(&nod.EdgeCore{Id: "missing-source-edge", SourceId: "missing-source", TargetId: "existing-target"}).Error
			require.Error(t, err)
		})

		t.Run("Target", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "existing-source")
			err := repo.DB().Create(&nod.EdgeCore{Id: "missing-target-edge", SourceId: "existing-source", TargetId: "missing-target"}).Error
			require.Error(t, err)
		})

		t.Run("EdgeKV", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.EdgeKV{EdgeId: "missing-edge", Key: "amount", ValueNumber: ptr(100.0)}).Error
			require.Error(t, err)
		})
	})

	t.Run("Cascades", func(t *testing.T) {
		t.Run("DeleteEdgeDeletesKVAndKeepsNodes", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "edge-delete-source")
			createEdgeNode(t, repo, "edge-delete-target")
			createEdgeCore(t, repo, "edge-delete", "edge-delete-source", "edge-delete-target")
			require.NoError(t, repo.DB().Create(&nod.EdgeKV{EdgeId: "edge-delete", Key: "amount", ValueNumber: ptr(100.0)}).Error)

			require.NoError(t, repo.DB().Delete(&nod.EdgeCore{}, "id = ?", "edge-delete").Error)
			requireTableCount(t, repo, &nod.EdgeKV{}, "edge_id = ?", "edge-delete", 0)
			requireTableCount(t, repo, &nod.NodeCore{}, "id IN ?", []string{"edge-delete-source", "edge-delete-target"}, 2)
		})

		t.Run("DeleteSourceDeletesEdgeAndKV", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "source-delete")
			createEdgeNode(t, repo, "source-delete-target")
			createEdgeCore(t, repo, "source-delete-edge", "source-delete", "source-delete-target")
			require.NoError(t, repo.DB().Create(&nod.EdgeKV{EdgeId: "source-delete-edge", Key: "amount", ValueNumber: ptr(100.0)}).Error)

			require.NoError(t, repo.DB().Delete(&nod.NodeCore{}, "id = ?", "source-delete").Error)
			requireTableCount(t, repo, &nod.EdgeCore{}, "id = ?", "source-delete-edge", 0)
			requireTableCount(t, repo, &nod.EdgeKV{}, "edge_id = ?", "source-delete-edge", 0)
			requireTableCount(t, repo, &nod.NodeCore{}, "id = ?", "source-delete-target", 1)
		})

		t.Run("DeleteTargetDeletesEdgeAndKV", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			createEdgeNode(t, repo, "target-delete-source")
			createEdgeNode(t, repo, "target-delete")
			createEdgeCore(t, repo, "target-delete-edge", "target-delete-source", "target-delete")
			require.NoError(t, repo.DB().Create(&nod.EdgeKV{EdgeId: "target-delete-edge", Key: "amount", ValueNumber: ptr(100.0)}).Error)

			require.NoError(t, repo.DB().Delete(&nod.NodeCore{}, "id = ?", "target-delete").Error)
			requireTableCount(t, repo, &nod.EdgeCore{}, "id = ?", "target-delete-edge", 0)
			requireTableCount(t, repo, &nod.EdgeKV{}, "edge_id = ?", "target-delete-edge", 0)
			requireTableCount(t, repo, &nod.NodeCore{}, "id = ?", "target-delete-source", 1)
		})
	})
}

func createEdgeNode(t *testing.T, repo *nod.Repository, id string) {
	t.Helper()
	require.NoError(t, repo.DB().Create(&nod.NodeCore{Id: id, Name: id, Kind: "test", Status: "active"}).Error)
}

func createEdgeCore(t *testing.T, repo *nod.Repository, id, sourceID, targetID string) {
	t.Helper()
	require.NoError(t, repo.DB().Create(&nod.EdgeCore{Id: id, SourceId: sourceID, TargetId: targetID, Name: id, Kind: "test", Status: "active"}).Error)
}

func requireTableCount(t *testing.T, repo *nod.Repository, model any, query string, value any, expected int64) {
	t.Helper()
	var count int64
	require.NoError(t, repo.DB().Model(model).Where(query, value).Count(&count).Error)
	require.Equal(t, expected, count)
}
