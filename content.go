package nod

import (
	"time"

	"gorm.io/gorm"
)

// Content represents a named text content block attached to a node.
type Content struct {
	NodeId    string    `gorm:"type:char(36);primaryKey;index:idx_content_node_id,priority:1"`
	Key       string    `gorm:"type:text;primaryKey;index:idx_content_key,priority:2"`
	Value     *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// ContentRepository provides methods for managing content blocks in the database.
type ContentRepository struct {
	DB *gorm.DB
}

func (r *ContentRepository) GetAllForNodes(nodeIds []string) (map[string]map[string]*Content, error) {
	var contents []Content
	result := make(map[string]map[string]*Content)

	if err := r.DB.Find(&contents, "node_id IN ?", nodeIds).Error; err != nil {
		return nil, err
	}

	for _, content := range contents {
		if _, exists := result[content.NodeId]; !exists {
			result[content.NodeId] = make(map[string]*Content)
		}
		contentCopy := content
		result[content.NodeId][content.Key] = &contentCopy
	}

	return result, nil
}

func (r *ContentRepository) Save(content *Content) error {
	return r.DB.Save(content).Error
}

func (r *ContentRepository) DeleteAll(nodeId string) error {
	return r.DB.Delete(&Content{}, "node_id = ?", nodeId).Error
}

// ConvertContentToStringMap extracts string values from a Content map.
func ConvertContentToStringMap(contents map[string]*Content) map[string]string {
	result := make(map[string]string)
	for key, content := range contents {
		if content.Value != nil {
			result[key] = *content.Value
		}
	}
	return result
}

// ConvertStringMapToContent creates a Content map from string key-value pairs.
func ConvertStringMapToContent(data map[string]string) map[string]*Content {
	result := make(map[string]*Content)
	for key, value := range data {
		result[key] = &Content{
			Key:   key,
			Value: &value,
		}
	}
	return result
}
