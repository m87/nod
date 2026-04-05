package contract

import (
	"testing"
	"time"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testConstraints(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("PrimaryKeyNodeCore", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.DB().Create(&nod.NodeCore{Id: "dup-node", Name: "n1", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
		require.NoError(t, err)

		err = repo.DB().Create(&nod.NodeCore{Id: "dup-node", Name: "n2", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
		require.Error(t, err)
	})

	t.Run("CompositePrimaryKeyKV", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.DB().Create(&nod.NodeCore{Id: "n-kv", Name: "node-kv", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
		require.NoError(t, err)

		err = repo.DB().Create(&nod.KV{NodeId: "n-kv", Key: "k1", ValueText: ptr("v1")}).Error
		require.NoError(t, err)

		err = repo.DB().Create(&nod.KV{NodeId: "n-kv", Key: "k1", ValueText: ptr("v2")}).Error
		require.Error(t, err)
	})

	t.Run("CompositePrimaryKeyContent", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.DB().Create(&nod.NodeCore{Id: "n-content", Name: "node-content", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
		require.NoError(t, err)

		err = repo.DB().Create(&nod.Content{NodeId: "n-content", Key: "k1", Value: ptr("v1")}).Error
		require.NoError(t, err)

		err = repo.DB().Create(&nod.Content{NodeId: "n-content", Key: "k1", Value: ptr("v2")}).Error
		require.Error(t, err)
	})

	t.Run("NotNullColumns", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.DB().Table("node_cores").Create(map[string]any{
			"id":         "null-name-node",
			"name":       nil,
			"kind":       "kind",
			"status":     "active",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}).Error
		require.Error(t, err)

		err = repo.DB().Table("tags").Create(map[string]any{
			"id":         "null-name-tag",
			"name":       nil,
			"created_at": time.Now(),
		}).Error
		require.Error(t, err)
	})

	t.Run("ForeignKeys", func(t *testing.T) {
		t.Run("NodeParent", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.NodeCore{Id: "child", ParentId: ptr("missing-parent"), Name: "child", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.Error(t, err)
		})

		t.Run("KVNode", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.KV{NodeId: "missing-node", Key: "k", ValueText: ptr("v")}).Error
			require.Error(t, err)
		})

		t.Run("ContentNode", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.Content{NodeId: "missing-node", Key: "k", Value: ptr("v")}).Error
			require.Error(t, err)
		})

		t.Run("NodeTag", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.NodeTag{NodeId: "missing-node", TagId: "missing-tag"}).Error
			require.Error(t, err)
		})

		t.Run("ParentDeleteSetsNull", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.NodeCore{Id: "parent", Name: "parent", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.NodeCore{Id: "child", ParentId: ptr("parent"), Name: "child", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Delete(&nod.NodeCore{}, "id = ?", "parent").Error
			require.NoError(t, err)

			var child nod.NodeCore
			err = repo.DB().First(&child, "id = ?", "child").Error
			require.NoError(t, err)
			require.Nil(t, child.ParentId)
		})

		t.Run("DeleteNodeCascadesChildrenRows", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.KV{NodeId: "node", Key: "k", ValueText: ptr("v")}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.Content{NodeId: "node", Key: "c", Value: ptr("v")}).Error
			require.NoError(t, err)

			err = repo.DB().Delete(&nod.NodeCore{}, "id = ?", "node").Error
			require.NoError(t, err)

			var nodeTagCount int64
			err = repo.DB().Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
			require.NoError(t, err)
			require.Equal(t, int64(0), nodeTagCount)

			var kvCount int64
			err = repo.DB().Model(&nod.KV{}).Where("node_id = ?", "node").Count(&kvCount).Error
			require.NoError(t, err)
			require.Equal(t, int64(0), kvCount)

			var contentCount int64
			err = repo.DB().Model(&nod.Content{}).Where("node_id = ?", "node").Count(&contentCount).Error
			require.NoError(t, err)
			require.Equal(t, int64(0), contentCount)
		})

		t.Run("DeleteTagKeepsNode", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.DB().Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.DB().Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
			require.NoError(t, err)

			err = repo.DB().Delete(&nod.Tag{}, "id = ?", "tag").Error
			require.NoError(t, err)

			var node nod.NodeCore
			err = repo.DB().First(&node, "id = ?", "node").Error
			require.NoError(t, err)

			var nodeTagCount int64
			err = repo.DB().Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
			require.NoError(t, err)
			require.Equal(t, int64(0), nodeTagCount)
		})
	})
}
