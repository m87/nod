package nod

import "gorm.io/gorm"

// NodeQuery represents a query for nodes in the repository, allowing for filtering based on various criteria.
type NodeQuery struct {
	repository   *Repository
	where        Expression
	fetchKV      bool
	fetchContent bool
	fetchTags    bool
}

// NewNodeQuery creates a new NodeQuery for the given repository.
func NewNodeQuery(repository *Repository) *NodeQuery {
	return &NodeQuery{
		repository: repository,
	}
}

func (q *NodeQuery) WithKV() *NodeQuery {
	q.fetchKV = true
	return q
}

func (q *NodeQuery) WithContent() *NodeQuery {
	q.fetchContent = true
	return q
}

func (q *NodeQuery) WithTags() *NodeQuery {
	q.fetchTags = true
	return q
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
	var kv map[string][]*NodeKV
	var contents map[string][]*NodeContent
	var tags map[string][]*Tag

	db := q.repository.db

	var err error
	if q.where != nil {
		db, err = applyExpression(db, q.where, ScopeNode)
		if err != nil {
			return nil, err
		}
	}

	result := db.Find(&cores)

	nodeIds := make([]string, 0, len(cores))
	for _, core := range cores {
		nodeIds = append(nodeIds, core.Id)
	}

	if q.fetchKV {
		kv, err = q.repository.getNodesKvs(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	if q.fetchContent {
		contents, err = q.repository.getNodesContents(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	if q.fetchTags {
		tags, err = q.repository.getNodesTags(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	var nodes []*Node
	for _, core := range cores {
		node := &Node{
			Core: *core,
		}

		if q.fetchKV {
			node.KV = make(map[string]*NodeKV)
			for _, kv := range kv[core.Id] {
				node.KV[kv.Key] = kv
			}
		}

		if q.fetchContent {
			node.Content = make(map[string]*NodeContent)
			for _, content := range contents[core.Id] {
				node.Content[content.Key] = content
			}
		}

		if q.fetchTags {
			node.Tags = tags[core.Id]
		}

		nodes = append(nodes, node)
	}

	return nodes, result.Error
}

func applyExpression(db *gorm.DB, expr Expression, scope Scope) (*gorm.DB, error) {
	compiler := queryCompiler{db: db, scope: scope}
	clauseExpr, err := compiler.compile(expr)
	if err != nil {
		panic(err)
	}

	return db.Where(clauseExpr), nil
}
