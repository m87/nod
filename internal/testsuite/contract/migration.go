package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testMigration(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer closeRepo(t, repo)

	migrator := repo.DB().Migrator()
	require.True(t, migrator.HasTable(&nod.NodeCore{}))
	require.True(t, migrator.HasTable(&nod.Tag{}))
	require.True(t, migrator.HasTable(&nod.NodeTag{}))
	require.True(t, migrator.HasTable(&nod.KV{}))
	require.True(t, migrator.HasTable(&nod.Content{}))
	require.True(t, migrator.HasTable(&nod.EdgeCore{}))
	require.True(t, migrator.HasTable(&nod.EdgeKV{}))
	require.True(t, migrator.HasTable(&nod.Property{}))

	for model, indexes := range map[any][]string{
		&nod.EdgeCore{}: {
			"idx_edge_namespace_id",
			"idx_edge_source_id",
			"idx_edge_target_id",
			"idx_edge_namespace_source_kind",
			"idx_edge_namespace_target_kind",
			"idx_edge_source_kind_name",
		},
		&nod.EdgeKV{}: {
			"idx_edge_kv_edge_id",
			"idx_edge_kv_key",
		},
	} {
		for _, index := range indexes {
			require.Truef(t, migrator.HasIndex(model, index), "missing index %s", index)
		}
	}

	err := repo.DB().AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{}, &nod.EdgeCore{}, &nod.EdgeKV{})
	require.NoError(t, err)

	var properties []nod.Property
	err = repo.DB().Find(&properties).Error
	require.NoError(t, err)
	require.Len(t, properties, 1)
	require.Equal(t, "version", properties[0].Key)
	require.Equal(t, "2", properties[0].Value)
}
