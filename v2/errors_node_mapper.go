package nod

type MapperRegistryIsNilError struct{}

func (e *MapperRegistryIsNilError) Error() string {
	return "mapper registry is nil"
}

func NewMapperRegistryIsNilError() *MapperRegistryIsNilError {
	return &MapperRegistryIsNilError{}
}

type MapperIsNilError struct {
	TypeName string
}

func (e *MapperIsNilError) Error() string {
	return "mapper is nil for type: " + e.TypeName
}

func NewMapperIsNilError(typeName string) *MapperIsNilError {
	return &MapperIsNilError{TypeName: typeName}
}

type ModelIsNilError struct {
	TypeName string
}

func (e *ModelIsNilError) Error() string {
	return "model is nil for type: " + e.TypeName
}

func NewModelIsNilError(typeName string) *ModelIsNilError {
	return &ModelIsNilError{TypeName: typeName}
}

type MapperNotFoundError struct {
	TypeName string
}

func (e *MapperNotFoundError) Error() string {
	return "mapper not found for type: " + e.TypeName
}

func NewMapperNotFoundError(typeName string) *MapperNotFoundError {
	return &MapperNotFoundError{TypeName: typeName}
}

type MapperInputTypeMismatchError struct {
	ExpectedType string
	ActualType   string
}

func (e *MapperInputTypeMismatchError) Error() string {
	return "mapper input type mismatch: expected " + e.ExpectedType + ", got " + e.ActualType
}

func NewMapperInputTypeMismatchError(expectedType, actualType string) *MapperInputTypeMismatchError {
	return &MapperInputTypeMismatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}

type MapperNotApplicableError struct {
	TypeName string
	NodeId   string
}

func (e *MapperNotApplicableError) Error() string {
	return "mapper for type " + e.TypeName + " is not applicable to node: " + e.NodeId
}

func NewMapperNotApplicableError(typeName, nodeId string) *MapperNotApplicableError {
	return &MapperNotApplicableError{
		TypeName: typeName,
		NodeId:   nodeId,
	}
}

type MapperReturnedNilNodeError struct {
	TypeName string
}

func (e *MapperReturnedNilNodeError) Error() string {
	return "mapper returned a nil node for type: " + e.TypeName
}

func NewMapperReturnedNilNodeError(typeName string) *MapperReturnedNilNodeError {
	return &MapperReturnedNilNodeError{TypeName: typeName}
}

type MapperReturnedNilModelError struct {
	TypeName string
}

func (e *MapperReturnedNilModelError) Error() string {
	return "mapper returned a nil model for type: " + e.TypeName
}

func NewMapperReturnedNilModelError(typeName string) *MapperReturnedNilModelError {
	return &MapperReturnedNilModelError{TypeName: typeName}
}

type MapperOutputTypeMismatchError struct {
	ExpectedType string
	ActualType   string
}

func (e *MapperOutputTypeMismatchError) Error() string {
	return "mapper output type mismatch: expected " + e.ExpectedType + ", got " + e.ActualType
}

func NewMapperOutputTypeMismatchError(expectedType, actualType string) *MapperOutputTypeMismatchError {
	return &MapperOutputTypeMismatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}
