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
	switch expr.Field.Source {
	case SourceCore:
		return c.compileCoreComparison(expr)
	case SourceKV:
		return c.compileKVComparison(expr)
	// case SourceContent:
	// 	return c.compileContentComparison(expr)
	// case SourceTag:
	// 	return c.compileTagComparison(expr)
	default:
		return nil, fmt.Errorf("unsupported field source: %v", expr.Field.Source)
	}
}

func (c queryCompiler) compileCoreComparison(expr *comparisionExpression) (clause.Expression, error) {
	column := clause.Column{Table: "node_cores", Name: expr.Field.Name}
	return compileScalarComparison(column, expr.Operator, expr.Value)
}

func kvColumnName(fieldType ValueType) (string, error) {
	switch fieldType {
	case ValueTypeString:
		return "value_text", nil
	case ValueTypeInt:
		return "value_int", nil
	default:
		return "", fmt.Errorf("unsupported KV field type: %v", fieldType)
	}
}

func (c queryCompiler) compileKVComparison(expr *comparisionExpression) (clause.Expression, error) {
	columnName, err := kvColumnName(expr.Field.Type)
	if err != nil {
		return nil, err
	}

	column := clause.Column{Table: "node_kvs", Name: columnName}
	scalarComperison, err := compileScalarComparison(column, expr.Operator, expr.Value)
	if err != nil {
		return nil, err
	}
	subquery := c.db.Session(&gorm.Session{NewDB: true}).Table("node_kvs").Select("1").Where("node_kvs.node_id = node_cores.id").Where("node_kvs.key = ?", expr.Field.Name).Where(scalarComperison)
	return clause.Expr{SQL: "EXISTS (?)", Vars: []interface{}{subquery}}, nil
}

// func (c queryCompiler) compileContentComparison(expr *comparisionExpression) (clause.Expression, error) {}
//
// func (c queryCompiler) compileTagComparison(expr *comparisionExpression) (clause.Expression, error) {}

func compileScalarComparison(column clause.Column, operator Operator, value any) (clause.Expression, error) {
	switch operator {
	case OperatorEqual:
		return clause.Eq{Column: column, Value: value}, nil
	case OperatorNotEqual:
		return clause.Neq{Column: column, Value: value}, nil
	case OperatorGreaterThan:
		return clause.Gt{Column: column, Value: value}, nil
	case OperatorLessThan:
		return clause.Lt{Column: column, Value: value}, nil
	case OperatorGreaterThanOrEqual:
		return clause.Gte{Column: column, Value: value}, nil
	case OperatorLessThanOrEqual:
		return clause.Lte{Column: column, Value: value}, nil
	case OperatorIn:
		values, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("value for 'in' operator must be a slice")
		}
		return clause.IN{Column: column, Values: values}, nil
	case OperatorNotIn:
		values, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("value for 'not in' operator must be a slice")
		}
		return clause.Not(clause.IN{Column: column, Values: values}), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %v", operator)
	}
}
