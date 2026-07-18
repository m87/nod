package nod

// TypedNodeQuery represents a node query that decodes matching nodes into models of type T.
type TypedNodeQuery[T any] struct {
	query *NodeQuery
}

// NewTypedNodeQuery creates a typed node query for the given repository.
func NewTypedNodeQuery[T any](repository *Repository) *TypedNodeQuery[T] {
	return &TypedNodeQuery[T]{
		query: NewNodeQuery(repository),
	}
}

// WithKV includes key-value attributes when decoding matching nodes.
func (q *TypedNodeQuery[T]) WithKV() *TypedNodeQuery[T] {
	q.query.WithKV()
	return q
}

// WithContent includes content when decoding matching nodes.
func (q *TypedNodeQuery[T]) WithContent() *TypedNodeQuery[T] {
	q.query.WithContent()
	return q
}

// WithTags includes tags when decoding matching nodes.
func (q *TypedNodeQuery[T]) WithTags() *TypedNodeQuery[T] {
	q.query.WithTags()
	return q
}

// Where adds an expression that matching nodes must satisfy.
func (q *TypedNodeQuery[T]) Where(expr Expression) *TypedNodeQuery[T] {
	q.query.Where(expr)
	return q
}

// FindAll returns all matching nodes decoded into models of type T.
func (q *TypedNodeQuery[T]) FindAll() ([]*T, error) {
	nodes, err := q.query.FindAll()
	if err != nil {
		return nil, err
	}

	models := make([]*T, 0, len(nodes))
	for _, node := range nodes {
		model, err := modelFromNode[T](q.query.repository.adapters, node)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

// FindFirst returns the first matching node decoded into T or
// gorm.ErrRecordNotFound when no node matches the query.
func (q *TypedNodeQuery[T]) FindFirst() (*T, error) {
	node, err := q.query.FindFirst()
	if err != nil {
		return nil, err
	}
	return modelFromNode[T](q.query.repository.adapters, node)
}

// DeleteAll deletes every node matching the query.
func (q *TypedNodeQuery[T]) DeleteAll() error {
	return q.query.DeleteAll()
}
