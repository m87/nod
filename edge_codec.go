package nod

// EdgeCodec defines an interface for converting between a domain model and an Edge.
type EdgeCodec interface {
	ToEdge() (*Edge, error)
	FromEdge(*Edge) error
	IsApplicable(*Edge) bool
}
