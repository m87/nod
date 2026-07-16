package nod

import (
	"reflect"

	"gorm.io/gorm"
)

// Operator defines the type of comparison to be used in filtering.
type Operator string

const (
	OperatorEquals     Operator = "equals"
	OperatorLike       Operator = "like"
	OperatorPrefix     Operator = "prefix"
	OperatorSuffix     Operator = "suffix"
	OperatorGreaterThan Operator = "greater_than"
	OperatorLessThan    Operator = "less_than"
)

// Filter represents a generic filter with a value, operator, and type information.
type Filter[T comparable] struct {
	Value    T
	Operator Operator
	Type		 string
}

// Equals creates a Filter that checks for equality with the specified value.
func Equals[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorEquals, Type: reflect.TypeOf(value).String()}
}

// Like creates a Filter that checks if the value is like the specified value (e.g., for string matching).
func Like[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorLike, Type: reflect.TypeOf(value).String()}
}

// Prefix creates a Filter that checks if the value has the specified prefix.
func Prefix[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorPrefix, Type: reflect.TypeOf(value).String()}
}

// Suffix creates a Filter that checks if the value has the specified suffix.
func Suffix[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorSuffix, Type: reflect.TypeOf(value).String()}
}

// GreaterThan creates a Filter that checks if the value is greater than the specified value.
func GreaterThan[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorGreaterThan, Type: reflect.TypeOf(value).String()}
}

// LessThan creates a Filter that checks if the value is less than the specified value.
func LessThan[T comparable](value T) *Filter[T] {
	return &Filter[T]{Value: value, Operator: OperatorLessThan, Type: reflect.TypeOf(value).String()}
}

// Gt is a shorthand for GreaterThan, creating a Filter that checks if the value is greater than the specified value.
func Gt[T comparable](value T) *Filter[T] {
	return GreaterThan(value)
}

// Lt is a shorthand for LessThan, creating a Filter that checks if the value is less than the specified value.
func Lt[T comparable](value T) *Filter[T] {
	return LessThan(value)
}

// Eq is a shorthand for Equals, creating a Filter that checks for equality with the specified value.
func Eq[T comparable](value T) *Filter[T] {
	return Equals(value)
}

func Between[T comparable](min, max T) []*Filter[T] {
	return []*Filter[T]{
		{Value: min, Operator: OperatorGreaterThan, Type: reflect.TypeOf(min).String()},
		{Value: max, Operator: OperatorLessThan, Type: reflect.TypeOf(max).String()},
	}
}

// escapeLike escapes special characters in a string for use in a SQL LIKE clause.
func escapeLike(s string) string {
	r := ""
	for _, c := range s {
		if c == '%' || c == '_' || c == '\\' {
			r += "\\"
		}
		r += string(c)
	}
	return r
}

// applyStringFilter applies a string filter to the specified column in the GORM database query.
func applyStringFilter(db *gorm.DB, column string, filter *Filter[string]) *gorm.DB {
	switch filter.Operator {
	case OperatorEquals:
		db = db.Where(column+" = ?", filter.Value)
	case OperatorLike:
		db = db.Where(column+" LIKE ?", "%"+escapeLike(filter.Value)+"%")
	case OperatorPrefix:
		db = db.Where(column+" LIKE ?", escapeLike(filter.Value)+"%")
	case OperatorSuffix:
		db = db.Where(column+" LIKE ?", "%"+escapeLike(filter.Value))
	case OperatorGreaterThan:
		db = db.Where(column+" > ?", filter.Value)
	case OperatorLessThan:
		db = db.Where(column+" < ?", filter.Value)
	}
	return db
}	
