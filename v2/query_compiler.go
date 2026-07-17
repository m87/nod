package nod

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



type queryCompiler struct {
	db *gorm.DB
}

func (c queryCompiler) compile(expr Expression) (clause.Expression, error) {
	switch expr := expr.(type) {
	case *comparisionExpression:
		return c.compileComparison(expr)
	case *andExpression:
		return c.compileAnd(expr)
	case *orExpression:
		return c.compileOr(expr)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (c queryCompiler) compileAnd(expr *andExpression) (clause.Expression, error) { 
	if len(expr.Expressions) == 0 {
		return nil, nil
	}

	if len(expr.Expressions) == 1 {
		return c.compile(expr.Expressions[0])
	}

	var clauses []clause.Expression
	for _, e := range expr.Expressions {
		clause, err := c.compile(e)
		if err != nil {
			return nil, err
		}
		if clause != nil {
			clauses = append(clauses, clause)
		}
	}

	if len(clauses) == 0 {
		return nil, nil
	}

	return clause.And(clauses...), nil
}

func (c queryCompiler) compileOr(expr *orExpression) (clause.Expression, error) {
	if len(expr.Expressions) == 0 {
		return nil, nil
	}

	if len(expr.Expressions) == 1 {
		return c.compile(expr.Expressions[0])
	}

	var clauses []clause.Expression
	for _, e := range expr.Expressions {
		clause, err := c.compile(e)
		if err != nil {
			return nil, err
		}
		if clause != nil {
			clauses = append(clauses, clause)
		}
	}

	if len(clauses) == 0 {
		return nil, nil
	}

	return clause.Or(clauses...), nil
}

func (c queryCompiler) compileComparison(expr *comparisionExpression) (clause.Expression, error) {
	fieldName := expr.Field.Name
	value := expr.Value

	switch expr.Operator {
	case OperatorEqual:
		return clause.Eq{Column: fieldName, Value: value}, nil
	case OperatorNotEqual:
		return clause.Neq{Column: fieldName, Value: value}, nil
	case OperatorGreaterThan:
		return clause.Gt{Column: fieldName, Value: value}, nil
	case OperatorLessThan:
		return clause.Lt{Column: fieldName, Value: value}, nil
	case OperatorGreaterThanOrEqual:
		return clause.Gte{Column: fieldName, Value: value}, nil
	case OperatorLessThanOrEqual:
		return clause.Lte{Column: fieldName, Value: value}, nil
	default:
		return nil, fmt.Errorf("unsupported operator: %v", expr.Operator)
	}
}
	
