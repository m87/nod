package nod

type InvalidSchemaVersionError struct {
	e error
}

func (e *InvalidSchemaVersionError) Error() string {
	return "invalid schema version: " + e.e.Error()
}

func NewInvalidSchemaVersionError(e error) *InvalidSchemaVersionError {
	return &InvalidSchemaVersionError{
		e: e,
	}
}
