package contract

import (
	"testing"
	"time"

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

	err := repo.DB().AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{}, &nod.KV{}, &nod.Content{})
	require.NoError(t, err)
}

func testTagRepositoryDelete(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer closeRepo(t, repo)

	err := repo.DB().Create(&nod.NodeCore{Id: "node", Name: "node", Kind: "kind", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
	require.NoError(t, err)

	err = repo.DB().Create(&nod.Tag{Id: "tag", Name: "tag", CreatedAt: time.Now()}).Error
	require.NoError(t, err)

	err = repo.DB().Create(&nod.NodeTag{NodeId: "node", TagId: "tag"}).Error
	require.NoError(t, err)

	tagRepo := &nod.TagRepository{DB: repo.DB()}
	err = tagRepo.Delete("tag")
	require.NoError(t, err)

	var nodeCount int64
	err = repo.DB().Model(&nod.NodeCore{}).Where("id = ?", "node").Count(&nodeCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), nodeCount)

	var nodeTagCount int64
	err = repo.DB().Model(&nod.NodeTag{}).Where("node_id = ?", "node").Count(&nodeTagCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), nodeTagCount)

	var tagCount int64
	err = repo.DB().Model(&nod.Tag{}).Where("id = ?", "tag").Count(&tagCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), tagCount)
}
