package nod

// NodeAdapter defines an interface for converting between a domain model of type T and a Node. Preferably, the domain model should implement NodeCodec, but if not, for example if the domain model is a struct from an external library, you can implement this interface in a separate type and register it with the AdapterRegistry.
type NodeAdapter[T any] interface {
	ToNode(*T) (*Node, error)
	FromNode(*Node) (*T, error)
	IsApplicable(*Node) bool
}

type anyNodeAdapter interface {
	toNode(any) (*Node, error)
	fromNode(*Node) (any, error)
	isApplicable(*Node) bool
}

type erasedNodeAdapter[T any] struct {
	mapper NodeAdapter[T]
}
