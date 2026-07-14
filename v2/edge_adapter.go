package nod

type EdgeAdapter[T any] interface {
	ToEdge(*T) (*Edge, error)
	FromEdge(*Edge) (*T, error)
	IsApplicable(*Edge) bool
}
