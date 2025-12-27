package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Repository struct {
	Db *gorm.DB
	Node *NodeRepository
}

func (r *Repository) Transaction(fc func(txRepo *Repository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{
			Db:   tx,
			Node: &NodeRepository{DB: tx},
		}
		return fc(txRepo)
	})
}

func (r *Repository) Save(node *Node, tags []Tag) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(node).Error; err != nil {
			return err
		}

		if err := tx.Where("node_id = ?", node.Id).Delete(&NodeTag{}).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			if err := tx.FirstOrCreate(&tag, Tag{Id: uuid.New().String(), Name: tag.Name}).Error; err != nil {
				return err
			}

			nodeTag := NodeTag{
				NodeId: node.Id,
				TagId:  tag.Id,
			}
			if err := tx.Create(&nodeTag).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
