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
	return q.find(0)
}

// FindFirst returns the first matching node or gorm.ErrRecordNotFound when no
// node matches the query.
func (q *NodeQuery) FindFirst() (*Node, error) {
	nodes, err := q.find(1)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return nodes[0], nil
}

// DeleteAll deletes every node matching the query. An empty query is rejected
// to prevent accidental deletion of all nodes.
func (q *NodeQuery) DeleteAll() error {
	if q.where == nil {
		return gorm.ErrMissingWhereClause
	}

	return q.repository.Transaction(func(txRepository *Repository) error {
		db, err := applyExpression(txRepository.db, q.where, ScopeNode)
		if err != nil {
			return err
		}
		return db.Delete(&NodeCore{}).Error
	})
}

func (q *NodeQuery) find(limit int) ([]*Node, error) {
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
	if limit > 0 {
		db = db.Limit(limit)
	}

	result := db.Find(&cores)
	if result.Error != nil {
		return nil, result.Error
	}

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

	return nodes, nil
}

func applyExpression(db *gorm.DB, expr Expression, scope Scope) (*gorm.DB, error) {
	compiler := queryCompiler{db: db, scope: scope}
	clauseExpr, err := compiler.compile(expr)
	if err != nil {
		return nil, err
	}

	return db.Where(clauseExpr), nil
}
