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
	require.True(t, migrator.HasTable(&nod.Property{}))

	err := repo.DB().AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
	require.NoError(t, err)

	var properties []nod.Property
	err = repo.DB().Find(&properties).Error
	require.NoError(t, err)
	require.Len(t, properties, 1)
	require.Equal(t, "version", properties[0].Key)
	require.Equal(t, "1", properties[0].Value)
}
