package nod

import (
	"fmt"
	"reflect"
)

type anyMapper interface {
	toNode(any) (*Node, error)
	fromNode(*Node) (any, error)
	isApplicable(*Node) bool
}

type anyEdgeMapper interface {
	toEdge(any) (*Edge, error)
	fromEdge(*Edge) (any, error)
	isApplicable(*Edge) bool
}

type erasedEdgeMapper[T any] struct {
	mapper EdgeMapper[T]
}

func (e erasedEdgeMapper[T]) isApplicable(edge *Edge) bool {
	return e.mapper.IsApplicable(edge)
}

func (e erasedEdgeMapper[T]) toEdge(v any) (*Edge, error) {
	p, ok := v.(*T)
	if !ok {
		return nil, fmt.Errorf("mapper expected *%v, got %T", reflect.TypeOf((*T)(nil)).Elem(), v)
	}
	return e.mapper.ToEdge(p)
}

func (e erasedEdgeMapper[T]) fromEdge(edge *Edge) (any, error) {
	return e.mapper.FromEdge(edge)
}

type erasedMapper[T any] struct {
	mapper NodeMapper[T]
}

func (e erasedMapper[T]) isApplicable(node *Node) bool {
	return e.mapper.IsApplicable(node)
}

func (e erasedMapper[T]) toNode(v any) (*Node, error) {
	p, ok := v.(*T)
	if !ok {
		return nil, fmt.Errorf("mapper expected *%v, got %T", reflect.TypeOf((*T)(nil)).Elem(), v)
	}
	return e.mapper.ToNode(p)
}

func (e erasedMapper[T]) fromNode(node *Node) (any, error) {
	return e.mapper.FromNode(node)
}

// MapperRegistry stores type-to-mapper associations for converting between domain models and nodes.
type MapperRegistry struct {
	nodeByType map[reflect.Type]anyMapper
	edgeByType map[reflect.Type]anyEdgeMapper
}

// NewMapperRegistry creates an empty MapperRegistry.
func NewMapperRegistry() *MapperRegistry {
	return &MapperRegistry{
		nodeByType: make(map[reflect.Type]anyMapper),
		edgeByType: make(map[reflect.Type]anyEdgeMapper),
	}
}

// RegisterMapper registers a NodeMapper for type T in the registry.
func RegisterMapper[T any](registry *MapperRegistry, mapper NodeMapper[T]) *MapperRegistry {
	t := reflect.TypeOf((*T)(nil)).Elem()
	registry.nodeByType[t] = &erasedMapper[T]{mapper: mapper}
	return registry
}

func RegisterEdgeMapper[T any](registry *MapperRegistry, mapper EdgeMapper[T]) *MapperRegistry {
	t := reflect.TypeOf((*T)(nil)).Elem()
	registry.edgeByType[t] = &erasedEdgeMapper[T]{mapper: mapper}
	return registry
}

func (r *MapperRegistry) forType(t reflect.Type) (anyMapper, error) {
	mapper, ok := r.nodeByType[t]
	if !ok {
		return nil, fmt.Errorf("no mapper registered for type %v", t)
	}
	return mapper, nil
}
