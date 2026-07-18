package nod

import "gorm.io/gorm"

func (r *Repository) getNodeContents(nodeId string) ([]*NodeContent, error) {
	var contents []*NodeContent
	err := r.db.Where("node_id = ?", nodeId).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func saveNodeContent(tx *gorm.DB, content *NodeContent) error {
	if content == nil {
		return NewNodeContentIsNilError()
	}
	return tx.Save(content).Error
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

func (r *Repository) getEdgesContents(edgeIds []string) (map[string][]*EdgeContent, error) {
	var contents []*EdgeContent
	err := r.db.Where("edge_id IN ?", edgeIds).Find(&contents).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*EdgeContent)
	for _, content := range contents {
		result[content.EdgeId] = append(result[content.EdgeId], content)
	}
	return result, nil
}

func (r *Repository) getNodesContents(nodeIds []string) (map[string][]*NodeContent, error) {
	var contents []*NodeContent
	err := r.db.Where("node_id IN ?", nodeIds).Find(&contents).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*NodeContent)
	for _, content := range contents {
		result[content.NodeId] = append(result[content.NodeId], content)
	}
	return result, nil
}

func (r *Repository) getEdgeContents(edgeId string) ([]*EdgeContent, error) {
	var contents []*EdgeContent
	err := r.db.Where("edge_id = ?", edgeId).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func saveEdgeContent(tx *gorm.DB, content *EdgeContent) error {
	if content == nil {
		return NewEdgeContentIsNilError()
	}
	return tx.Save(content).Error
}

func saveEdgeContents(tx *gorm.DB, contents []*EdgeContent) error {
	for _, content := range contents {
		if err := saveEdgeContent(tx, content); err != nil {
			return err
		}
	}
	return nil
}

func deleteEdgeContents(tx *gorm.DB, edgeId string) error {
	return tx.Where("edge_id = ?", edgeId).Delete(&EdgeContent{}).Error
}
