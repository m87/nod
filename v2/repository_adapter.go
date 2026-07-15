package nod

import (
	"reflect"
	"sync"
)

// AdapterRegistry stores node and edge adapters by domain model type.
type AdapterRegistry struct {
	mu           sync.RWMutex
	nodeAdapters map[reflect.Type]anyNodeAdapter
	edgeAdapters map[reflect.Type]anyEdgeAdapter
}

// Create a new AdapterRegistry instance.
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		nodeAdapters: make(map[reflect.Type]anyNodeAdapter),
		edgeAdapters: make(map[reflect.Type]anyEdgeAdapter),
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
