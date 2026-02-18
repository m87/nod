package nod

import (
	"time"
)

type NodeCore struct {
	Id          string    `gorm:"type:char(36);primaryKey"`
	NamespaceId *string   `gorm:"type:char(36);index:idx_namespace_id,priority:1;index"`
	ParentId    *string   `gorm:"type:char(36);index:idx_parent_id,priority:2;index"`
	Kind				string    `gorm:"type:text;not null;index;default:''"`
	Status      string    `gorm:"type:text;not null;index;default:''"`
	Name        string    `gorm:"type:text;not null;index"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;autoUpdateTime"`
}
