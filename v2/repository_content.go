package nod

import "gorm.io/gorm"



func (r *Repository) getNodeContent(nodeId string, key string) (*NodeContent, error) {
	var content NodeContent
	err := r.db.Where("node_id = ? AND key = ?", nodeId, key).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *Repository) getNodeContents(nodeId string) ([]*NodeContent, error) {
	var contents []*NodeContent
	err := r.db.Where("node_id = ?", nodeId).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (r *Repository) saveNodeContent(content *NodeContent) error {
	return saveNodeContent(r.db, content)
}

func (r *Repository) deleteNodeContent(nodeId string, key string) error {
	return deleteNodeContent(r.db, nodeId, key)
}

func saveNodeContent(tx *gorm.DB, content *NodeContent) error {
	if content == nil {
		return NewNodeContentIsNilError()
	}
	return tx.Save(content).Error
}

func deleteNodeContent(tx *gorm.DB, nodeId string, key string) error {
	return tx.Where("node_id = ? AND key = ?", nodeId, key).Delete(&NodeContent{}).Error
}

func saveNodeContents(tx *gorm.DB, contents []*NodeContent) error {
	for _, content := range contents {
		if err := saveNodeContent(tx, content); err != nil {
			return err
		}
	}
	return nil
}

func deleteNodeContents(tx *gorm.DB, nodeId string) error {
	return tx.Where("node_id = ?", nodeId).Delete(&NodeContent{}).Error
}

func (r *Repository) saveNodeContents(contents []*NodeContent) error {
	return saveNodeContents(r.db, contents)
}

func (r *Repository) deleteNodeContents(nodeId string) error {
	return deleteNodeContents(r.db, nodeId)
}

