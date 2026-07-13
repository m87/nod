package nod

import (
	"reflect"
	"sync"
)

type anyNodeMapper interface {
	toNode(any) (*Node, error)
	fromNode(*Node) (any, error)
	isApplicable(*Node) bool
}

type erasedNodeMapper[T any] struct {
	mapper NodeMapper[T]
}

func (e erasedNodeMapper[T]) isApplicable(node *Node) bool {
	return e.mapper.IsApplicable(node)
}

func (e erasedNodeMapper[T]) toNode(v any) (*Node, error) {
	p, ok := v.(*T)
	if !ok {
		return nil, NewMapperInputTypeMismatchError(pointerTypeName[T](), valueTypeName(v))
	}
	return e.mapper.ToNode(p)
}

func (e erasedNodeMapper[T]) fromNode(node *Node) (any, error) {
	return e.mapper.FromNode(node)
}

// MapperRegistry is a registry that holds node mappers for different types.
type MapperRegistry struct {
	mu          sync.RWMutex
	nodeMappers map[reflect.Type]anyNodeMapper
}

// Create a new MapperRegistry instance.
func NewMapperRegistry() *MapperRegistry {
	return &MapperRegistry{
		nodeMappers: make(map[reflect.Type]anyNodeMapper),
	}
}

// RegisterNodeMapper registers a node mapper for a specific type T in the registry.
func RegisterNodeMapper[T any](registry *MapperRegistry, mapper NodeMapper[T]) error {
	if registry == nil {
		return NewMapperRegistryIsNilError()
	}
	if isNilValue(mapper) {
		return NewMapperIsNilError(modelTypeName[T]())
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.nodeMappers == nil {
		registry.nodeMappers = make(map[reflect.Type]anyNodeMapper)
	}
	registry.nodeMappers[reflect.TypeFor[T]()] = &erasedNodeMapper[T]{mapper: mapper}
	return nil
}

// nodeFromModel converts a raw Node directly and uses the registered mapper for domain models.
func nodeFromModel[T any](registry *MapperRegistry, model *T) (*Node, error) {
	if model == nil {
		return nil, NewModelIsNilError(modelTypeName[T]())
	}

	if node, ok := any(model).(*Node); ok {
		return node, nil
	}

	mapper, err := nodeMapperFor[T](registry)
	if err != nil {
		return nil, err
	}

	node, err := mapper.toNode(model)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, NewMapperReturnedNilNodeError(modelTypeName[T]())
	}
	return node, nil
}

// modelFromNode returns a raw Node directly and uses the registered mapper for domain models.
func modelFromNode[T any](registry *MapperRegistry, node *Node) (*T, error) {
	if node == nil {
		return nil, NewNodeIsNilError()
	}

	if model, ok := any(node).(*T); ok {
		return model, nil
	}

	mapper, err := nodeMapperFor[T](registry)
	if err != nil {
		return nil, err
	}
	if !mapper.isApplicable(node) {
		return nil, NewMapperNotApplicableError(modelTypeName[T](), node.Core.Id)
	}

	model, err := mapper.fromNode(node)
	if err != nil {
		return nil, err
	}
	if isNilValue(model) {
		return nil, NewMapperReturnedNilModelError(modelTypeName[T]())
	}

	typed, ok := model.(*T)
	if !ok {
		return nil, NewMapperOutputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return typed, nil
}

func modelTypeName[T any]() string {
	return reflect.TypeFor[T]().String()
}

func pointerTypeName[T any]() string {
	return reflect.TypeFor[*T]().String()
}

func valueTypeName(value any) string {
	if value == nil {
		return "<nil>"
	}
	return reflect.TypeOf(value).String()
}

func isNilValue(value any) bool {
	if value == nil {
		return true
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func nodeMapperFor[T any](registry *MapperRegistry) (anyNodeMapper, error) {
	t := reflect.TypeFor[T]()
	if registry == nil {
		return nil, NewMapperRegistryIsNilError()
	}

	registry.mu.RLock()
	mapper := registry.nodeMappers[t]
	registry.mu.RUnlock()
	if mapper == nil {
		return nil, NewMapperNotFoundError(t.String())
	}
	return mapper, nil
}
