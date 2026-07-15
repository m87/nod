package nod

import "time"

// NodeContent represents a named text content block attached to a node. Should be used for storing large text data associated with a node.
type NodeContent struct {
	NodeId    string    `gorm:"type:varchar(36);primaryKey;index:idx_content_node_id,priority:1"`
	Node      *NodeCore `gorm:"foreignKey:NodeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Key       string    `gorm:"type:text;primaryKey;index:idx_content_key,priority:2"`
	Value     *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}


