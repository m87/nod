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

func testFullNodeSave(t *testing.T, factory RepositoryFactory) {
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
		Tags: []*nod.Tag{
			{Name: "tag1"},
			{Name: "tag2"},
		},
		// KV: map[string]*nod.NodeKV{
		// 	"key1": {Key: "key1", Value: "value1"},
		// 	"key2": {Key: "key2", Value: "value2"},
		// },
		Content: map[string]*nod.NodeContent{
			"content1": {Key: "content1", Value: nod.Ptr("content value 1")},
			"content2": {Key: "content2", Value: nod.Ptr("content value 2")},
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

	// Check tags
	require.Len(t, savedNode.Tags, 2)
	tagNames := []string{savedNode.Tags[0].Name, savedNode.Tags[1].Name}
	require.Contains(t, tagNames, "tag1")
	require.Contains(t, tagNames, "tag2")

	// Check key-value attributes
	// require.Len(t, savedNode.KV, 2)
	// require.Equal(t, "value1", *savedNode.KV["key1"].ValueText)
	// require.Equal(t, "value2", *savedNode.KV["key2"].ValueText)

	// Check content
	require.Len(t, savedNode.Content, 2)
	require.Equal(t, "content value 1", string(*savedNode.Content["content1"].Value))
	require.Equal(t, "content value 2", string(*savedNode.Content["content2"].Value))
}
