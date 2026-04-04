package nod

import (
	"time"

	"gorm.io/gorm"
)

// Tag represents a label that can be attached to one or more nodes.
type Tag struct {
	Id          string    `gorm:"type:char(36);primaryKey"`
	NamespaceId *string   `gorm:"type:char(36);index:idx_tags_namespace_id,priority:1;index"`
	Name        string    `gorm:"type:text;not null;index"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

// NodeTag represents the many-to-many relationship between nodes and tags.
type NodeTag struct {
	NodeId string    `gorm:"type:char(36);primaryKey;index:idx_node_tag,priority:1"`
	Node   *NodeCore `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:NodeId;references:Id"`
	TagId  string    `gorm:"type:char(36);primaryKey;index:idx_node_tag,priority:2"`
	Tag    *Tag      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TagId;references:Id"`
}

// TagRepository provides methods for managing tags in the database.
type TagRepository struct {
	DB *gorm.DB
}

func (r *TagRepository) Create(tag *Tag) error {
	return r.DB.Create(tag).Error
}

func (r *TagRepository) Save(tag *Tag) error {
	return r.DB.Save(tag).Error
}

func (r *TagRepository) GetAll(nodeId string) ([]Tag, error) {
	var tags []Tag
	err := r.DB.Joins("JOIN node_tags ON node_tags.tag_id = tags.id").
		Where("node_tags.node_id = ?", nodeId).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetAllForNodes(nodeIds []string) (map[string][]Tag, error) {
	var nodeTags []NodeTag
	result := make(map[string][]Tag)

	if err := r.DB.Find(&nodeTags, "node_id IN ?", nodeIds).Error; err != nil {
		return nil, err
	}

	tagIds := make([]string, 0, len(nodeTags))
	for _, nt := range nodeTags {
		tagIds = append(tagIds, nt.TagId)
	}

	var tags []Tag
	if err := r.DB.Find(&tags, "id IN ?", tagIds).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[string]Tag)
	for _, tag := range tags {
		tagMap[tag.Id] = tag
	}

	for _, nt := range nodeTags {
		if tag, exists := tagMap[nt.TagId]; exists {
			result[nt.NodeId] = append(result[nt.NodeId], tag)
		}
	}

	return result, nil
}

func (r *TagRepository) BindNodeTag(nodeId string, tagId string) error {
	nodeTag := NodeTag{
		NodeId: nodeId,
		TagId:  tagId,
	}
	return r.DB.FirstOrCreate(&nodeTag, nodeTag).Error
}

func (r *TagRepository) Delete(tagId string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tag_id = ?", tagId).Delete(&NodeTag{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Tag{}, "id = ?", tagId).Error
	})
}

// ConvertTagsToStringSlice extracts tag names from a slice of Tag pointers.
func ConvertTagsToStringSlice(tags []*Tag) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		result[i] = tag.Name
	}
	return result
}

// ConvertStringSliceToTags creates Tag pointers from a slice of tag names.
func ConvertStringSliceToTags(names []string) []*Tag {
	result := make([]*Tag, len(names))
	for i, name := range names {
		result[i] = &Tag{
			Name: name,
		}
	}
	return result
}
