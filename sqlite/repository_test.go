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

func (s *RepositoryTestSuite) TestSaveAndQuery() {
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
		"key1": &nod.Content{
			Key:   "key1",
			Value: ptr("value1"),
		},
	}

	testTime, _ := time.Parse("2006-12-12 12:12:12", "2006-12-12 12:12:12")
	node.KV = map[string]*nod.KV{
		"kv1": &nod.KV{
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
		&nod.Tag{Name: "tag1"},
		&nod.Tag{Name: "tag2"},
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
	s.Equal("tag1", foundNode.Tags[0].Name)
	s.Equal("tag2", foundNode.Tags[1].Name)

}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
