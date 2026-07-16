package nod

type NodeIsNilError struct {
}

func (e *NodeIsNilError) Error() string {
	return "node is nil"
}

func NewNodeIsNilError() *NodeIsNilError {
	return &NodeIsNilError{}
}

type CodecIsNilError struct {
	modelType string
}

func (e *CodecIsNilError) Error() string {
	return "codec is nil " + e.modelType
}

func NewCodecIsNilError(modelType string) *CodecIsNilError {
	return &CodecIsNilError{modelType: modelType}
}

type MultipleNodesFoundError struct {
}

func (e *MultipleNodesFoundError) Error() string {
	return "multiple nodes found"
}

func NewMultipleNodesFoundError() *MultipleNodesFoundError {
	return &MultipleNodesFoundError{}
}
