package nod

import "time"

// EdgeCore holds the core attributes of a directed edge between two nodes.
type EdgeCore struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	NamespaceId *string   `gorm:"type:varchar(36);index:idx_edge_namespace_id;index:idx_edge_namespace_source_kind,priority:1;index:idx_edge_namespace_target_kind,priority:1"`
	SourceId    string    `gorm:"type:varchar(36);not null;index:idx_edge_source_id;index:idx_edge_namespace_source_kind,priority:2;index:idx_edge_source_kind_name,priority:1"`
	Source      *NodeCore `gorm:"foreignKey:SourceId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TargetId    string    `gorm:"type:varchar(36);not null;index:idx_edge_target_id;index:idx_edge_namespace_target_kind,priority:2"`
	Target      *NodeCore `gorm:"foreignKey:TargetId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string    `gorm:"type:text;not null;index;default:'';index:idx_edge_source_kind_name,priority:3"`
	Kind        string    `gorm:"type:text;not null;index;default:'';index:idx_edge_namespace_source_kind,priority:3;index:idx_edge_namespace_target_kind,priority:3;index:idx_edge_source_kind_name,priority:2"`
	Status      string    `gorm:"type:text;not null;index;default:''"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

// Edge represents a directed edge with core data and key-value attributes.
type Edge struct {
	Core EdgeCore
	KV   map[string]*EdgeKV
}
