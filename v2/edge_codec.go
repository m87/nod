package nod

type EdgeCodec interface {
	ToEdge() (*Edge, error)
	FromEdge(*Edge) error
	IsApplicable(*Edge) bool
}
