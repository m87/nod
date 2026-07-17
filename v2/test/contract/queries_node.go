package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)






func testFindAllNodes(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()


	repo.Nodes().SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:   "node1",
			Kind:   "test",
			Status: "active",
		},
	})
	
	query := nod.NewNodeQuery(repo)

	nodes, err := query.Where(nod.CoreFields.Name.Equals("node1")).FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Equal(t, "node1", nodes[0].Core.Name)
	require.Equal(t, "test", nodes[0].Core.Kind)
	require.Equal(t, "active", nodes[0].Core.Status)

}
