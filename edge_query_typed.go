package nod

// TypedEdgeQuery represents an edge query that decodes matching edges into models of type T.
type TypedEdgeQuery[T any] struct {
	query *EdgeQuery
}

// NewTypedEdgeQuery creates a typed edge query for the given repository.
func NewTypedEdgeQuery[T any](repository *Repository) *TypedEdgeQuery[T] {
	return &TypedEdgeQuery[T]{
		query: NewEdgeQuery(repository),
	}
}

// WithKV includes key-value attributes when decoding matching edges.
func (q *TypedEdgeQuery[T]) WithKV() *TypedEdgeQuery[T] {
	q.query.WithKV()
	return q
}

// WithContent includes content when decoding matching edges.
func (q *TypedEdgeQuery[T]) WithContent() *TypedEdgeQuery[T] {
	q.query.WithContent()
	return q
}

// WithTags includes tags when decoding matching edges.
func (q *TypedEdgeQuery[T]) WithTags() *TypedEdgeQuery[T] {
	q.query.WithTags()
	return q
}

// Where adds an expression that matching edges must satisfy.
func (q *TypedEdgeQuery[T]) Where(expr Expression) *TypedEdgeQuery[T] {
	q.query.Where(expr)
	return q
}

// FindAll returns all matching edges decoded into models of type T.
func (q *TypedEdgeQuery[T]) FindAll() ([]*T, error) {
	edges, err := q.query.FindAll()
	if err != nil {
		return nil, err
	}

	models := make([]*T, 0, len(edges))
	for _, edge := range edges {
		model, err := modelFromEdge[T](q.query.repository.adapters, edge)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}
