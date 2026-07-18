package nod


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
)

type FieldRef struct {
	Source FiledSource
	Type  ValueType
	Name	string
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
	Field FieldRef
	Operator Operator
	Value any
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

func coreStringField(name string) StringField  {
	return StringField{
		ref: FieldRef{
			Source: SourceCore,
			Type: ValueTypeString,
			Name: name,
		},
	}	
}

func kvString(name string) StringField  {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type: ValueTypeString,
			Name: name,
		},
	}	
}

func kvInt(name string) StringField  {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type: ValueTypeInt,
			Name: name,
		},
	}	
}

func content(name string) StringField  {
	return StringField{
		ref: FieldRef{
			Source: SourceContent,
			Type: ValueTypeString,
			Name: name,
		},
	}	
}


type TagsField struct {}

func Tags() TagsField {
	return TagsField{}
}

func (f TagsField) Has(tagName string) Expression {
	return &comparisionExpression{
		Field: FieldRef{
			Source: SourceTag,
			Type: ValueTypeString,
			Name: tagName,
		},
		Operator: OperatorEqual,
		Value: true,
	}
}

func valuesToAny[T any](values []T) []any {
    result := make([]any, len(values))

    for i, value := range values {
        result[i] = value
    }

    return result
}



var CoreFields = struct {
	Id StringField
	Name StringField
	NamespaceId StringField
	ParentId StringField
	Status StringField
	Kind StringField

}{
	Id: coreStringField("id"),
	Name: coreStringField("name"),
	NamespaceId: coreStringField("namespace_id"),
	ParentId: coreStringField("parent_id"),
	Status: coreStringField("status"),
	Kind: coreStringField("kind"),
}

func KvString(name string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceKV,
			Type: ValueTypeString,
			Name: name,
		},
	}
}

func Content(key string) StringField {
	return StringField{
		ref: FieldRef{
			Source: SourceContent,
			Type: ValueTypeString,
			Name: key,
		},
	}
}

func (f StringField) Equals(value string) Expression {
	return &comparisionExpression{
		Field: f.ref,
		Operator: OperatorEqual,
		Value: value,
	}
}

func (f StringField) In(value []string) Expression {
	return &comparisionExpression{
		Field: f.ref,
		Operator: OperatorIn,
		Value: valuesToAny(value),
	}
}

func (f StringField) NotIn(value []string) Expression {
	return &comparisionExpression{
		Field: f.ref,
		Operator: OperatorNotIn,
		Value: valuesToAny(value),
	}
}


