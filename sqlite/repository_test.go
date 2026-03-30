package sqlite

import (
	"log/slog"
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	repo *nod.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	registry := nod.NewMapperRegistry()
	s.repo, _ = NewRepository(":memory:", slog.Default(), registry)
}

func (s *RepositoryTestSuite) TestSaveAndQuery() {
	parentId := "1"
	namespaceId := "2"

	node := &nod.Node{
		Core: nod.NodeCore{
			Name:        "test-node",
			Kind:        "test-kind",
			Status:      "active",
			ParentId:    &parentId,
			NamespaceId: &namespaceId,
		},
	}
	id, err := s.repo.Save(node)
	s.Require().NoError(err)

	q := s.repo.Query().NameEquals("test-node")
	found, err := q.First()
	s.Require().NoError(err)
	s.NotNil(id)
	s.Equal("test-node", found.Core.Name)
	s.Equal("test-kind", found.Core.Kind)
	s.Equal("active", found.Core.Status)
	s.NotEmpty(found.Core.CreatedAt)
	s.NotEmpty(found.Core.UpdatedAt)
	s.Equal(&parentId, found.Core.ParentId)
	s.Equal(&namespaceId, found.Core.NamespaceId)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
