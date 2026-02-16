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

func fetchDescendantNodes(q *NodeQuery, rootID string) ([]*Node, error) {
	db := q.db.Model(&NodeCore{})

	sql := `
WITH RECURSIVE tree AS (
  SELECT * FROM node_cores WHERE id = ?
  UNION ALL
  SELECT n.* FROM node_cores n
  JOIN tree t ON n.parent_id = t.id
)
SELECT * FROM tree;
`
	var nodeCores []NodeCore
	if err := db.Raw(sql, rootID).Scan(&nodeCores).Error; err != nil {
		return nil, err
	}
	if len(nodeCores) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	nodes := make([]*Node, 0, len(nodeCores))
	for _, nc := range nodeCores {
		nodes = append(nodes, &Node{Core: nc})
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Tags = tagsByNode[n.Core.Id]
		}
	}

	if q.includeKV {
		ids := make([]string, 0, len(nodes))
		for _, n := range nodes {
			ids = append(ids, n.Core.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(ids)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.KV = kvsByNode[n.Core.Id]
		}
	}

	if q.includeContent {
		ids := make([]string, 0, len(nodes))
		for _, n := range nodes {
			ids = append(ids, n.Core.Id)
		}
		contentsByNode, err := (&ContentRepository{DB: q.db}).GetAllForNodes(ids)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Content = contentsByNode[n.Core.Id]
		}
	}

	return nodes, nil
}

func fetchAncestorNodes(q *NodeQuery, childID string) ([]*Node, error) {
	db := q.db.Model(&NodeCore{})

	sql := `
WITH RECURSIVE path AS (
  SELECT * FROM node_cores WHERE id = ?
  UNION ALL
  SELECT p.* FROM node_cores p
  JOIN path c ON p.id = c.parent_id
)
SELECT * FROM path;
`
	var nodeCores []NodeCore
	if err := db.Raw(sql, childID).Scan(&nodeCores).Error; err != nil {
		return nil, err
	}
	if len(nodeCores) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	nodes := make([]*Node, 0, len(nodeCores))
	for _, nc := range nodeCores {
		nodes = append(nodes, &Node{Core: nc})
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Tags = tagsByNode[n.Core.Id]
		}
	}

	if q.includeKV {
		ids := make([]string, 0, len(nodes))
		for _, n := range nodes {
			ids = append(ids, n.Core.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(ids)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.KV = kvsByNode[n.Core.Id]
		}
	}

	if q.includeContent {
		ids := make([]string, 0, len(nodes))
		for _, n := range nodes {
			ids = append(ids, n.Core.Id)
		}
		contentsByNode, err := (&ContentRepository{DB: q.db}).GetAllForNodes(ids)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Content = contentsByNode[n.Core.Id]
		}
	}

	return nodes, nil
}

func DescendantTreeAs[T any](q *NodeQuery, rootID string) (*TypedTreeNode[T], error) {
	nodes, err := fetchDescendantNodes(q, rootID)
	if err != nil {
		return nil, err
	}
	return buildTreeFromNodes[T](q, nodes, rootID)
}

func DescendantsAs[T any](q *NodeQuery, onlyRoots bool) ([]*TypedTreeNode[T], error) {
	roots, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	out := make([]*TypedTreeNode[T], 0)

	for _, n := range roots {
		isRoot := (n.Core.ParentId == nil || *n.Core.ParentId == "")
		if onlyRoots && !isRoot {
			continue
		}
		if !onlyRoots && !isRoot {
		}

		tree, err := DescendantTreeAs[T](q, n.Core.Id)
		if err != nil {
			return nil, err
		}
		out = append(out, tree)
	}

	return out, nil
}

func AncestorTreeAs[T any](q *NodeQuery, childID string) (*TypedTreeNode[T], error) {
	nodes, err := fetchAncestorNodes(q, childID)
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
