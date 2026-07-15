package nod

type TagIsNilError struct {
}

func (e *TagIsNilError) Error() string {
	return "tag is nil"
}

func NewTagIsNilError() *TagIsNilError {
	return &TagIsNilError{}
}
