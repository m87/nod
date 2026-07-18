package nod

type AdapterRegistryIsNilError struct{}

func (e *AdapterRegistryIsNilError) Error() string {
	return "adapter registry is nil"
}

func NewAdapterRegistryIsNilError() *AdapterRegistryIsNilError {
	return &AdapterRegistryIsNilError{}
}

type AdapterIsNilError struct {
	TypeName string
}

func (e *AdapterIsNilError) Error() string {
	return "adapter is nil for type: " + e.TypeName
}

func NewAdapterIsNilError(typeName string) *AdapterIsNilError {
	return &AdapterIsNilError{TypeName: typeName}
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

type AdapterNotFoundError struct {
	TypeName string
}

func (e *AdapterNotFoundError) Error() string {
	return "adapter not found for type: " + e.TypeName
}

func NewAdapterNotFoundError(typeName string) *AdapterNotFoundError {
	return &AdapterNotFoundError{TypeName: typeName}
}

type AdapterInputTypeMismatchError struct {
	ExpectedType string
	ActualType   string
}

func (e *AdapterInputTypeMismatchError) Error() string {
	return "adapter input type mismatch: expected " + e.ExpectedType + ", got " + e.ActualType
}

func NewAdapterInputTypeMismatchError(expectedType, actualType string) *AdapterInputTypeMismatchError {
	return &AdapterInputTypeMismatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}

type AdapterNotApplicableError struct {
	TypeName string
	NodeId   string
}

func (e *AdapterNotApplicableError) Error() string {
	return "adapter for type " + e.TypeName + " is not applicable to node: " + e.NodeId
}

func NewAdapterNotApplicableError(typeName, nodeId string) *AdapterNotApplicableError {
	return &AdapterNotApplicableError{
		TypeName: typeName,
		NodeId:   nodeId,
	}
}

type AdapterReturnedNilNodeError struct {
	TypeName string
}

func (e *AdapterReturnedNilNodeError) Error() string {
	return "adapter returned a nil node for type: " + e.TypeName
}

func NewAdapterReturnedNilNodeError(typeName string) *AdapterReturnedNilNodeError {
	return &AdapterReturnedNilNodeError{TypeName: typeName}
}

type CodecReturnedNilNodeError struct {
	TypeName string
}

func (e *CodecReturnedNilNodeError) Error() string {
	return "codec returned a nil node for type: " + e.TypeName
}

func NewCodecReturnedNilNodeError(typeName string) *CodecReturnedNilNodeError {
	return &CodecReturnedNilNodeError{TypeName: typeName}
}

type CodecNotApplicableError struct {
	TypeName string
	NodeId   string
}

func (e *CodecNotApplicableError) Error() string {
	return "codec for type " + e.TypeName + " is not applicable to node: " + e.NodeId
}

func NewCodecNotApplicableError(typeName, nodeId string) *CodecNotApplicableError {
	return &CodecNotApplicableError{
		TypeName: typeName,
		NodeId:   nodeId,
	}
}

type CodecOutputTypeMismatchError struct {
	ExpectedType string
	ActualType   string
}

func (e *CodecOutputTypeMismatchError) Error() string {
	return "codec output type mismatch: expected " + e.ExpectedType + ", got " + e.ActualType
}

func NewCodecOutputTypeMismatchError(expectedType, actualType string) *CodecOutputTypeMismatchError {
	return &CodecOutputTypeMismatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}

type AdapterReturnedNilModelError struct {
	TypeName string
}

func (e *AdapterReturnedNilModelError) Error() string {
	return "adapter returned a nil model for type: " + e.TypeName
}

func NewAdapterReturnedNilModelError(typeName string) *AdapterReturnedNilModelError {
	return &AdapterReturnedNilModelError{TypeName: typeName}
}

type AdapterOutputTypeMismatchError struct {
	ExpectedType string
	ActualType   string
}

func (e *AdapterOutputTypeMismatchError) Error() string {
	return "adapter output type mismatch: expected " + e.ExpectedType + ", got " + e.ActualType
}

func NewAdapterOutputTypeMismatchError(expectedType, actualType string) *AdapterOutputTypeMismatchError {
	return &AdapterOutputTypeMismatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}
