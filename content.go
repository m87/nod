package nod

import (
	"time"

	"gorm.io/gorm"
)

type Content struct {
	NodeId    string    `gorm:"type:char(36);primaryKey;index:idx_kv_node_id,priority:1"`
	Key       string    `gorm:"type:text;primaryKey;index:idx_kv_key,priority:2"`
	Value     *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:datetime;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;autoUpdateTime"`
}

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

