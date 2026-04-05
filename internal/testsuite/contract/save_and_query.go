package contract

import (
	"testing"
	"time"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testSaveAndQueryFullModel(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer closeRepo(t, repo)

	_, err := repo.Save(&nod.Node{
		Core: nod.NodeCore{
			Id:     "parent-id",
			Name:   "parent-node",
			Kind:   "kind",
			Status: "active",
		},
	})
	require.NoError(t, err)

	testTime, err := time.Parse("2006-01-02 15:04:05", "2006-12-12 12:12:12")
	require.NoError(t, err)

	node := &nod.Node{
		Core: nod.NodeCore{
			Id:          "test-id",
			ParentId:    ptr("parent-id"),
			NamespaceId: ptr("namespace-id"),
			Name:        "test-node",
			Kind:        "kind",
			Status:      "active",
		},
		Tags: []*nod.Tag{
			{Id: "tag-id-1", Name: "tag-1"},
			{Id: "tag-id-2", Name: "tag-2"},
		},
		KV: map[string]*nod.KV{
			"kv1": {
				Key:         "kv1",
				ValueText:   ptr("value1"),
				ValueNumber: ptr(42.0),
				ValueBool:   ptr(true),
				ValueTime:   ptr(testTime),
				ValueInt:    ptr(100),
				ValueInt64:  ptr(int64(200)),
			},
		},
		Content: map[string]*nod.Content{
			"content1": {
				Key:   "content1",
				Value: ptr("content-value"),
			},
		},
	}

	id, err := repo.Save(node)
	require.NoError(t, err)

	found, err := repo.Query().NodeId(id).Tags().KV().Content().First()
	require.NoError(t, err)

	require.Equal(t, node.Core.Id, found.Core.Id)
	require.Equal(t, node.Core.ParentId, found.Core.ParentId)
	require.Equal(t, node.Core.NamespaceId, found.Core.NamespaceId)
	require.Equal(t, node.Core.Name, found.Core.Name)
	require.Equal(t, node.Core.Kind, found.Core.Kind)
	require.Equal(t, node.Core.Status, found.Core.Status)

	require.Len(t, found.Tags, 2)
	tagNames := []string{found.Tags[0].Name, found.Tags[1].Name}
	require.ElementsMatch(t, []string{"tag-1", "tag-2"}, tagNames)

	require.Len(t, found.KV, 1)
	require.Equal(t, "value1", *found.KV["kv1"].ValueText)

	require.Len(t, found.Content, 1)
	require.Equal(t, "content-value", *found.Content["content1"].Value)
}
