package nod

// EdgeQuery represents a query for edges in the repository, allowing for filtering based on various criteria.
type EdgeQuery struct {
	repository *Repository

	edgeIds      []string
	sourceIds    []string
	targetIds    []string
	namespaceIds []string
	name         *StringFilter
	status       *StringFilter
	kind         *StringFilter
	createdDate  *TimeFilter
	updatedDate  *TimeFilter
	limit        int
	page         int
	pageSize     int
}

// NewEdgeQuery creates a new EdgeQuery for the given repository.
func NewEdgeQuery(repository *Repository) *EdgeQuery {
	return &EdgeQuery{
		repository: repository,
	}
}

// Clone creates a copy of the EdgeQuery, allowing for modifications without affecting the original query.
func (q *EdgeQuery) Clone() *EdgeQuery {
	return &EdgeQuery{
		repository:   q.repository,
		edgeIds:      append([]string{}, q.edgeIds...),
		sourceIds:    append([]string{}, q.sourceIds...),
		targetIds:    append([]string{}, q.targetIds...),
		namespaceIds: append([]string{}, q.namespaceIds...),
		name:         q.name,
		status:       q.status,
		kind:         q.kind,
		createdDate:  q.createdDate,
		updatedDate:  q.updatedDate,
		limit:        q.limit,
		page:         q.page,
		pageSize:     q.pageSize,
	}
}

func (q *EdgeQuery) EdgeIds(ids ...string) *EdgeQuery {
	q.edgeIds = append(q.edgeIds, ids...)
	return q
}

func (q *EdgeQuery) EdgeId(id string) *EdgeQuery {
	q.edgeIds = append(q.edgeIds, id)
	return q
}

func (q *EdgeQuery) SourceIds(ids ...string) *EdgeQuery {
	q.sourceIds = append(q.sourceIds, ids...)
	return q
}

func (q *EdgeQuery) SourceId(id string) *EdgeQuery {
	q.sourceIds = append(q.sourceIds, id)
	return q
}

func (q *EdgeQuery) TargetIds(ids ...string) *EdgeQuery {
	q.targetIds = append(q.targetIds, ids...)
	return q
}

func (q *EdgeQuery) TargetId(id string) *EdgeQuery {
	q.targetIds = append(q.targetIds, id)
	return q
}

func (q *EdgeQuery) NamespaceIds(ids ...string) *EdgeQuery {
	q.namespaceIds = append(q.namespaceIds, ids...)
	return q
}

func (q *EdgeQuery) NamespaceId(id string) *EdgeQuery {
	q.namespaceIds = append(q.namespaceIds, id)
	return q
}

func (q *EdgeQuery) Name(filter *StringFilter) *EdgeQuery {
	q.name = filter
	return q
}

func (q *EdgeQuery) Status(filter *StringFilter) *EdgeQuery {
	q.status = filter
	return q
}

func (q *EdgeQuery) Kind(filter *StringFilter) *EdgeQuery {
	q.kind = filter
	return q
}

func (q *EdgeQuery) CreatedDate(filter *TimeFilter) *EdgeQuery {
	q.createdDate = filter
	return q
}

func (q *EdgeQuery) UpdatedDate(filter *TimeFilter) *EdgeQuery {
	q.updatedDate = filter
	return q
}

func (q *EdgeQuery) Limit(limit int) *EdgeQuery {
	q.limit = limit
	return q
}

func (q *EdgeQuery) Page(page int) *EdgeQuery {
	q.page = page
	return q
}

func (q *EdgeQuery) PageSize(pageSize int) *EdgeQuery {
	q.pageSize = pageSize
	return q
}
