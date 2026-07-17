package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)



func createTestNodes(t *testing.T, repo *nod.Repository) {
	t.Helper()

	repo.Nodes().SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:   "node1",
			Kind:   "test1",
			Status: "active",
		},
		KV: map[string]*nod.NodeKV{
			"key1": {Key: "key1", ValueText: nod.Ptr("value1")},
			"key2": {Key: "key2", ValueText: nod.Ptr("value2")},
		},
	})

	repo.Nodes().SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:   "node2",
			Kind:   "test1",
			Status: "active",
		},
		KV: map[string]*nod.NodeKV{
			"key1": {Key: "key1", ValueText: nod.Ptr("value1")},
			"key2": {Key: "key2", ValueText: nod.Ptr("value4")},
		},
	})

	repo.Nodes().SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:   "node3",
			Kind:   "test3",
			Status: "active",
		},
		KV: map[string]*nod.NodeKV{
			"key1": {Key: "key1", ValueText: nod.Ptr("value1")},
			"key2": {Key: "key2", ValueText: nod.Ptr("value3")},
		},
	})
}

func testFindByCoreAndKv(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()

	createTestNodes(t, repo)

	query := nod.NewNodeQuery(repo)

	nodes, err := query.Where(
		nod.And(
			nod.CoreFields.Kind.Equals("test1"),
			nod.KvString("key1").Equals("value1"),
		),
	).FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 2)
}

func testFindByKv(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()

	createTestNodes(t, repo)

	query := nod.NewNodeQuery(repo)

	nodes, err := query.Where(
		nod.And(
			nod.KvString("key2").Equals("value4"),
		),
	).FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Equal(t, "node2", nodes[0].Core.Name)
	require.Equal(t, "test1", nodes[0].Core.Kind)
	require.Equal(t, "active", nodes[0].Core.Status)
}


func testFindAllNodes(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()

	createTestNodes(t, repo)

	query := nod.NewNodeQuery(repo)

	nodes, err := query.Where(
		nod.And(
			nod.CoreFields.Name.Equals("node2"),
			nod.CoreFields.Kind.Equals("test1"),
			nod.CoreFields.Status.Equals("active"),
		),
	).FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Equal(t, "node2", nodes[0].Core.Name)
	require.Equal(t, "test1", nodes[0].Core.Kind)
	require.Equal(t, "active", nodes[0].Core.Status)
}

func testFindAllNodesWithNoFilter(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()

	createTestNodes(t, repo)

	query := nod.NewNodeQuery(repo)

	nodes, err := query.FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 3)
}	

func testFindMultipleNodes(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer repo.Close()

	createTestNodes(t, repo)

	query := nod.NewNodeQuery(repo)

	nodes, err := query.Where(
		nod.And(
			nod.CoreFields.Kind.Equals("test1"),
			nod.CoreFields.Status.Equals("active"),
		),
	).FindAll()
	require.NoError(t, err)
	require.Len(t, nodes, 2)
}	
