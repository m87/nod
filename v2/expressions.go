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

func coreStringField(name string) StringField  {
	return StringField{
		ref: FieldRef{
			Source: SourceCore,
			Type: ValueTypeString,
			Name: name,
		},
	}	
}

var CoreFields = struct {
	Name StringField
}{
	Name: coreStringField("name"),
}

func (f StringField) Equals(value string) Expression {
	return &comparisionExpression{
		Field: f.ref,
		Operator: OperatorEqual,
		Value: value,
	}
}
