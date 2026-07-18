package nod

type NodeContentIsNilError struct {
}

func (e *NodeContentIsNilError) Error() string {
	return "new node content is nil"
}

func NewNodeContentIsNilError() *NodeContentIsNilError {
	return &NodeContentIsNilError{}
}

type EdgeContentIsNilError struct {
}

func (e *EdgeContentIsNilError) Error() string {
	return "new edge content is nil"
}

func NewEdgeContentIsNilError() *EdgeContentIsNilError {
	return &EdgeContentIsNilError{}
}
