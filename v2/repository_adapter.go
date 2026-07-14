package nod

import (
	"reflect"
	"sync"
)

// AdapterRegistry stores node and edge adapters by domain model type.
type AdapterRegistry struct {
	mu          sync.RWMutex
	nodeMappers map[reflect.Type]anyNodeAdapter
	edgeMappers map[reflect.Type]anyEdgeMapper
}

// Create a new AdapterRegistry instance.
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		nodeMappers: make(map[reflect.Type]anyNodeAdapter),
		edgeMappers: make(map[reflect.Type]anyEdgeMapper),
	}
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
