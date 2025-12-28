package nod

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type NodeCore struct {
	Id          string         `gorm:"type:char(36);primaryKey"`
	NamespaceId *string        `gorm:"type:char(36);index:idx_namespace_id,priority:1;index"`
	ParentId    *string        `gorm:"type:char(36);index:idx_parent_id,priority:2;index"`
	Type        string         `gorm:"type:text;not null;index;default:''"`
	Kind        string         `gorm:"type:text;not null;index;default:''"`
	Status      string         `gorm:"type:text;not null;index;default:''"`
	Name        string         `gorm:"type:text;not null;index"`
	Metadata    datatypes.JSON `gorm:"type:json;not null;default:'{}'"`
	CreatedAt   time.Time      `gorm:"type:datetime;not null;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null;autoUpdateTime"`
}

type NodeRepository struct {
	DB *gorm.DB
}

func (r *NodeRepository) Create(node *Node) error {
	return r.DB.Create(node).Error
}

func (r *NodeRepository) Save(node *Node) error {
	return r.DB.Save(node).Error
}

func (r *NodeRepository) Delete(nodeId string) error {
	return r.DB.Delete(&Node{}, "id = ?", nodeId).Error
}

