package nod

import "time"

// NodeKV represents a key-value attribute attached to a node.
type NodeKV struct {
	NodeId      string    `gorm:"type:varchar(36);primaryKey;index:idx_kv_node_id,priority:1"`
	Node        *NodeCore `gorm:"foreignKey:NodeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Key         string    `gorm:"type:text;primaryKey;index:idx_kv_key,priority:2"`
	ValueText   *string   `gorm:"type:text"`
	ValueNumber *float64
	ValueInt    *int   `gorm:"type:integer"`
	ValueInt64  *int64 `gorm:"type:bigint"`
	ValueBool   *bool  `gorm:"type:boolean"`
	ValueTime   *time.Time
}

type EdgeKV struct {
	EdgeId      string    `gorm:"type:varchar(36);not null;primaryKey;index:idx_edge_kv_edge_id,priority:1"`
	Edge        *EdgeCore `gorm:"foreignKey:EdgeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Key         string    `gorm:"type:text;not null;primaryKey;index:idx_edge_kv_key,priority:2"`
	ValueText   *string   `gorm:"type:text"`
	ValueNumber *float64
	ValueInt    *int   `gorm:"type:integer"`
	ValueInt64  *int64 `gorm:"type:bigint"`
	ValueBool   *bool  `gorm:"type:boolean"`
	ValueTime   *time.Time
}
