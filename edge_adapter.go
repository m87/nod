package nod

// EdgeAdapter defines an interface for converting between a domain model of type T and an Edge. Preferably, the domain model should implement EdgeCodec, but if not, for example if the domain model is a struct from an external library, you can implement this interface in a separate type and register it with the AdapterRegistry.
type EdgeAdapter[T any] interface {
	ToEdge(*T) (*Edge, error)
	FromEdge(*Edge) (*T, error)
	IsApplicable(*Edge) bool
}

type anyEdgeAdapter interface {
	toEdge(any) (*Edge, error)
	fromEdge(*Edge) (any, error)
	isApplicable(*Edge) bool
}

type erasedEdgeAdapter[T any] struct {
	adapter EdgeAdapter[T]
}
