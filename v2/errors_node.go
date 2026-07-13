package nod

type NodeIsNilError struct {
}

func (e *NodeIsNilError) Error() string {
	return "node is nil"
}

func NewNodeIsNilError() *NodeIsNilError {
	return &NodeIsNilError{}
}
