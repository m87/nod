package nod

// NodeScope is a generic struct that provides methods for managing nodes of type T within a repository.
type NodeScope[T any] struct {
	repository *Repository
}

// Nodes returns a NodeScope for the given repository, allowing for operations on nodes.
func (repository *Repository) Nodes() *NodeScope[Node] {
	return &NodeScope[Node]{
		repository: repository,
	}
}

// Nodes is a generic function that returns a NodeScope for the given repository, allowing for operations on nodes of type T.
func Nodes[T any](repository *Repository) *NodeScope[T] {
	return &NodeScope[T]{
		repository: repository,
	}
}

// SaveNode saves the given node to the repository.
func (scope *NodeScope[T]) SaveNode(node *T) error {
	if node == nil {
		return NewNodeIsNilError()
	}

	return nil
}

// DeleteNode deletes the given node from the repository.
func (scope *NodeScope[T]) DeleteNode(node *T) error {
	if node == nil {
		return NewNodeIsNilError()
	}

	return nil
}

func (scope *NodeScope[T]) GetNode(id string) (*T, error) {
	// Placeholder for actual implementation to retrieve a node by ID.
	return nil, nil
}
