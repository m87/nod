package nod

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


func Save[T any](r *Repository, model *T) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		t := reflect.TypeOf((*T)(nil)).Elem()
		mapper, err := r.Mappers.forType(t)
		if err != nil {
			return err
		}

		node, err := mapper.toNode(model)
		if err != nil {
			return err
		}

		if node.Core.Id == "" {
			node.Core.Id = uuid.New().String()
		}

		if err := tx.Save(&node.Core).Error; err != nil {
			return err
		}

		if err := tx.Model(&NodeTag{}).Where("node_id = ?", node.Core.Id).Delete(&NodeTag{}).Error; err != nil {
			return err
		}

		for _, tag := range node.Tags {
			if tag.Id == "" {
				tag.Id = uuid.New().String()
			}
			if err := tx.FirstOrCreate(tag, Tag{Id: tag.Id}).Error; err != nil {
				return err
			}
			nodeTag := &NodeTag{
				NodeId: node.Core.Id,
				TagId:  tag.Id,
			}
			if err := tx.Create(nodeTag).Error; err != nil {
				return err
			}
		}

		kvRepo := &KVRepository{DB: tx}
		if err := kvRepo.DeleteAll(node.Core.Id); err != nil {
			return err
		}
		for _, kv := range node.KV {
			kv.NodeId = node.Core.Id
			if err := kvRepo.Set(kv); err != nil {
				return err
			}
		}

		contentRepo := &ContentRepository{DB: tx}
		if err := contentRepo.DeleteAll(node.Core.Id); err != nil {
			return err
		}
		for _, content := range node.Content {
			content.NodeId = node.Core.Id
			if err := contentRepo.Save(content); err != nil {
				return err
			}
		}

		return nil
	})
}

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

	out := make([]*T, 0, len(nodes))
	for _, n := range nodes {
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
