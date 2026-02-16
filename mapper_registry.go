package nod

import (
	"fmt"
	"reflect"
)

type anyMapper interface {
	toNode(any) (*Node, error)
	fromNode(*Node) (any, error)
}

type erasedMapper[T any] struct {
	mapper NodeMapper[T]
}

func (e *erasedMapper[T]) toNode(v any) (*Node, error) {
	p, ok := v.(*T)
	if !ok {
		return nil, fmt.Errorf("mapper expected *%v, got %T", reflect.TypeOf((*T)(nil)).Elem(), v)
	}
	return e.mapper.ToNode(p)
}

func (e *erasedMapper[T]) fromNode(node *Node) (any, error) {
	return e.mapper.FromNode(node)
}

type MapperRegistry struct {
	byType map[reflect.Type]anyMapper
}

func NewMapperRegistry() *MapperRegistry {
	return &MapperRegistry{
		byType: make(map[reflect.Type]anyMapper),
	}
}

func RegisterMapper[T any](registry *MapperRegistry, mapper NodeMapper[T]) *MapperRegistry {
	t := reflect.TypeOf((*T)(nil)).Elem()
	registry.byType[t] = &erasedMapper[T]{mapper: mapper}
	return registry
}

func (r *MapperRegistry) forType(t reflect.Type) (anyMapper, error) {
	mapper, ok := r.byType[t]
	if !ok {
		return nil, fmt.Errorf("no mapper registered for type %v", t)
	}
	return mapper, nil
}
