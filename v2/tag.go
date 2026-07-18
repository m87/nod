package nod

import "time"

// Tag represents a label that can be attached to one or more nodes.
type Tag struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	NamespaceId *string   `gorm:"type:varchar(36);index:idx_tags_namespace_id,priority:1;index;uniqueIndex:idx_tags_namespace_id_name,priority:1"`
	Name        string    `gorm:"type:text;not null;index;uniqueIndex:idx_tags_namespace_id_name,priority:2"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

// NodeTag represents the many-to-many relationship between nodes and tags.
type NodeTag struct {
	NodeId string    `gorm:"type:varchar(36);primaryKey;index:idx_node_tag,priority:1"`
	Node   *NodeCore `gorm:"foreignKey:NodeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TagId  string    `gorm:"type:varchar(36);primaryKey;index:idx_node_tag,priority:2"`
	Tag    *Tag      `gorm:"foreignKey:TagId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// EdgeTag represents the many-to-many relationship between edges and tags.
type EdgeTag struct {
	EdgeId string    `gorm:"type:varchar(36);not null;primaryKey;index:idx_edge_tag,priority:1"`
	Edge   *EdgeCore `gorm:"foreignKey:EdgeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TagId  string    `gorm:"type:varchar(36);not null;primaryKey;index:idx_edge_tag,priority:2"`
	Tag    *Tag      `gorm:"foreignKey:TagId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
