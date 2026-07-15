package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EdgeScope is a generic struct that provides methods for managing edges in a repository.
type EdgeScope[T any] struct {
	repository *Repository
}

// Edges returns an EdgeScope for the given repository, allowing for operations on edges.
func (repository *Repository) Edges() *EdgeScope[Edge] {
	return &EdgeScope[Edge]{
		repository: repository,
	}
}

// Edges is a generic function that returns an EdgeScope for the given repository, allowing for operations on edges of type T.
func Edges[T any](repository *Repository) *EdgeScope[T] {
	return &EdgeScope[T]{
		repository: repository,
	}
}

// SaveEdge saves the given edge to the repository.
func (scope *EdgeScope[T]) SaveEdge(model *T) (string, error) {
	if model == nil {
		return "", NewEdgeIsNilError()
	}

	edge, err := edgeFromModel(scope.repository.adapters, model)
	if err != nil {
		return "", err
	}

	id := ensureEdgeID(edge)
	err = scope.repository.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&edge.Core).Error
	})
	return id, err
}

// DeleteEdge deletes the given edge from the repository.
func (scope *EdgeScope[T]) DeleteEdge(model *T) error {
	if model == nil {
		return NewEdgeIsNilError()
	}

	edge, err := edgeFromModel(scope.repository.adapters, model)
	if err != nil {
		return err
	}
	return scope.repository.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&edge.Core).Error
	})
}

func (scope *EdgeScope[T]) GetEdge(id string) (*T, error) {
	edge := &Edge{}
	if err := scope.repository.db.First(&edge.Core, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return modelFromEdge[T](scope.repository.adapters, edge)
}

func ensureEdgeID(edge *Edge) string {
	if edge.Core.Id == "" {
		edge.Core.Id = uuid.New().String()
	}
	return edge.Core.Id
}
