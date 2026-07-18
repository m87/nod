package sqlite

type ForeignKeysDisabledError struct{}

func (e *ForeignKeysDisabledError) Error() string {
	return "sqlite foreign keys are disabled"
}

func NewForeignKeysDisabledError() *ForeignKeysDisabledError {
	return &ForeignKeysDisabledError{}
}
