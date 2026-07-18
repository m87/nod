package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func saveTagIfNotExists(tx *gorm.DB, namespaceId *string, name string) (*Tag, error) {
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

func (r *Repository) getEdgesTags(edgeIds []string) (map[string][]*Tag, error) {
	var edgeTags []EdgeTag
	err := r.db.Where("edge_id IN ?", edgeIds).Find(&edgeTags).Error
	if err != nil {
		return nil, err
	}

	tagIds := make([]string, 0, len(edgeTags))
	for _, edgeTag := range edgeTags {
		tagIds = append(tagIds, edgeTag.TagId)
	}

	var tags []*Tag
	err = r.db.Where("id IN ?", tagIds).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]*Tag)
	for _, tag := range tags {
		tagMap[tag.Id] = tag
	}

	result := make(map[string][]*Tag)
	for _, edgeTag := range edgeTags {
		if tag, exists := tagMap[edgeTag.TagId]; exists {
			result[edgeTag.EdgeId] = append(result[edgeTag.EdgeId], tag)
		}
	}
	return result, nil
}

func (r *Repository) getNodesTags(nodeIds []string) (map[string][]*Tag, error) {
	var nodeTags []NodeTag
	err := r.db.Where("node_id IN ?", nodeIds).Find(&nodeTags).Error
	if err != nil {
		return nil, err
	}

	tagIds := make([]string, 0, len(nodeTags))
	for _, nodeTag := range nodeTags {
		tagIds = append(tagIds, nodeTag.TagId)
	}

	var tags []*Tag
	err = r.db.Where("id IN ?", tagIds).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]*Tag)
	for _, tag := range tags {
		tagMap[tag.Id] = tag
	}

	result := make(map[string][]*Tag)
	for _, nodeTag := range nodeTags {
		if tag, exists := tagMap[nodeTag.TagId]; exists {
			result[nodeTag.NodeId] = append(result[nodeTag.NodeId], tag)
		}
	}
	return result, nil
}

func getEdgeTags(tx *gorm.DB, edgeId string) ([]*Tag, error) {
	var tags []*Tag
	err := tx.Joins("JOIN edge_tags ON edge_tags.tag_id = tags.id").
		Where("edge_tags.edge_id = ?", edgeId).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) getEdgeTags(edgeId string) ([]*Tag, error) {
	return getEdgeTags(r.db, edgeId)
}

func bindEdgeTagToEdge(tx *gorm.DB, edgeId string, tagId string) error {
	binding := &EdgeTag{
		EdgeId: edgeId,
		TagId:  tagId,
	}
	return tx.Save(binding).Error
}

func unbindEdgeTagsFromEdge(tx *gorm.DB, edgeId string) error {
	return tx.Where("edge_id = ?", edgeId).Delete(&EdgeTag{}).Error
}
