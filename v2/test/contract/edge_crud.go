package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeCrud(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("BasicEdgeSave", func(t *testing.T) { testBasicEdgeSave(t, factory) })
	t.Run("DeleteEdge", func(t *testing.T) { testDeleteEdge(t, factory) })
	t.Run("DeleteEdgeIfSourceDeleted", func(t *testing.T) { testDeleteEdgeIfSourceDeleted(t, factory) })
	t.Run("DeleteEdgeIfTargetDeleted", func(t *testing.T) { testDeleteEdgeIfTargetDeleted(t, factory) })
}

func createEdgeEndpoints(t *testing.T, repo *nod.Repository) (string, string) {
	t.Helper()
	nodes := repo.Nodes()
	sourceID, err := nodes.SaveNode(&nod.Node{Core: nod.NodeCore{Name: "source", Kind: "test"}})
	require.NoError(t, err)
	targetID, err := nodes.SaveNode(&nod.Node{Core: nod.NodeCore{Name: "target", Kind: "test"}})
	require.NoError(t, err)
	return sourceID, targetID
}
