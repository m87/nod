package nod

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Save persists a model of type T by converting it to a Node using the registered mapper.
func Save[T any](r *Repository, model *T) (string, error) {
	var nodeId string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		t := reflect.TypeOf((*T)(nil)).Elem()
		mapper, err := r.mappers.forType(t)
		if err != nil {
			return err
		}

		node, err := mapper.toNode(model)
		if err != nil {
			return err
		}

		nodeId = ensureNodeID(node)
		return saveNodeGraph(tx, node)
	})
	if err != nil {
		return "", err
	}
	return nodeId, nil
}

// ListAs fetches nodes matching the query and converts them to type T using the registered mapper.
func ListAs[T any](q *NodeQuery) ([]*T, error) {
	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf((*T)(nil)).Elem()
	mapper, err := q.mappers.forType(t)
	if err != nil {
		return nil, err
	}

	out := []*T{}
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
		out = append(out, p)
	}
	return out, nil
}

// FirstAs returns the first node matching the query converted to type T, or ErrRecordNotFound.
func FirstAs[T any](q *NodeQuery) (*T, error) {
	items, err := ListAs[T](q.Limit(1))
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return items[0], nil
}
