package nod

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type NodeMapper[T any] interface {
	ToNode(*T) (*Node, error)       
	FromNode(*Node) (*T, error)    
}

type TypedQuery[T any] struct {
	q      *NodeQuery
	mapper NodeMapper[T]
}

type TypedTreeNode[T any] struct {
	Node     *T
	Children []*TreeNode
}

func NewTypedQuery[T any](db *gorm.DB, log *slog.Logger, mapper NodeMapper[T]) *TypedQuery[T] {
	return &TypedQuery[T]{
		q:      Query(db, log),
		mapper: mapper,
	}
}

func (t *TypedQuery[T]) Clone() *TypedQuery[T] {
	clone := *t
	clone.q.nodeIds = append([]string{}, t.q.nodeIds...)
	clone.q.parentIds = append([]string{}, t.q.parentIds...)
	clone.q.namespaceIds = append([]string{}, t.q.namespaceIds...)
	return &clone
}

func (t *TypedQuery[T]) Roots() *TypedQuery[T] { 
	t.q.Roots()
	return t
}

func (t *TypedQuery[T]) ExcludeRoot() *TypedQuery[T] {
	t.q.ExcludeRoot()
	return t
}

func (t *TypedQuery[T]) NodeId(nodeId string) *TypedQuery[T] {
	t.q.NodeId(nodeId)
	return t
}

func (t *TypedQuery[T]) ParentId(parentId string) *TypedQuery[T] {
	t.q.ParentId(parentId)
	return t
}

func (t *TypedQuery[T]) NamespaceId(namespaceId string) *TypedQuery[T] {
	t.q.NamespaceId(namespaceId)
	return t
}

func (t *TypedQuery[T]) NodeIds(nodeIds []string) *TypedQuery[T] {
	t.q.NodeIds(nodeIds)
	return t
}

func (t *TypedQuery[T]) ParentIds(parentIds []string) *TypedQuery[T] {
	t.q.ParentIds(parentIds)
	return t
}

func (t *TypedQuery[T]) NamespaceIds(namespaceIds []string) *TypedQuery[T] {
	t.q.NamespaceIds(namespaceIds)
	return t
}

func (t *TypedQuery[T]) Tags() *TypedQuery[T] {
	t.q.Tags()
	return t
}

func (t *TypedQuery[T]) KV() *TypedQuery[T] {
	t.q.KV()
	return t
}

func (t *TypedQuery[T]) Content() *TypedQuery[T] {
	t.q.Content()
	return t
}

func (t *TypedQuery[T]) Limit(limit int) *TypedQuery[T] {
	t.q.Limit(limit)
	return t
}

func (t *TypedQuery[T]) Page(page int, pageSize int) *TypedQuery[T] {
	t.q.Page(page, pageSize)
	return t
}

func (t *TypedQuery[T]) Name(filter *StringFilter) *TypedQuery[T] {
	t.q.Name(filter)
	return t
}

func (t *TypedQuery[T]) Type(filter *StringFilter) *TypedQuery[T] {
	t.q.Type(filter)
	return t
}

func (t *TypedQuery[T]) Kind(filter *StringFilter) *TypedQuery[T] {
	t.q.Kind(filter)
	return t
}

func (t *TypedQuery[T]) Status(filter *StringFilter) *TypedQuery[T] {
	t.q.Status(filter)
	return t
}

func (t *TypedQuery[T]) CreatedDate(filter *TimeFilter) *TypedQuery[T] {
	t.q.CreatedDate(filter)
	return t
}

func (t *TypedQuery[T]) UpdatedDate(filter *TimeFilter) *TypedQuery[T] {
	t.q.UpdatedDate(filter)
	return t
}

func TApplyCommonFilters[T any](db *gorm.DB, t *TypedQuery[T]) *gorm.DB {
	if len(t.q.nodeIds) > 0 {
		db = db.Where("id IN ?", t.q.nodeIds)
	}
	if t.q.onlyRoots {
		db = db.Where("parent_id IS NULL or parent_id = \"\"")
	}
	if t.q.excludeRoot {
		db = db.Where("parent_id IS NOT NULL")
	}
	if len(t.q.parentIds) > 0 {
		db = db.Where("parent_id IN ?", t.q.parentIds)
	}
	if len(t.q.namespaceIds) > 0 {
		db = db.Where("namespace_id IN ?", t.q.namespaceIds)
	}
	if t.q.name != nil {
		db = ApplyStringFilter(db, "name", t.q.name)
	}
	if t.q.type_ != nil {
		db = ApplyStringFilter(db, "type", t.q.type_)
	}
	if t.q.kind != nil {
		db = ApplyStringFilter(db, "kind", t.q.kind)
	}
	if t.q.status != nil {
		db = ApplyStringFilter(db, "status", t.q.status)
	}
	if t.q.createdDate != nil {
		db = ApplyTimeFilter(db, "created_at", t.q.createdDate)
	}
	if t.q.updatedDate != nil {
		db = ApplyTimeFilter(db, "updated_at", t.q.updatedDate)
	}
	return db
}

func (t *TypedQuery[T]) ApplyConditions(db *gorm.DB) *gorm.DB {
	t.q.log.Debug(fmt.Sprintf("TypedQuery current filters: nodeIds=%v, parentIds=%v, namespaceIds=%v, name=%v, type_=%v, kind=%v, status=%v, createdDate=%v, updatedDate=%v, onlyRoots=%v, excludeRoot=%v",
		t.q.nodeIds, t.q.parentIds, t.q.namespaceIds, t.q.name, t.q.type_, t.q.kind, t.q.status, t.q.createdDate, t.q.updatedDate, t.q.onlyRoots, t.q.excludeRoot))

	db = TApplyCommonFilters(db, t)

	if t.q.limit > 0 {
		db = db.Limit(t.q.limit)
	}
	if t.q.page > 0 && t.q.pageSize > 0 {
		offset := (t.q.page - 1) * t.q.pageSize
		db = db.Offset(offset).Limit(t.q.pageSize)
	}

	return db
}

func (t *TypedQuery[T]) List() ([]*T, error) {
	nodes, err := t.q.List()
	if err != nil {
		return nil, err
	}
	
	var results []*T
	for _, n := range nodes {
		mapped, err := t.mapper.FromNode(n)
		if err != nil {
			return nil, err
		}
		results = append(results, mapped)
	}
	return results, nil
}

func (t *TypedQuery[T]) First() (*T, error) {
	node, err := t.q.First()
	if err != nil {
		return nil, err
	}
	mapped, err := t.mapper.FromNode(node)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (t *TypedQuery[T]) Count() (int64, error) {
	return t.q.Count()
}

func (t *TypedQuery[T]) Descendants(onlyRoots bool) ([]*TypedTreeNode[T], error) {
	trees := make([]*TypedTreeNode[T], 0)

	nodes, err := t.q.List()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if !onlyRoots && (n.Core.ParentId == nil || *n.Core.ParentId == "") {
			tree, err := t.q.buildTree(n.Core.Id)
			if err != nil {
				return nil, err
			}

			mappedNode, err := t.mapper.FromNode(tree.Node)
			if err != nil {
				return nil, err
			}

			mappedTree := &TypedTreeNode[T]{
				Node:     mappedNode,
				Children: tree.Children,
			}

			trees = append(trees, mappedTree)
		}
	}

	return trees, nil
}

func (t *TypedQuery[T]) DescendantTree(rootID string) (*TypedTreeNode[T], error) {
	tree, err := t.q.buildTree(rootID)
	if err != nil {
		return nil, err
	}

	mappedNode, err := t.mapper.FromNode(tree.Node)
	if err != nil {
		return nil, err
	}

	mappedTree := &TypedTreeNode[T]{
		Node:     mappedNode,
		Children: tree.Children,
	}

	return mappedTree, nil
}

func (t *TypedQuery[T]) Ancestors() ([]*TypedTreeNode[T], error) {
	trees := make([]*TypedTreeNode[T], 0)

	nodes, err := t.q.List()		
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		tree, err := t.q.buildAncestorTree(n.Core.Id)
		if err != nil {
			return nil, err
		}

		mappedNode, err := t.mapper.FromNode(tree.Node)
		if err != nil {
			return nil, err
		}

		mappedTree := &TypedTreeNode[T]{
			Node:     mappedNode,
			Children: tree.Children,
		}

		trees = append(trees, mappedTree)
	}		
	return trees, nil
}

func (t *TypedQuery[T]) AncestorTree(childID string) (*TypedTreeNode[T], error) {
	tree, err := t.q.buildAncestorTree(childID)
	if err != nil {
		return nil, err
	}

	mappedNode, err := t.mapper.FromNode(tree.Node)
	if err != nil {
		return nil, err
	}

	mappedTree := &TypedTreeNode[T]{
		Node:     mappedNode,
		Children: tree.Children,
	}

	return mappedTree, nil
}

func (t *TypedQuery[T]) Delete() error {
	return t.q.Delete()
}

func (t *TypedQuery[T]) HasChildren() bool {
	return t.q.HasChildren()
}

