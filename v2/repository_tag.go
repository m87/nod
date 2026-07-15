package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func saveNodeTagIfNotExists(tx *gorm.DB, namespaceId *string, name string) (*Tag, error) {
	var tag Tag
	err := tx.Where("namespace_id = ? AND name = ?", namespaceId, name).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tag = Tag{
				Id:          uuid.NewString(),
				NamespaceId: namespaceId,
				Name:        name,
			}
			if err := tx.Create(&tag).Error; err != nil {
				return nil, err
			}
			return &tag, nil
		}
		return nil, err
	}
	return &tag, nil
}

func getNodeTags(tx *gorm.DB, nodeId string) ([]*Tag, error) {
	var tags []*Tag
	err := tx.Joins("JOIN node_tags ON node_tags.tag_id = tags.id").
		Where("node_tags.node_id = ?", nodeId).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) getNodeTags(nodeId string) ([]*Tag, error) {
	return getNodeTags(r.db, nodeId)
}

func bindNodeTagToNode(tx *gorm.DB, nodeId string, tagId string) error {
	binding := &NodeTag{
		NodeId: nodeId,
		TagId:  tagId,
	}
	return tx.Save(binding).Error
}

func unbindNodeTagsFromNode(tx *gorm.DB, nodeId string) error {
	return tx.Where("node_id = ?", nodeId).Delete(&NodeTag{}).Error
}
