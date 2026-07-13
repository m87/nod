package nod

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
func (scope *EdgeScope[T]) SaveEdge(edge *T) error {
	if edge == nil {
		return NewEdgeIsNilError()
	}

	return nil
}

// DeleteEdge deletes the given edge from the repository.
func (scope *EdgeScope[T]) DeleteEdge(edge *T) error {
	if edge == nil {
		return NewEdgeIsNilError()
	}

	return nil
}

func (scope *EdgeScope[T]) GetEdge(id string) (*T, error) {
	// Placeholder for actual implementation to retrieve an edge by ID.
	return nil, nil
}
