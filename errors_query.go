package nod

import "strconv"

type ExpressionIsNilError struct{}

func (e *ExpressionIsNilError) Error() string {
	return "expression is nil"
}

func NewExpressionIsNilError() *ExpressionIsNilError {
	return &ExpressionIsNilError{}
}

type UnsupportedExpressionTypeError struct {
	TypeName string
}

func (e *UnsupportedExpressionTypeError) Error() string {
	return "unsupported expression type: " + e.TypeName
}

func NewUnsupportedExpressionTypeError(typeName string) *UnsupportedExpressionTypeError {
	return &UnsupportedExpressionTypeError{TypeName: typeName}
}

type UnsupportedScopeError struct {
	Scope Scope
}

func (e *UnsupportedScopeError) Error() string {
	return "unsupported scope: " + strconv.Itoa(int(e.Scope))
}

func NewUnsupportedScopeError(scope Scope) *UnsupportedScopeError {
	return &UnsupportedScopeError{Scope: scope}
}
