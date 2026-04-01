package sqlite

import (
	"log/slog"
	"testing"
	"time"

	"github.com/m87/nod"
	"github.com/stretchr/testify/suite"
)

func ptr[T any](v T) *T {
	return &v
}

type RepositoryTestSuite struct {
	suite.Suite
	repo *nod.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	registry := nod.NewMapperRegistry()
	s.repo, _ = NewRepository(":memory:", slog.Default(), registry)
}

func (s *RepositoryTestSuite) TestSaveAndQuery_FullModel() {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:          "test-id",
			ParentId:    ptr("parent-id"),
			NamespaceId: ptr("namespace-id"),
			Name:        "test-node",
			Kind:        "test-kind",
			Status:      "active",
		},
	}

	node.Content = map[string]*nod.Content{
		"key1": {
			Key:   "key1",
			Value: ptr("value1"),
		},
	}

	testTime, _ := time.Parse("2006-12-12 12:12:12", "2006-12-12 12:12:12")
	node.KV = map[string]*nod.KV{
		"kv1": {
			Key:         "kv1",
			ValueText:   ptr("value1"),
			ValueNumber: ptr(42.0),
			ValueBool:   ptr(true),
			ValueTime:   ptr(testTime),
			ValueInt:    ptr(100),
			ValueInt64:  ptr(int64(200)),
		},
	}

	node.Tags = []*nod.Tag{
		{Id: "tag-id-1", Name: "tag1"},
		{Id: "tag-id-2", Name: "tag2"},
	}

	id, err := s.repo.Save(node)
	s.Require().NoError(err)

	foundNode, err := s.repo.Query().NodeId(id).KV().Content().Tags().First()

	s.Require().NoError(err)
	s.Equal(node.Core.Id, foundNode.Core.Id)
	s.Equal(node.Core.ParentId, foundNode.Core.ParentId)
	s.Equal(node.Core.NamespaceId, foundNode.Core.NamespaceId)
	s.Equal(node.Core.Name, foundNode.Core.Name)
	s.Equal(node.Core.Kind, foundNode.Core.Kind)
	s.Equal(node.Core.Status, foundNode.Core.Status)

	s.Require().Len(foundNode.Content, 1)
	content := foundNode.Content["key1"]
	s.NotNil(content)
	s.Equal("key1", content.Key)
	s.Equal("value1", *content.Value)

	s.Require().Len(foundNode.KV, 1)
	kv := foundNode.KV["kv1"]
	s.NotNil(kv)
	s.Equal("kv1", kv.Key)
	s.Equal("value1", *kv.ValueText)
	s.Equal(42.0, *kv.ValueNumber)
	s.Equal(true, *kv.ValueBool)
	s.Equal(testTime, *kv.ValueTime)
	s.Equal(100, *kv.ValueInt)
	s.Equal(int64(200), *kv.ValueInt64)

	s.Require().Len(foundNode.Tags, 2)
	foundTagNames := []string{foundNode.Tags[0].Name, foundNode.Tags[1].Name}
	s.ElementsMatch([]string{"tag1", "tag2"}, foundTagNames)

}

func (s *RepositoryTestSuite) TestGormConstraints_PrimaryKeyNodeCore() {
	first := &nod.NodeCore{
		Id:        "dup-node-id",
		Name:      "node-1",
		Kind:      "kind",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.Db.Create(first).Error
	s.Require().NoError(err)

	second := &nod.NodeCore{
		Id:        "dup-node-id",
		Name:      "node-2",
		Kind:      "kind",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.repo.Db.Create(second).Error
	s.Error(err)
}

func (s *RepositoryTestSuite) TestGormConstraints_CompositePrimaryKeyKV() {
	base := &nod.KV{
		NodeId:    "n-kv-1",
		Key:       "k1",
		ValueText: ptr("v1"),
	}
	err := s.repo.Db.Create(base).Error
	s.Require().NoError(err)

	duplicate := &nod.KV{
		NodeId:    "n-kv-1",
		Key:       "k1",
		ValueText: ptr("v2"),
	}
	err = s.repo.Db.Create(duplicate).Error
	s.Error(err)
}

func (s *RepositoryTestSuite) TestGormConstraints_CompositePrimaryKeyContent() {
	base := &nod.Content{
		NodeId: "n-content-1",
		Key:    "c1",
		Value:  ptr("value-1"),
	}
	err := s.repo.Db.Create(base).Error
	s.Require().NoError(err)

	duplicate := &nod.Content{
		NodeId: "n-content-1",
		Key:    "c1",
		Value:  ptr("value-2"),
	}
	err = s.repo.Db.Create(duplicate).Error
	s.Error(err)
}

func (s *RepositoryTestSuite) TestGormConstraints_NotNullColumns() {
	err := s.repo.Db.Table("node_cores").Create(map[string]any{
		"id":         "null-name-node",
		"name":       nil,
		"kind":       "kind",
		"status":     "active",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}).Error
	s.Error(err)

	err = s.repo.Db.Table("tags").Create(map[string]any{
		"id":         "null-name-tag",
		"name":       nil,
		"created_at": time.Now(),
	}).Error
	s.Error(err)

}

func (s *RepositoryTestSuite) TestMigration_AutoMigrateCreatesTables() {
	migrator := s.repo.Db.Migrator()

	s.True(migrator.HasTable(&nod.NodeCore{}))
	s.True(migrator.HasTable(&nod.Tag{}))
	s.True(migrator.HasTable(&nod.NodeTag{}))
	s.True(migrator.HasTable(&nod.KV{}))
	s.True(migrator.HasTable(&nod.Content{}))

}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
