package nod

import "time"


type Tag struct {
	Id string `gorm:"type:char(36);primaryKey"`
	NamespaceId *string `gorm:"type:char(36);index:idx_namespace_id,priority:1;index"`
	Name string `gorm:"type:text;not null;index"`
	CreatedAt time.Time `gorm:"type:datetime;not null;autoCreateTime"`
}

type NodeTag struct {
	NodeId string `gorm:"type:char(36);primaryKey;index:idx_node_tag,priority:1"`
	TagId string `gorm:"type:char(36);primaryKey;index:idx_node_tag,priority:2"`
}
