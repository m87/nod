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

type CodecReturnedNilEdgeError struct {
	TypeName string
}

func (e *CodecReturnedNilEdgeError) Error() string {
	return "codec returned a nil edge for type: " + e.TypeName
}

func NewCodecReturnedNilEdgeError(typeName string) *CodecReturnedNilEdgeError {
	return &CodecReturnedNilEdgeError{TypeName: typeName}
}

type EdgeCodecNotApplicableError struct {
	TypeName string
	EdgeId   string
}

func (e *EdgeCodecNotApplicableError) Error() string {
	return "codec for type " + e.TypeName + " is not applicable to edge: " + e.EdgeId
}

func NewEdgeCodecNotApplicableError(typeName, edgeId string) *EdgeCodecNotApplicableError {
	return &EdgeCodecNotApplicableError{
		TypeName: typeName,
		EdgeId:   edgeId,
	}
}
