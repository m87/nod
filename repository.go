package nod

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

// Node represents a tree node with core data, tags, key-value attributes, and content.
type Node struct {
	Core    NodeCore
	Tags    []*Tag
	KV      map[string]*KV
	Content map[string]*Content
}

// Repository provides access to nodes and related data in the database.
type Repository struct {
	db      *gorm.DB
	log     *slog.Logger
	mappers *MapperRegistry
}

// NewRepository creates a new Repository instance.
func NewRepository(db *gorm.DB, log *slog.Logger, mappers *MapperRegistry) *Repository {
	return &Repository{
		db:      db,
		log:     log,
		mappers: mappers,
	}
}

// DB returns the underlying GORM database connection.
func (r *Repository) DB() *gorm.DB { return r.db }

// Log returns the repository's logger.
func (r *Repository) Log() *slog.Logger { return r.log }

// Mappers returns the mapper registry used by this repository.
func (r *Repository) Mappers() *MapperRegistry { return r.mappers }

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Transaction executes a function within a database transaction.
func (r *Repository) Transaction(fc func(txRepo *Repository) error) error {
	r.log.Debug(">> new transaction")
	return r.db.Transaction(func(tx *gorm.DB) error {
		r.log.Debug(">> new repository in transaction")
		txRepo := &Repository{
			db:      tx,
			log:     r.log,
			mappers: r.mappers,
		}
		r.log.Debug(">> execute function in transaction")
		err := fc(txRepo)
		if err != nil {
			r.log.Debug("<< rollback transaction due to error:", slog.String("error", err.Error()))
			return err
		}
		r.log.Debug("<< end repository in transaction")
		return err
	})
}

func (r *Repository) Delete(nodeId string) error {
	if nodeId == "" {
		return fmt.Errorf("nod: nodeId must not be empty")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
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
	return NewNodeQuery(r.db, r.log, r.mappers)
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

func (r *Repository) Save(node *Node) (string, error) {
	nodeID := ensureNodeID(node)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return saveNodeGraph(tx, node)
	})
	if err != nil {
		return "", err
	}
	return nodeID, nil
}
