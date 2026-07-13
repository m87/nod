package nod

import "reflect"

// NodeMapper defines the interface for converting between a domain model T and a Node.
type NodeMapper[T any] interface {
	ToNode(*T) (*Node, error)
	FromNode(*Node) (*T, error)
	IsApplicable(*Node) bool
}

type anyNodeMapper interface {
	toNode(any) (*Node, error)
	fromNode(*Node) (any, error)
	isApplicable(*Node) bool
}

type erasedNodeMapper[T any] struct {
	mapper NodeMapper[T]
}

func (mapper erasedNodeMapper[T]) toNode(model any) (*Node, error) {
	typed, ok := model.(*T)
	if !ok {
		return nil, NewMapperInputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return mapper.mapper.ToNode(typed)
}

func (mapper erasedNodeMapper[T]) fromNode(node *Node) (any, error) {
	return mapper.mapper.FromNode(node)
}

func (mapper erasedNodeMapper[T]) isApplicable(node *Node) bool {
	return mapper.mapper.IsApplicable(node)
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
