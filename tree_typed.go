package nod

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// TypedTreeNode represents a typed node in a tree structure with its children.
type TypedTreeNode[T any] struct {
	Node     *T
	Children []*TypedTreeNode[T]
}

// AncestorTreeAs builds an ancestor tree for the given child ID, mapping nodes to type T.
func AncestorTreeAs[T any](q *NodeQuery, childID string) (*TypedTreeNode[T], error) {
	nodes, err := q.fetchAncestorNodes(childID)
	if err != nil {
		return nil, err
	}
	return buildAncestorTreeFromNodes[T](q, nodes)
}

// AncestorsAs returns ancestor trees for all nodes matching the query, mapped to type T.
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

// DescendantTreeAs builds a descendant tree rooted at the given ID, mapping nodes to type T.
func DescendantTreeAs[T any](q *NodeQuery, rootID string) (*TypedTreeNode[T], error) {
	nodes, err := q.fetchDescendantNodes(rootID)
	if err != nil {
		return nil, err
	}
	return buildTreeFromNodes[T](q, nodes, rootID)
}

// DescendantsAs returns descendant trees for nodes matching the query, mapped to type T.
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

// mapNodesToTyped maps a slice of nodes to typed tree nodes using the registered mapper.
func mapNodesToTyped[T any](q *NodeQuery, nodes []*Node) (map[string]*TypedTreeNode[T], error) {
	mapper, t, err := mapperForT[T](q)
	if err != nil {
		return nil, err
	}

	byID := map[string]*TypedTreeNode[T]{}
	for _, n := range nodes {
		if !mapper.isApplicable(n) {
			continue
		}
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
	return byID, nil
}

// linkTypedTree links typed tree nodes into a tree structure and returns the root.
func linkTypedTree[T any](nodes []*Node, byID map[string]*TypedTreeNode[T], isRoot func(*Node) bool) (*TypedTreeNode[T], error) {
	var root *TypedTreeNode[T]
	for _, n := range nodes {
		cur := byID[n.Core.Id]
		if cur == nil {
			continue
		}
		if isRoot(n) {
			root = cur
		}
		if n.Core.ParentId != nil && *n.Core.ParentId != "" {
			if parent := byID[*n.Core.ParentId]; parent != nil {
				parent.Children = append(parent.Children, cur)
			}
		}
	}
	if root == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil
}

func buildTreeFromNodes[T any](q *NodeQuery, nodes []*Node, rootID string) (*TypedTreeNode[T], error) {
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	byID, err := mapNodesToTyped[T](q, nodes)
	if err != nil {
		return nil, err
	}
	return linkTypedTree(nodes, byID, func(n *Node) bool {
		return n.Core.Id == rootID
	})
}

func buildAncestorTreeFromNodes[T any](q *NodeQuery, nodes []*Node) (*TypedTreeNode[T], error) {
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	byID, err := mapNodesToTyped[T](q, nodes)
	if err != nil {
		return nil, err
	}
	return linkTypedTree(nodes, byID, func(n *Node) bool {
		return n.Core.ParentId == nil || *n.Core.ParentId == ""
	})
}
