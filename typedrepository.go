package nod

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TRepository[T any] struct {
	Db   *gorm.DB
	Node *NodeRepository
	Log  *slog.Logger
	Mapper NodeMapper[T]
}

func NewTRepository[T any](db *gorm.DB, log *slog.Logger, mapper NodeMapper[T]) *TRepository[T] {
	return &TRepository[T]{
		Db:   db,
		Node: &NodeRepository{DB: db},
		Log:  log,
		Mapper: mapper,
	}
}

func (r *TRepository[T]) Transaction(fc func(txRepo *TRepository[T]) error) error {
	r.Log.Debug(">> new transaction")
	return r.Db.Transaction(func(tx *gorm.DB) error {
		r.Log.Debug(">> new repository in transaction")
		txRepo := &TRepository[T]{
			Db:   tx,
			Node: &NodeRepository{DB: tx},
			Log:	r.Log,
			Mapper: r.Mapper,
		}
		r.Log.Debug(">> execute function in transaction")
		err := fc(txRepo)
		if err != nil {
			r.Log.Debug("<< rollback transaction due to error:", slog.String("error", err.Error()))
			return err
		}
		r.Log.Debug("<< end repository in transaction")
    return err
	})
}

func (r *TRepository[T]) Save(model *T) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		node, err := r.Mapper.ToNode(model)
		if err != nil {
			return err
		}
		if node.Core.Id == "" {
			node.Core.Id = uuid.New().String()
		}
		err = r.Db.Save(&node.Core).Error
		if err != nil {
			return err
		}

		if err := r.Db.Model(&NodeTag{}).Where("node_id = ?", node.Core.Id).Delete(&NodeTag{}).Error; err != nil {
			return err
		}
		for _, tag := range node.Tags {
			if tag.Id == "" {
				tag.Id = uuid.New().String()
			}
			if err := r.Db.FirstOrCreate(tag, Tag{Id: tag.Id}).Error; err != nil {
				return err
			}
			nodeTag := &NodeTag{
				NodeId: node.Core.Id,
				TagId:  tag.Id,
			}
			if err := r.Db.Create(nodeTag).Error; err != nil {
				return err
			}
		}

		kvRepo := &KVRepository{DB: r.Db}
		if err := kvRepo.DeleteAll(node.Core.Id); err != nil {
			return err
		}
		for _, kv := range node.KV {
			kv.NodeId = node.Core.Id
			if err := kvRepo.Set(kv); err != nil {
				return err
			}
		}

		contentRepo := &ContentRepository{DB: r.Db}
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

func (r *TRepository[T]) Delete(nodeId string) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		count := int64(0)
		db := tx.Model(&NodeCore{})
		if err := db.Where("parent_id = ?", nodeId).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("cannot delete node with children")
		}

		if err := tx.Delete(&NodeCore{}, "id = ?", nodeId).Error; err != nil {
			return err
		}
		if err := tx.Delete(&NodeTag{}, "node_id = ?", nodeId).Error; err != nil {
			return err
		}
		if err := tx.Delete(&KV{}, "node_id = ?", nodeId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TRepository[T]) Query() *TypedQuery[T] {
	return NewTypedQuery(r.Db, r.Log, r.Mapper)
}

