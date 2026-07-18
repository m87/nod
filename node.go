package nod

import "time"

// Node represents a node with core data, tags, key-value attributes, and content.
type Node struct {
	Core    NodeCore
	Tags    []*Tag
	KV      map[string]*NodeKV
	Content map[string]*NodeContent
}

// NodeCore holds the core attributes of a node stored in the database.
type NodeCore struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	NamespaceId *string   `gorm:"type:varchar(36);index:idx_namespace_id,priority:1;index"`
	ParentId    *string   `gorm:"type:varchar(36);index:idx_parent_id,priority:2;index"`
	Parent      *NodeCore `gorm:"foreignKey:ParentId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Kind        string    `gorm:"type:text;not null;index;default:''"`
	Status      string    `gorm:"type:text;not null;index;default:''"`
	Name        string    `gorm:"type:text;not null;index"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}
