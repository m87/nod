package nod

import "gorm.io/gorm"


type Repository struct {
	db *gorm.DB
	Node *NodeRepository
}

func (r *Repository) Transaction(fc func(txRepo *Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{
			db:   tx,
			Node: &NodeRepository{DB: tx},
		}
		return fc(txRepo)
	})
}

func (r *Repository) Save(node *Node, tags []Tag) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(node).Error; err != nil {
			return err
		}

		if err := tx.Where("node_id = ?", node.Id).Delete(&NodeTag{}).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			if err := tx.FirstOrCreate(&tag, Tag{Name: tag.Name}).Error; err != nil {
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
