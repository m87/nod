package nod

// NodeQuery represents a query for nodes in the repository, allowing for filtering based on various criteria.
type NodeQuery struct {
	repository *Repository

	nodeIds      []string
	parentIds    []string
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

// NewNodeQuery creates a new NodeQuery for the given repository.
func NewNodeQuery(repository *Repository) *NodeQuery {
	return &NodeQuery{
		repository: repository,
	}
}

// Clone creates a copy of the NodeQuery, allowing for modifications without affecting the original query.
func (q *NodeQuery) Clone() *NodeQuery {
	return &NodeQuery{
		repository:   q.repository,
		nodeIds:      append([]string{}, q.nodeIds...),
		parentIds:    append([]string{}, q.parentIds...),
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

func (q *NodeQuery) NodeIds(ids ...string) *NodeQuery {
	q.nodeIds = append(q.nodeIds, ids...)
	return q
}

func (q *NodeQuery) NodeId(id string) *NodeQuery {
	q.nodeIds = append(q.nodeIds, id)
	return q
}

func (q *NodeQuery) ParentIds(ids ...string) *NodeQuery {
	q.parentIds = append(q.parentIds, ids...)
	return q
}

func (q *NodeQuery) ParentId(id string) *NodeQuery {
	q.parentIds = append(q.parentIds, id)
	return q
}

func (q *NodeQuery) NamespaceIds(ids ...string) *NodeQuery {
	q.namespaceIds = append(q.namespaceIds, ids...)
	return q
}

func (q *NodeQuery) NamespaceId(id string) *NodeQuery {
	q.namespaceIds = append(q.namespaceIds, id)
	return q
}

func (q *NodeQuery) Name(filter *StringFilter) *NodeQuery {
	q.name = filter
	return q
}

func (q *NodeQuery) Status(filter *StringFilter) *NodeQuery {
	q.status = filter
	return q
}

func (q *NodeQuery) Kind(filter *StringFilter) *NodeQuery {
	q.kind = filter
	return q
}

func (q *NodeQuery) CreatedDate(filter *TimeFilter) *NodeQuery {
	q.createdDate = filter
	return q
}

func (q *NodeQuery) UpdatedDate(filter *TimeFilter) *NodeQuery {
	q.updatedDate = filter
	return q
}

func (q *NodeQuery) Limit(limit int) *NodeQuery {
	q.limit = limit
	return q
}

func (q *NodeQuery) Page(page int) *NodeQuery {
	q.page = page
	return q
}

func (q *NodeQuery) PageSize(pageSize int) *NodeQuery {
	q.pageSize = pageSize
	return q
}
