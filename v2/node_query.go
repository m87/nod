package nod

import "gorm.io/gorm"

// NodeQuery represents a query for nodes in the repository, allowing for filtering based on various criteria.
type NodeQuery struct {
	repository *Repository
	where Expression
}

// NewNodeQuery creates a new NodeQuery for the given repository.
func NewNodeQuery(repository *Repository) *NodeQuery {
	return &NodeQuery{
		repository: repository,
	}
}

func And(exprs ...Expression) Expression {
	if len(exprs) == 0 {
		return nil
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	return &andExpression{
		Expressions: exprs,
	}
}

func Or(exprs ...Expression) Expression {
	if len(exprs) == 0 {
		return nil
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	return &orExpression{
		Expressions: exprs,
	}
}

func (q *NodeQuery) Where(expr Expression) *NodeQuery {
	if expr == nil {
		return q
	}
	
	if q.where == nil {
		q.where = expr
	} else {
		q.where = And(q.where, expr)
	}
	return q
}

func (q *NodeQuery) FindAll() ([]*Node, error) {
	var cores []*NodeCore
	db := q.repository.db

	var err error
	if q.where != nil {
		db, err = applyExpression(db, q.where)
		if err != nil {
			return nil, err
		}
	}

	result := db.Find(&cores)
	
	var nodes []*Node
	for _, core := range cores {
		node := &Node{
			Core: *core,
		}
		nodes = append(nodes, node)
	}

	return nodes, result.Error
}

func applyExpression(db *gorm.DB, expr Expression) (*gorm.DB, error) {
	compiler := queryCompiler{db: db}
	clauseExpr, err := compiler.compile(expr)
	if err != nil {
		panic(err) 	
	}

	return db.Where(clauseExpr), nil
}





