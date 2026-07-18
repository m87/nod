package nod

import "time"

type Scope uint8

const (
	ScopeNode Scope = iota
	ScopeEdge
)

type Expression interface {
}

type FiledSource uint8

const (
	SourceCore FiledSource = iota
	SourceKV
	SourceTag
	SourceContent
)

type ValueType uint8

const (
	ValueTypeString ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeBool
	ValueTypeTime
)

type FieldRef struct {
	Source FiledSource
	Type   ValueType
	Name   string
}

type Operator uint8

const (
	OperatorEqual Operator = iota
	OperatorNotEqual
	OperatorGreaterThan
	OperatorLessThan
	OperatorGreaterThanOrEqual
	OperatorLessThanOrEqual
	OperatorIn
	OperatorNotIn
)

type comparisionExpression struct {
	Field    FieldRef
	Operator Operator
	Value    any
}

type orExpression struct {
	Expressions []Expression
}

type andExpression struct {
	Expressions []Expression
}

type notExpression struct {
	Expression Expression
}

type StringField struct {
	ref FieldRef
}

type StringListField struct {
	ref FieldRef
}

type TimeField struct {
	ref FieldRef
}

func coreStringField(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceCore,
			Type:   ValueTypeString,
			Name:   name,
		},
	}
}

func kvString(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type:   ValueTypeString,
			Name:   name,
		},
	}
}

func kvInt(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type:   ValueTypeInt,
			Name:   name,
		},
	}
}

func content(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceContent,
			Type:   ValueTypeString,
			Name:   name,
		},
	}
}

type TagsField struct{}

func Tags() TagsField {
	return TagsField{}
}

func (f TagsField) Has(tagName string) Expression {
	return &comparisionExpression{
		Field: FieldRef{
			Source: SourceTag,
			Type:   ValueTypeString,
			Name:   tagName,
		},
		Operator: OperatorEqual,
		Value:    true,
	}
}

func valuesToAny[T any](values []T) []any {
	result := make([]any, len(values))

	for i, value := range values {
		result[i] = value
	}

	return result
}

var EdgeFields = struct {
	Id          StringField
	Name        StringField
	NamespaceId StringField
	SourceId    StringField
	TargetId    StringField
	Status      StringField
	Kind        StringField
}{
	Id:          coreStringField("id"),
	Name:        coreStringField("name"),
	NamespaceId: coreStringField("namespace_id"),
	SourceId:    coreStringField("source_id"),
	TargetId:    coreStringField("target_id"),
	Status:      coreStringField("status"),
	Kind:        coreStringField("kind"),
}

var NodeFields = struct {
	Id          StringField
	Name        StringField
	NamespaceId StringField
	ParentId    StringField
	Status      StringField
	Kind        StringField
}{
	Id:          coreStringField("id"),
	Name:        coreStringField("name"),
	NamespaceId: coreStringField("namespace_id"),
	ParentId:    coreStringField("parent_id"),
	Status:      coreStringField("status"),
	Kind:        coreStringField("kind"),
}

func KvString(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type:   ValueTypeString,
			Name:   name,
		},
	}
}

// KvTime references a time-valued node or edge KV field.
func KvTime(name string) TimeField {
	return TimeField{
		ref: FieldRef{
			Source: SourceKV,
			Type:   ValueTypeTime,
			Name:   name,
		},
	}
}

func Content(key string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceContent,
			Type:   ValueTypeString,
			Name:   key,
		},
	}
}

func (f StringField) Equals(value string) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorEqual,
		Value:    value,
	}
}

func (f StringField) In(value []string) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorIn,
		Value:    valuesToAny(value),
	}
}

func (f StringField) NotIn(value []string) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorNotIn,
		Value:    valuesToAny(value),
	}
}

func (f TimeField) Equals(value time.Time) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorEqual,
		Value:    value,
	}
}

func (f TimeField) GreaterThan(value time.Time) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorGreaterThan,
		Value:    value,
	}
}

func (f TimeField) LessThan(value time.Time) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorLessThan,
		Value:    value,
	}
}

func (f TimeField) GreaterThanOrEqual(value time.Time) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorGreaterThanOrEqual,
		Value:    value,
	}
}

func (f TimeField) LessThanOrEqual(value time.Time) Expression {
	return &comparisionExpression{
		Field:    f.ref,
		Operator: OperatorLessThanOrEqual,
		Value:    value,
	}
}
