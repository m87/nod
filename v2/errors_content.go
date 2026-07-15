package nod


type NodeContentIsNilError struct {
}

func (e *NodeContentIsNilError) Error() string {
	return "new node content is nil"
}

func NewNodeContentIsNilError() *NodeContentIsNilError {
	return &NodeContentIsNilError{}
}
