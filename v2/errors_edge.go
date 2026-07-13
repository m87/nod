package nod

type EdgeIsNilError struct {
}

func (e *EdgeIsNilError) Error() string {
	return "edge is nil"
}

func NewEdgeIsNilError() *EdgeIsNilError {
	return &EdgeIsNilError{}
}
