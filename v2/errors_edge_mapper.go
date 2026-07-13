package nod

type EdgeMapperNotApplicableError struct {
	TypeName string
	EdgeId   string
}

func (e *EdgeMapperNotApplicableError) Error() string {
	return "mapper for type " + e.TypeName + " is not applicable to edge: " + e.EdgeId
}

func NewEdgeMapperNotApplicableError(typeName, edgeId string) *EdgeMapperNotApplicableError {
	return &EdgeMapperNotApplicableError{
		TypeName: typeName,
		EdgeId:   edgeId,
	}
}

type MapperReturnedNilEdgeError struct {
	TypeName string
}

func (e *MapperReturnedNilEdgeError) Error() string {
	return "mapper returned a nil edge for type: " + e.TypeName
}

func NewMapperReturnedNilEdgeError(typeName string) *MapperReturnedNilEdgeError {
	return &MapperReturnedNilEdgeError{TypeName: typeName}
}
