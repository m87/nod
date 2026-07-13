package nod

// NodeMapper defines the interface for converting between a domain model T and a Node.
type NodeMapper[T any] interface {
	ToNode(*T) (*Node, error)
	FromNode(*Node) (*T, error)
	IsApplicable(*Node) bool
}
