package nod


type EdgeQuery struct {
	repository *Repository
	where Expression
}

func NewEdgeQuery(repository *Repository) *EdgeQuery {
	return &EdgeQuery{
		repository: repository,
	}
}

func (q *EdgeQuery) Where(expr Expression) *EdgeQuery {
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

func (q *EdgeQuery) FindAll() ([]*Edge, error) {
	var cores []*EdgeCore
	db := q.repository.db

	var err error
	if q.where != nil {
		db, err = applyExpression(db, q.where, ScopeEdge)
		if err != nil {
			return nil, err
		}
	}

	err = db.Find(&cores).Error
	if err != nil {
		return nil, err
	}

	var edges []*Edge
	for _, core := range cores {
		edge := &Edge{
			Core: *core,
		}
		edges = append(edges, edge)
	}

	return edges, nil
}
