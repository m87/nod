package contract

import (
	"testing"
	"time"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

type RepositoryFactory func(t *testing.T) *nod.Repository

func RunRepositoryContractTests(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("SaveAndQueryFullModel", func(t *testing.T) {
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
	})

	t.Run("Constraints", func(t *testing.T) {
		t.Run("PrimaryKeyNodeCore", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.Db.Create(&nod.NodeCore{Id: "dup-node", Name: "n1", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.Db.Create(&nod.NodeCore{Id: "dup-node", Name: "n2", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.Error(t, err)
		})

		t.Run("CompositePrimaryKeyKV", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.Db.Create(&nod.NodeCore{Id: "n-kv", Name: "node-kv", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.Db.Create(&nod.KV{NodeId: "n-kv", Key: "k1", ValueText: ptr("v1")}).Error
			require.NoError(t, err)

			err = repo.Db.Create(&nod.KV{NodeId: "n-kv", Key: "k1", ValueText: ptr("v2")}).Error
			require.Error(t, err)
		})

		t.Run("CompositePrimaryKeyContent", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.Db.Create(&nod.NodeCore{Id: "n-content", Name: "node-content", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
			require.NoError(t, err)

			err = repo.Db.Create(&nod.Content{NodeId: "n-content", Key: "k1", Value: ptr("v1")}).Error
			require.NoError(t, err)

			err = repo.Db.Create(&nod.Content{NodeId: "n-content", Key: "k1", Value: ptr("v2")}).Error
			require.Error(t, err)
		})

		t.Run("NotNullColumns", func(t *testing.T) {
			repo := factory(t)
			defer closeRepo(t, repo)

			err := repo.Db.Table("node_cores").Create(map[string]any{
				"id":         "null-name-node",
				"name":       nil,
				"kind":       "kind",
				"status":     "active",
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}).Error
			require.Error(t, err)

			err = repo.Db.Table("tags").Create(map[string]any{
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

				err := repo.Db.Create(&nod.NodeCore{Id: "child", ParentId: ptr("missing-parent"), Name: "child", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
				require.Error(t, err)
			})

			t.Run("KVNode", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.KV{NodeId: "missing-node", Key: "k", ValueText: ptr("v")}).Error
				require.Error(t, err)
			})

			t.Run("ContentNode", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.Content{NodeId: "missing-node", Key: "k", Value: ptr("v")}).Error
				require.Error(t, err)
			})

			t.Run("NodeTag", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.NodeTag{NodeId: "missing-node", TagId: "missing-tag"}).Error
				require.Error(t, err)
			})

			t.Run("ParentDeleteSetsNull", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.NodeCore{Id: "parent", Name: "parent", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.NodeCore{Id: "child", ParentId: ptr("parent"), Name: "child", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Delete(&nod.NodeCore{}, "id = ?", "parent").Error
				require.NoError(t, err)

				var child nod.NodeCore
				err = repo.Db.First(&child, "id = ?", "child").Error
				require.NoError(t, err)
				require.Nil(t, child.ParentId)
			})

			t.Run("DeleteNodeCascadesChildrenRows", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.KV{NodeId: "node", Key: "k", ValueText: ptr("v")}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.Content{NodeId: "node", Key: "c", Value: ptr("v")}).Error
				require.NoError(t, err)

				err = repo.Db.Delete(&nod.NodeCore{}, "id = ?", "node").Error
				require.NoError(t, err)

				var nodeTagCount int64
				err = repo.Db.Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
				require.NoError(t, err)
				require.Equal(t, int64(0), nodeTagCount)

				var kvCount int64
				err = repo.Db.Model(&nod.KV{}).Where("node_id = ?", "node").Count(&kvCount).Error
				require.NoError(t, err)
				require.Equal(t, int64(0), kvCount)

				var contentCount int64
				err = repo.Db.Model(&nod.Content{}).Where("node_id = ?", "node").Count(&contentCount).Error
				require.NoError(t, err)
				require.Equal(t, int64(0), contentCount)
			})

			t.Run("DeleteTagKeepsNode", func(t *testing.T) {
				repo := factory(t)
				defer closeRepo(t, repo)

				err := repo.Db.Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
				require.NoError(t, err)

				err = repo.Db.Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
				require.NoError(t, err)

				err = repo.Db.Delete(&nod.Tag{}, "id = ?", "tag").Error
				require.NoError(t, err)

				var node nod.NodeCore
				err = repo.Db.First(&node, "id = ?", "node").Error
				require.NoError(t, err)

				var nodeTagCount int64
				err = repo.Db.Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
				require.NoError(t, err)
				require.Equal(t, int64(0), nodeTagCount)
			})
		})
	})

	t.Run("Migration", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		migrator := repo.Db.Migrator()
		require.True(t, migrator.HasTable(&nod.NodeCore{}))
		require.True(t, migrator.HasTable(&nod.Tag{}))
		require.True(t, migrator.HasTable(&nod.NodeTag{}))
		require.True(t, migrator.HasTable(&nod.KV{}))
		require.True(t, migrator.HasTable(&nod.Content{}))

		err := repo.Db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
		require.NoError(t, err)
	})

	t.Run("TagRepositoryDelete", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.Db.Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
		require.NoError(t, err)

		err = repo.Db.Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
		require.NoError(t, err)

		err = repo.Db.Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
		require.NoError(t, err)

		tagRepo := &nod.TagRepository{DB: repo.Db}
		err = tagRepo.Delete("tag")
		require.NoError(t, err)

		var nodeCount int64
		err = repo.Db.Model(&nod.NodeCore{}).Where("id = ?", "node").Count(&nodeCount).Error
		require.NoError(t, err)
		require.Equal(t, int64(1), nodeCount)

		var nodeTagCount int64
		err = repo.Db.Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
		require.NoError(t, err)
		require.Equal(t, int64(0), nodeTagCount)

		var tagCount int64
		err = repo.Db.Model(&nod.Tag{}).Where("id = ?", "tag").Count(&tagCount).Error
		require.NoError(t, err)
		require.Equal(t, int64(0), tagCount)
	})

	t.Run("TypedRepositorySaveAndQuery", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		nod.RegisterMapper(repo.Mappers, contractModelMapper{})

		typed := nod.As[contractModel](repo)
		model := &contractModel{
			Name: "typed-name",
			Note: "typed-note",
			Tag:  "typed-tag",
		}

		id, err := typed.Save(model)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		found, err := typed.Query().NameEquals("typed-name").KV().Content().Tags().First()
		require.NoError(t, err)
		require.Equal(t, id, found.ID)
		require.Equal(t, "typed-name", found.Name)
		require.Equal(t, "typed-note", found.Note)
		require.Equal(t, "typed-tag", found.Tag)
	})

	t.Run("RepositoryClose", func(t *testing.T) {
		repo := factory(t)
		require.NoError(t, repo.Close())
	})
}

type contractModel struct {
	ID   string
	Name string
	Note string
	Tag  string
}

type contractModelMapper struct{}

func (m contractModelMapper) ToNode(model *contractModel) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:     model.ID,
			Name:   model.Name,
			Kind:   "contract-kind",
			Status: "active",
		},
		Tags: []*nod.Tag{{Name: model.Tag}},
		KV: map[string]*nod.KV{
			"note": {
				Key:       "note",
				ValueText: ptr(model.Note),
			},
		},
		Content: map[string]*nod.Content{
			"note": {
				Key:   "note",
				Value: ptr(model.Note),
			},
		},
	}
	return node, nil
}

func (m contractModelMapper) FromNode(node *nod.Node) (*contractModel, error) {
	model := &contractModel{
		ID:   node.Core.Id,
		Name: node.Core.Name,
	}

	if kv, ok := node.KV["note"]; ok && kv.ValueText != nil {
		model.Note = *kv.ValueText
	}
	if len(node.Tags) > 0 {
		model.Tag = node.Tags[0].Name
	}

	return model, nil
}

func (m contractModelMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "contract-kind"
}

func closeRepo(t *testing.T, repo *nod.Repository) {
	t.Helper()
	if repo == nil {
		return
	}
	require.NoError(t, repo.Close())
}

func ptr[T any](value T) *T {
	return &value
}
