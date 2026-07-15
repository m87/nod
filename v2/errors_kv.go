package nod

type NodeKVIsNilError struct {
}

func (e *NodeKVIsNilError) Error() string {
	return "node kv is nil"
}

func NewNodeKVIsNilError() *NodeKVIsNilError {
	return &NodeKVIsNilError{}
}

type EdgeKVIsNilError struct {
}

func (e *EdgeKVIsNilError) Error() string {
	return "edge kv is nil"
}

func NewEdgeKVIsNilError() *EdgeKVIsNilError {
	return &EdgeKVIsNilError{}
}
