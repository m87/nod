package nod

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type TypedTreeNode[T any] struct {
	Node     *T
	Children []*TypedTreeNode[T]
}

func AncestorTreeAs[T any](q *NodeQuery, childID string) (*TypedTreeNode[T], error) {
	nodes, err := q.fetchAncestorNodes(childID)
	if err != nil {
		return nil, err
	}
	return buildAncestorTreeFromNodes[T](q, nodes)
}

func AncestorsAs[T any](q *NodeQuery) ([]*TypedTreeNode[T], error) {
	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	out := make([]*TypedTreeNode[T], 0, len(nodes))
	for _, n := range nodes {
		tree, err := AncestorTreeAs[T](q, n.Core.Id)
		if err != nil {
			return nil, err
		}
		out = append(out, tree)
	}
	return out, nil
}

func DescendantTreeAs[T any](q *NodeQuery, rootID string) (*TypedTreeNode[T], error) {
	nodes, err := q.fetchDescendantNodes(rootID)
	if err != nil {
		return nil, err
	}
	return buildTreeFromNodes[T](q, nodes, rootID)
}

func DescendantsAs[T any](q *NodeQuery, onlyRoots bool) ([]*TypedTreeNode[T], error) {
	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	out := make([]*TypedTreeNode[T], 0, len(nodes))
	for _, n := range nodes {
		if onlyRoots && n.Core.ParentId != nil && *n.Core.ParentId != "" {
			continue
		}
		tree, err := DescendantTreeAs[T](q, n.Core.Id)
		if err != nil {
			return nil, err
		}
		out = append(out, tree)
	}
	return out, nil
}

func mapperForT[T any](q *NodeQuery) (anyMapper, reflect.Type, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	m, err := q.mappers.forType(t)
	return m, t, err
}

func buildTreeFromNodes[T any](q *NodeQuery, nodes []*Node, rootID string) (*TypedTreeNode[T], error) {
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	mapper, t, err := mapperForT[T](q)
	if err != nil {
		return nil, err
	}

	byID := make(map[string]*TypedTreeNode[T], len(nodes))

	for _, n := range nodes {
		v, err := mapper.fromNode(n)
		if err != nil {
			return nil, err
		}
		p, ok := v.(*T)
		if !ok {
			return nil, fmt.Errorf("nod: mapper returned %T, expected *%v", v, t)
		}
		byID[n.Core.Id] = &TypedTreeNode[T]{Node: p}
	}

	var root *TypedTreeNode[T]

	for _, n := range nodes {
		cur := byID[n.Core.Id]
		if n.Core.Id == rootID {
			root = cur
		}
		if n.Core.ParentId == nil || *n.Core.ParentId == "" {
			continue
		}
		parent := byID[*n.Core.ParentId]
		if parent != nil {
			parent.Children = append(parent.Children, cur)
		}
	}

	if root == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil
}

func buildAncestorTreeFromNodes[T any](q *NodeQuery, nodes []*Node) (*TypedTreeNode[T], error) {
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	mapper, t, err := mapperForT[T](q)
	if err != nil {
		return nil, err
	}

	byID := make(map[string]*TypedTreeNode[T], len(nodes))

	for _, n := range nodes {
		v, err := mapper.fromNode(n)
		if err != nil {
			return nil, err
		}
		p, ok := v.(*T)
		if !ok {
			return nil, fmt.Errorf("nod: mapper returned %T, expected *%v", v, t)
		}
		byID[n.Core.Id] = &TypedTreeNode[T]{Node: p}
	}

	var root *TypedTreeNode[T]

	for _, n := range nodes {
		cur := byID[n.Core.Id]
		if n.Core.ParentId == nil || *n.Core.ParentId == "" {
			root = cur
			continue
		}
		parent := byID[*n.Core.ParentId]
		if parent != nil {
			parent.Children = append(parent.Children, cur)
		}
	}

	if root == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil
}
