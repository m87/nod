package nod

type EdgeAdapterNotApplicableError struct {
	TypeName string
	EdgeId   string
}

func (e *EdgeAdapterNotApplicableError) Error() string {
	return "adapter for type " + e.TypeName + " is not applicable to edge: " + e.EdgeId
}

func NewEdgeAdapterNotApplicableError(typeName, edgeId string) *EdgeAdapterNotApplicableError {
	return &EdgeAdapterNotApplicableError{
		TypeName: typeName,
		EdgeId:   edgeId,
	}
}

type AdapterReturnedNilEdgeError struct {
	TypeName string
}

func (e *AdapterReturnedNilEdgeError) Error() string {
	return "adapter returned a nil edge for type: " + e.TypeName
}

func NewAdapterReturnedNilEdgeError(typeName string) *AdapterReturnedNilEdgeError {
	return &AdapterReturnedNilEdgeError{TypeName: typeName}
}
