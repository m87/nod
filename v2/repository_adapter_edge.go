package nod

import "reflect"

// EdgeMapper defines the interface for converting between a domain model T and an Edge.
type EdgeMapper[T any] interface {
	ToEdge(*T) (*Edge, error)
	FromEdge(*Edge) (*T, error)
	IsApplicable(*Edge) bool
}

type anyEdgeMapper interface {
	toEdge(any) (*Edge, error)
	fromEdge(*Edge) (any, error)
	isApplicable(*Edge) bool
}

type erasedEdgeMapper[T any] struct {
	mapper EdgeMapper[T]
}

func (mapper erasedEdgeMapper[T]) toEdge(model any) (*Edge, error) {
	typed, ok := model.(*T)
	if !ok {
		return nil, NewAdapterInputTypeMismatchError(pointerTypeName[T](), valueTypeName(model))
	}
	return mapper.mapper.ToEdge(typed)
}

func (mapper erasedEdgeMapper[T]) fromEdge(edge *Edge) (any, error) {
	return mapper.mapper.FromEdge(edge)
}

func (mapper erasedEdgeMapper[T]) isApplicable(edge *Edge) bool {
	return mapper.mapper.IsApplicable(edge)
}

// RegisterEdgeMapper registers an edge mapper for a specific type T in the registry.
func RegisterEdgeMapper[T any](registry *AdapterRegistry, mapper EdgeMapper[T]) error {
	if registry == nil {
		return NewAdapterRegistryIsNilError()
	}
	if isNilValue(mapper) {
		return NewAdapterIsNilError(modelTypeName[T]())
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.edgeMappers == nil {
		registry.edgeMappers = make(map[reflect.Type]anyEdgeMapper)
	}
	registry.edgeMappers[reflect.TypeFor[T]()] = &erasedEdgeMapper[T]{mapper: mapper}
	return nil
}

func edgeFromModel[T any](registry *AdapterRegistry, model *T) (*Edge, error) {
	if model == nil {
		return nil, NewModelIsNilError(modelTypeName[T]())
	}
	if edge, ok := any(model).(*Edge); ok {
		return edge, nil
	}

	mapper, err := edgeMapperFor[T](registry)
	if err != nil {
		return nil, err
	}
	edge, err := mapper.toEdge(model)
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

	mapper, err := edgeMapperFor[T](registry)
	if err != nil {
		return nil, err
	}
	if !mapper.isApplicable(edge) {
		return nil, NewEdgeAdapterNotApplicableError(modelTypeName[T](), edge.Core.Id)
	}

	model, err := mapper.fromEdge(edge)
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

func edgeMapperFor[T any](registry *AdapterRegistry) (anyEdgeMapper, error) {
	t := reflect.TypeFor[T]()
	if registry == nil {
		return nil, NewAdapterRegistryIsNilError()
	}

	registry.mu.RLock()
	mapper := registry.edgeMappers[t]
	registry.mu.RUnlock()
	if mapper == nil {
		return nil, NewAdapterNotFoundError(t.String())
	}
	return mapper, nil
}
