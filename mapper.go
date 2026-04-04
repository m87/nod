package nod

// NodeMapper defines the interface for converting between a domain model T and a Node.
type NodeMapper[T any] interface {
	// ToNode converts a model instance to a Node representation.
	ToNode(*T) (*Node, error)
	// FromNode reconstructs a model instance from a Node.
	FromNode(*Node) (*T, error)
	// IsApplicable returns true if the given Node can be converted by this mapper.
	IsApplicable(*Node) bool
}
