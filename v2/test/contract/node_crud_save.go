package contract

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testNodeSaveWithParent(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	nodeScope := repo.Nodes()

	parentID, err := nodeScope.SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:   "parent",
			Kind:   "test",
			Status: "active",
		},
	})
	require.NoError(t, err)

	id, err := nodeScope.SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:        "node1",
			Kind:        "test",
			Status:      "active",
			NamespaceId: nod.Ptr("namespace"),
			ParentId:    &parentID,
		},
	})
	require.NoError(t, err)

	savedNode, err := nodeScope.GetNode(id)
	require.NoError(t, err)
	require.Equal(t, "node1", savedNode.Core.Name)
	require.Equal(t, "test", savedNode.Core.Kind)
	require.Equal(t, "active", savedNode.Core.Status)
	require.Equal(t, "namespace", *savedNode.Core.NamespaceId)
	require.Equal(t, parentID, *savedNode.Core.ParentId)
	require.NotEmpty(t, savedNode.Core.CreatedAt)
	require.NotEmpty(t, savedNode.Core.UpdatedAt)
	err = uuid.Validate(savedNode.Core.Id)
	require.NoError(t, err)
}

func testBasicNodeSave(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	defer repo.Close()

	nodeScope := repo.Nodes()

	id, err := nodeScope.SaveNode(&nod.Node{
		Core: nod.NodeCore{
			Name:        "node1",
			Kind:        "test",
			Status:      "active",
			NamespaceId: nod.Ptr("namespace"),
		},
	})
	require.NoError(t, err)

	savedNode, err := nodeScope.GetNode(id)
	require.NoError(t, err)
	require.Equal(t, "node1", savedNode.Core.Name)
	require.Equal(t, "test", savedNode.Core.Kind)
	require.Equal(t, "active", savedNode.Core.Status)
	require.Equal(t, "namespace", *savedNode.Core.NamespaceId)
	require.NotEmpty(t, savedNode.Core.CreatedAt)
	require.NotEmpty(t, savedNode.Core.UpdatedAt)
	err = uuid.Validate(savedNode.Core.Id)
	require.NoError(t, err)
}
