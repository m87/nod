package nod

import "reflect"

func (mapper erasedNodeAdapter[T]) toNode(model any) (*Node, error) {
	typed, ok := model.(*T)
	if !ok {
		return nil, NewAdapterInputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return mapper.mapper.ToNode(typed)
}

func (mapper erasedNodeAdapter[T]) fromNode(node *Node) (any, error) {
	return mapper.mapper.FromNode(node)
}

func (mapper erasedNodeAdapter[T]) isApplicable(node *Node) bool {
	return mapper.mapper.IsApplicable(node)
}

// RegisterNodeMapper registers a node mapper for a specific type T in the registry.
func RegisterNodeAdapter[T any](registry *AdapterRegistry, mapper NodeAdapter[T]) error {
	if registry == nil {
		return NewAdapterRegistryIsNilError()
	}
	if isNilValue(mapper) {
		return NewAdapterIsNilError(modelTypeName[T]())
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.nodeMappers == nil {
		registry.nodeMappers = make(map[reflect.Type]anyNodeAdapter)
	}
	registry.nodeMappers[reflect.TypeFor[T]()] = &erasedNodeAdapter[T]{mapper: mapper}
	return nil
}

func nodeFromModel[T any](registry *AdapterRegistry, model *T) (*Node, error) {
	if model == nil {
		return nil, NewModelIsNilError(modelTypeName[T]())
	}
	if node, ok := any(model).(*Node); ok {
		return node, nil
	}

	if codec, ok := any(model).(NodeCodec); ok {
		node, err := codec.ToNode()
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, NewCodecReturnedNilNodeError(modelTypeName[T]())
		}
		return node, nil
	}

	mapper, err := nodeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	node, err := mapper.toNode(model)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, NewAdapterReturnedNilNodeError(modelTypeName[T]())
	}
	return node, nil
}

func modelFromNode[T any](registry *AdapterRegistry, node *Node) (*T, error) {
	if node == nil {
		return nil, NewNodeIsNilError()
	}
	if model, ok := any(node).(*T); ok {
		return model, nil
	}

	if codec, ok := any(new(T)).(NodeCodec); ok {
		err := codec.FromNode(node)
		if err != nil {
			return nil, err
		}
		if !codec.IsApplicable(node) {
			return nil, NewCodecNotApplicableError(modelTypeName[T](), node.Core.Id)
		}
		typed, ok := any(codec).(*T)
		if !ok {
			return nil, NewCodecOutputTypeMismatchError(pointerTypeName[T](), valueTypeName(codec))
		}
		return typed, nil
	}

	mapper, err := nodeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	if !mapper.isApplicable(node) {
		return nil, NewAdapterNotApplicableError(modelTypeName[T](), node.Core.Id)
	}

	model, err := mapper.fromNode(node)
	if err != nil {
		return nil, err
	}
	if isNilValue(model) {
		return nil, NewAdapterReturnedNilModelError(modelTypeName[T]())
	}

	typed, ok := model.(*T)
	if !ok {
		return nil, NewAdapterOutputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return typed, nil
}

func nodeAdapterFor[T any](registry *AdapterRegistry) (anyNodeAdapter, error) {
	t := reflect.TypeFor[T]()
	if registry == nil {
		return nil, NewAdapterRegistryIsNilError()
	}

	registry.mu.RLock()
	mapper := registry.nodeMappers[t]
	registry.mu.RUnlock()
	if mapper == nil {
		return nil, NewAdapterNotFoundError(t.String())
	}
	return mapper, nil
}
