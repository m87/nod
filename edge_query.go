package nod

type EdgeQuery struct {
	repository   *Repository
	where        Expression
	fetchKV      bool
	fetchContent bool
	fetchTags    bool
}

func NewEdgeQuery(repository *Repository) *EdgeQuery {
	return &EdgeQuery{
		repository: repository,
	}
}

func (q *EdgeQuery) WithKV() *EdgeQuery {
	q.fetchKV = true
	return q
}

func (q *EdgeQuery) WithContent() *EdgeQuery {
	q.fetchContent = true
	return q
}

func (q *EdgeQuery) WithTags() *EdgeQuery {
	q.fetchTags = true
	return q
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
	var kvs map[string][]*EdgeKV
	var contents map[string][]*EdgeContent
	var tags map[string][]*Tag
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

	nodeIds := make([]string, 0, len(cores))
	for _, core := range cores {
		nodeIds = append(nodeIds, core.Id)
	}

	if q.fetchKV {
		kvs, err = q.repository.getEdgesKvs(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	if q.fetchContent {
		contents, err = q.repository.getEdgesContents(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	if q.fetchTags {
		tags, err = q.repository.getEdgesTags(nodeIds)
		if err != nil {
			return nil, err
		}
	}

	var edges []*Edge
	for _, core := range cores {
		edge := &Edge{
			Core: *core,
		}

		if q.fetchKV {
			edge.KV = make(map[string]*EdgeKV)
			for _, kv := range kvs[core.Id] {
				edge.KV[kv.Key] = kv
			}
		}

		if q.fetchContent {
			edge.Content = make(map[string]*EdgeContent)
			for _, content := range contents[core.Id] {
				edge.Content[content.Key] = content
			}
		}

		if q.fetchTags {
			edge.Tags = tags[core.Id]
		}

		edges = append(edges, edge)
	}

	return edges, nil
}
