package nod

type TypedRepository[T any] struct {
	repository *Repository
}

func As[T any](repository *Repository) TypedRepository[T] {
	return TypedRepository[T]{repository: repository}
}

func (tr TypedRepository[T]) Transaction(fn func(repo *TypedRepository[T]) error) error {
	return tr.repository.Transaction(func(txRepo *Repository) error {
		return fn(&TypedRepository[T]{repository: txRepo})
	})
}

func (tr TypedRepository[T]) Save(model *T) error {
	return Save(tr.repository, model)
}

func (tr TypedRepository[T]) Query() *TypedQuery[T] {
	return &TypedQuery[T]{query: tr.repository.Query()}
}

type TypedQuery[T any] struct {
	query *NodeQuery
}

func (tq *TypedQuery[T]) Roots() *TypedQuery[T]                { tq.query.Roots(); return tq }
func (tq *TypedQuery[T]) ExcludeRoot() *TypedQuery[T]          { tq.query.ExcludeRoot(); return tq }
func (tq *TypedQuery[T]) NodeId(id string) *TypedQuery[T]      { tq.query.NodeId(id); return tq }
func (tq *TypedQuery[T]) ParentId(id string) *TypedQuery[T]    { tq.query.ParentId(id); return tq }
func (tq *TypedQuery[T]) NamespaceId(id string) *TypedQuery[T] { tq.query.NamespaceId(id); return tq }
func (tq *TypedQuery[T]) Tags() *TypedQuery[T]                 { tq.query.Tags(); return tq }
func (tq *TypedQuery[T]) KV() *TypedQuery[T]                   { tq.query.KV(); return tq }
func (tq *TypedQuery[T]) Content() *TypedQuery[T]              { tq.query.Content(); return tq }
func (tq *TypedQuery[T]) Limit(n int) *TypedQuery[T]           { tq.query.Limit(n); return tq }
func (tq *TypedQuery[T]) Page(p, size int) *TypedQuery[T]      { tq.query.Page(p, size); return tq }
func (tq *TypedQuery[T]) NameEquals(v string) *TypedQuery[T]   { tq.query.NameEquals(v); return tq }
func (tq *TypedQuery[T]) NameContains(v string) *TypedQuery[T] { tq.query.NameContains(v); return tq }
func (tq *TypedQuery[T]) NameStartsWith(v string) *TypedQuery[T] {
	tq.query.NameStartsWith(v)
	return tq
}
func (tq *TypedQuery[T]) NameEndsWith(v string) *TypedQuery[T] { tq.query.NameEndsWith(v); return tq }
func (tq *TypedQuery[T]) KindEquals(v string) *TypedQuery[T]   { tq.query.KindEquals(v); return tq }
func (tq *TypedQuery[T]) KindContains(v string) *TypedQuery[T] { tq.query.KindContains(v); return tq }
func (tq *TypedQuery[T]) KindStartsWith(v string) *TypedQuery[T] {
	tq.query.KindStartsWith(v)
	return tq
}
func (tq *TypedQuery[T]) KindEndsWith(v string) *TypedQuery[T] { tq.query.KindEndsWith(v); return tq }
func (tq *TypedQuery[T]) StatusEquals(v string) *TypedQuery[T] { tq.query.StatusEquals(v); return tq }
func (tq *TypedQuery[T]) StatusContains(v string) *TypedQuery[T] {
	tq.query.StatusContains(v)
	return tq
}
func (tq *TypedQuery[T]) List() ([]*T, error)   { return ListAs[T](tq.query) }
func (tq *TypedQuery[T]) First() (*T, error)    { return FirstAs[T](tq.query) }
func (tq *TypedQuery[T]) Count() (int64, error) { return tq.query.Count() }
func (tq *TypedQuery[T]) Exists() (bool, error) { return tq.query.Exists() }
func (tq *TypedQuery[T]) Delete() error         { return tq.query.Delete() }

func (tq *TypedQuery[T]) DescendantTree(rootID string) (*TypedTreeNode[T], error) {
	return DescendantTreeAs[T](tq.query, rootID)
}

func (tq *TypedQuery[T]) Descendants(onlyRoots bool) ([]*TypedTreeNode[T], error) {
	return DescendantsAs[T](tq.query, onlyRoots)
}

func (tq *TypedQuery[T]) AncestorTree(childID string) (*TypedTreeNode[T], error) {
	return AncestorTreeAs[T](tq.query, childID)
}

func (tq *TypedQuery[T]) Ancestors() ([]*TypedTreeNode[T], error) {
	return AncestorsAs[T](tq.query)
}
