package nod

import "time"

// Tag represents a label that can be attached to one or more nodes.
type Tag struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	NamespaceId *string   `gorm:"type:varchar(36);index:idx_tags_namespace_id,priority:1;index"`
	Name        string    `gorm:"type:text;not null;index"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

// NodeTag represents the many-to-many relationship between nodes and tags.
type NodeTag struct {
	NodeId string    `gorm:"type:varchar(36);primaryKey;index:idx_node_tag,priority:1"`
	Node   *NodeCore `gorm:"foreignKey:NodeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TagId  string    `gorm:"type:varchar(36);primaryKey;index:idx_node_tag,priority:2"`
	Tag    *Tag      `gorm:"foreignKey:TagId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
