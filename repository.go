package nod

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Node struct {
	Core    NodeCore
	Tags    []*Tag
	KV      map[string]*KV
	Content map[string]*Content
}

type Repository struct {
	Db      *gorm.DB
	Log     *slog.Logger
	Mappers *MapperRegistry
}

func NewRepository(db *gorm.DB, log *slog.Logger, mappers *MapperRegistry) *Repository {
	return &Repository{
		Db:      db,
		Log:     log,
		Mappers: mappers,
	}
}

func (r *Repository) Transaction(fc func(txRepo *Repository) error) error {
	r.Log.Debug(">> new transaction")
	return r.Db.Transaction(func(tx *gorm.DB) error {
		r.Log.Debug(">> new repository in transaction")
		txRepo := &Repository{
			Db:      tx,
			Log:     r.Log,
			Mappers: r.Mappers,
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

func (r *Repository) Delete(nodeId string) error {
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

func (r *Repository) Query() *NodeQuery {
	return NewNodeQuery(r.Db, r.Log, r.Mappers)
}

func loadTagsByNode(db *gorm.DB, nodes []*Node) (map[string][]*Tag, error) {
	ids := make([]string, 0, len(nodes))
	for _, n := range nodes {
		ids = append(ids, n.Core.Id)
	}

	type row struct {
		NodeID string
		TagID  string
		Name   string
	}
	var rows []row
	if err := db.Table("node_tags nt").
		Select("nt.node_id as node_id, t.id as tag_id, t.name as name").
		Joins("JOIN tags t ON t.id = nt.tag_id").
		Where("nt.node_id IN ?", ids).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make(map[string][]*Tag, len(ids))
	for _, r := range rows {
		out[r.NodeID] = append(out[r.NodeID], &Tag{Id: r.TagID, Name: r.Name})
	}
	return out, nil
}

func (r *Repository) Save(node *Node) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {

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

		return nil
	})
}
