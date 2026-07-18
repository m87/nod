package nod

// NodeCodec defines an interface for converting between a domain model and a Node. Preferably, the domain model should implement this interface directly, but if not, for example if the domain model is a struct from an external library, you can implement this interface in a separate type and register it with the AdapterRegistry.
type NodeCodec interface {
	ToNode() (*Node, error)
	FromNode(*Node) error
	IsApplicable(*Node) bool
}
