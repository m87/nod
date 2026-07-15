package nod

import "reflect"

func (adapter erasedNodeAdapter[T]) toNode(model any) (*Node, error) {
	typed, ok := model.(*T)
	if !ok {
		return nil, NewAdapterInputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return adapter.adapter.ToNode(typed)
}

func (adapter erasedNodeAdapter[T]) fromNode(node *Node) (any, error) {
	return adapter.adapter.FromNode(node)
}

func (adapter erasedNodeAdapter[T]) isApplicable(node *Node) bool {
	return adapter.adapter.IsApplicable(node)
}

// RegisterNodeAdapter registers a node adapter for a specific type T in the registry.
func RegisterNodeAdapter[T any](registry *AdapterRegistry, adapter NodeAdapter[T]) error {
	if registry == nil {
		return NewAdapterRegistryIsNilError()
	}
	if isNilValue(adapter) {
		return NewAdapterIsNilError(modelTypeName[T]())
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.nodeAdapters == nil {
		registry.nodeAdapters = make(map[reflect.Type]anyNodeAdapter)
	}
	registry.nodeAdapters[reflect.TypeFor[T]()] = &erasedNodeAdapter[T]{adapter: adapter}
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

	adapter, err := nodeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	node, err := adapter.toNode(model)
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
		if codec == nil {
			return nil, NewCodecIsNilError(modelTypeName[T]())
		}
		if !codec.IsApplicable(node) {
			return nil, NewCodecNotApplicableError(modelTypeName[T](), node.Core.Id)
		}
		err := codec.FromNode(node)
		if err != nil {
			return nil, err
		}
		typed, ok := any(codec).(*T)
		if !ok {
			return nil, NewCodecOutputTypeMismatchError(pointerTypeName[T](), valueTypeName(codec))
		}
		return typed, nil
	}

	adapter, err := nodeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	if !adapter.isApplicable(node) {
		return nil, NewAdapterNotApplicableError(modelTypeName[T](), node.Core.Id)
	}

	model, err := adapter.fromNode(node)
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
	adapter := registry.nodeAdapters[t]
	registry.mu.RUnlock()
	if adapter == nil {
		return nil, NewAdapterNotFoundError(t.String())
	}
	return adapter, nil
}
