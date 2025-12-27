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
