package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryRegressions(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("FirstDoesNotMutateQueryLimit", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		for _, id := range []string{"first-a", "first-b"} {
			_, err := repo.Save(&nod.Node{
				Core: nod.NodeCore{Id: id, Name: "first-same", Kind: "kind"},
			})
			require.NoError(t, err)
		}

		query := repo.Query().NameEquals("first-same")
		_, err := query.First()
		require.NoError(t, err)

		nodes, err := query.List()
		require.NoError(t, err)
		require.Len(t, nodes, 2)
	})

	t.Run("DeleteRespectsLimit", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		for _, id := range []string{"delete-a", "delete-b", "delete-c"} {
			_, err := repo.Save(&nod.Node{
				Core: nod.NodeCore{Id: id, Name: "delete-limited", Kind: "kind"},
			})
			require.NoError(t, err)
		}

		err := repo.Query().NameEquals("delete-limited").Limit(1).Delete()
		require.NoError(t, err)

		count, err := repo.Query().NameEquals("delete-limited").Count()
		require.NoError(t, err)
		require.Equal(t, int64(2), count)
	})

	t.Run("MultipleKVFiltersMatchAcrossRowsForSameNode", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		_, err := repo.Save(&nod.Node{
			Core: nod.NodeCore{Id: "kv-match", Name: "kv-match", Kind: "kind"},
			KV: map[string]*nod.KV{
				"color": {Key: "color", ValueText: ptr("red")},
				"size":  {Key: "size", ValueText: ptr("xl")},
			},
		})
		require.NoError(t, err)

		_, err = repo.Save(&nod.Node{
			Core: nod.NodeCore{Id: "kv-miss", Name: "kv-miss", Kind: "kind"},
			KV: map[string]*nod.KV{
				"color": {Key: "color", ValueText: ptr("red")},
				"size":  {Key: "size", ValueText: ptr("s")},
			},
		})
		require.NoError(t, err)

		nodes, err := repo.Query().
			KVFilter(&nod.KVFilter{Key: ptr("color"), TextEquals: ptr("red")}).
			KVFilter(&nod.KVFilter{Key: ptr("size"), TextEquals: ptr("xl")}).
			List()
		require.NoError(t, err)
		require.Len(t, nodes, 1)
		require.Equal(t, "kv-match", nodes[0].Core.Id)
	})

	t.Run("KVLikeFiltersEscapeWildcards", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		_, err := repo.Save(&nod.Node{
			Core: nod.NodeCore{Id: "literal-percent", Name: "literal-percent", Kind: "kind"},
			KV: map[string]*nod.KV{
				"label": {Key: "label", ValueText: ptr("100%")},
			},
		})
		require.NoError(t, err)

		_, err = repo.Save(&nod.Node{
			Core: nod.NodeCore{Id: "wildcard-miss", Name: "wildcard-miss", Kind: "kind"},
			KV: map[string]*nod.KV{
				"label": {Key: "label", ValueText: ptr("1000")},
			},
		})
		require.NoError(t, err)

		nodes, err := repo.Query().
			KVFilter(&nod.KVFilter{Key: ptr("label"), TextContains: ptr("%")}).
			List()
		require.NoError(t, err)
		require.Len(t, nodes, 1)
		require.Equal(t, "literal-percent", nodes[0].Core.Id)
	})

	t.Run("TypedQueryFiltersBeforeFirstCountExistsAndDelete", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		err := repo.DB().Create(&nod.NodeCore{Id: "typed-other", Name: "typed-mixed", Kind: "other", Status: "active"}).Error
		require.NoError(t, err)

		nod.RegisterMapper(repo.Mappers(), contractModelMapper{})
		typed := nod.As[contractModel](repo)

		id, err := typed.Save(&contractModel{
			ID:   "typed-match",
			Name: "typed-mixed",
			Note: "note",
			Tag:  "tag",
		})
		require.NoError(t, err)

		found, err := typed.Query().NameEquals("typed-mixed").KV().Tags().First()
		require.NoError(t, err)
		require.Equal(t, id, found.ID)

		count, err := typed.Query().NameEquals("typed-mixed").Count()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		exists, err := typed.Query().NameEquals("typed-mixed").Exists()
		require.NoError(t, err)
		require.True(t, exists)

		err = typed.Query().NameEquals("typed-mixed").Delete()
		require.NoError(t, err)

		remaining, err := repo.Query().NameEquals("typed-mixed").Count()
		require.NoError(t, err)
		require.Equal(t, int64(1), remaining)

		otherExists, err := repo.Query().NodeId("typed-other").Exists()
		require.NoError(t, err)
		require.True(t, otherExists)
	})

	t.Run("TagsWithoutIDReuseNamespaceAndName", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		for _, id := range []string{"tag-node-a", "tag-node-b"} {
			_, err := repo.Save(&nod.Node{
				Core: nod.NodeCore{Id: id, Name: id, Kind: "kind"},
				Tags: []*nod.Tag{
					{Name: "shared-tag"},
				},
			})
			require.NoError(t, err)
		}

		var count int64
		err := repo.DB().Model(&nod.Tag{}).Where("name = ?", "shared-tag").Count(&count).Error
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
	})

	t.Run("NilNodeSaveReturnsError", func(t *testing.T) {
		repo := factory(t)
		defer closeRepo(t, repo)

		_, err := repo.Save(nil)
		require.Error(t, err)
	})
}
