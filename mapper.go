package nod

type NodeMapper[T any] interface {
	ToNode(*T) (*Node, error)
	FromNode(*Node) (*T, error)
}
