package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testNodeDelete(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	nodeID, err := repo.Nodes().SaveNode(&nod.Node{Core: nod.NodeCore{Name: "test", Kind: "test"}})
	require.NoError(t, err)

	err = repo.Nodes().DeleteNode(&nod.Node{Core: nod.NodeCore{Id: nodeID}})
	require.NoError(t, err)

	_, err = repo.Nodes().GetNode(nodeID)
	require.Error(t, err)
}
