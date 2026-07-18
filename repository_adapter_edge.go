package nod

import "reflect"

func (adapter erasedEdgeAdapter[T]) toEdge(model any) (*Edge, error) {
	typed, ok := model.(*T)
	if !ok {
		return nil, NewAdapterInputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return adapter.adapter.ToEdge(typed)
}

func (adapter erasedEdgeAdapter[T]) fromEdge(edge *Edge) (any, error) {
	return adapter.adapter.FromEdge(edge)
}

func (adapter erasedEdgeAdapter[T]) isApplicable(edge *Edge) bool {
	return adapter.adapter.IsApplicable(edge)
}

// RegisterEdgeAdapter registers an edge adapter for a specific type T in the registry.
func RegisterEdgeAdapter[T any](registry *AdapterRegistry, adapter EdgeAdapter[T]) error {
	if registry == nil {
		return NewAdapterRegistryIsNilError()
	}
	if isNilValue(adapter) {
		return NewAdapterIsNilError(modelTypeName[T]())
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.edgeAdapters == nil {
		registry.edgeAdapters = make(map[reflect.Type]anyEdgeAdapter)
	}
	registry.edgeAdapters[reflect.TypeFor[T]()] = &erasedEdgeAdapter[T]{adapter: adapter}
	return nil
}

func edgeFromModel[T any](registry *AdapterRegistry, model *T) (*Edge, error) {
	if model == nil {
		return nil, NewModelIsNilError(modelTypeName[T]())
	}
	if edge, ok := any(model).(*Edge); ok {
		return edge, nil
	}

	if codec, ok := any(model).(EdgeCodec); ok {
		edge, err := codec.ToEdge()
		if err != nil {
			return nil, err
		}
		if edge == nil {
			return nil, NewCodecReturnedNilEdgeError(modelTypeName[T]())
		}
		return edge, nil
	}

	adapter, err := edgeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	edge, err := adapter.toEdge(model)
	if err != nil {
		return nil, err
	}
	if edge == nil {
		return nil, NewAdapterReturnedNilEdgeError(modelTypeName[T]())
	}
	return edge, nil
}

func modelFromEdge[T any](registry *AdapterRegistry, edge *Edge) (*T, error) {
	if edge == nil {
		return nil, NewEdgeIsNilError()
	}
	if model, ok := any(edge).(*T); ok {
		return model, nil
	}

	if codec, ok := any(new(T)).(EdgeCodec); ok {
		if codec == nil {
			return nil, NewCodecIsNilError(modelTypeName[T]())
		}
		if !codec.IsApplicable(edge) {
			return nil, NewEdgeCodecNotApplicableError(modelTypeName[T](), edge.Core.Id)
		}
		if err := codec.FromEdge(edge); err != nil {
			return nil, err
		}
		typed, ok := any(codec).(*T)
		if !ok {
			return nil, NewCodecOutputTypeMismatchError(pointerTypeName[T](), valueTypeName(codec))
		}
		return typed, nil
	}

	adapter, err := edgeAdapterFor[T](registry)
	if err != nil {
		return nil, err
	}
	if !adapter.isApplicable(edge) {
		return nil, NewEdgeAdapterNotApplicableError(modelTypeName[T](), edge.Core.Id)
	}

	model, err := adapter.fromEdge(edge)
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

func edgeAdapterFor[T any](registry *AdapterRegistry) (anyEdgeAdapter, error) {
	t := reflect.TypeFor[T]()
	if registry == nil {
		return nil, NewAdapterRegistryIsNilError()
	}

	registry.mu.RLock()
	adapter := registry.edgeAdapters[t]
	registry.mu.RUnlock()
	if adapter == nil {
		return nil, NewAdapterNotFoundError(t.String())
	}
	return adapter, nil
}
