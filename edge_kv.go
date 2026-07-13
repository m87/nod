package nod

import (
	"time"

	"gorm.io/gorm"
)

// EdgeKV represents a key-value attribute attached to an edge.
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

type EdgeKVFilter struct {
	Key               string
	TextContains      *string
	NumberEquals      *float64
	IntEquals         *int
	Int64Equals       *int64
	BoolEquals        *bool
	TimeFrom          *time.Time
	TimeTo            *time.Time
	TextEquals        *string
	TextStartsWith    *string
	TextEndsWith      *string
	NumberGreaterThan *float64
	NumberLessThan    *float64
	IntGreaterThan    *int
	IntLessThan       *int
}

type EdgeKVRepository struct {
	DB *gorm.DB
}

func (r *EdgeKVRepository) Set(kv *EdgeKV) error {
	return r.DB.Save(kv).Error
}

func (r *EdgeKVRepository) Get(edgeId string, key string) (*EdgeKV, error) {
	var kv EdgeKV
	if err := r.DB.First(&kv, "edge_id = ? AND key = ?", edgeId, key).Error; err != nil {
		return nil, err
	}
	return &kv, nil
}

func (r *EdgeKVRepository) GetAll(edgeId string) (map[string]*EdgeKV, error) {
	var kvs []EdgeKV
	result := make(map[string]*EdgeKV)
	if err := r.DB.Where("edge_id = ?", edgeId).Find(&kvs).Error; err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		result[kv.Key] = &kv
	}
	return result, nil
}

func (r *EdgeKVRepository) DeleteAll(edgeId string) error {
	return r.DB.Where("edge_id = ?", edgeId).Delete(&EdgeKV{}).Error
}

func (r *EdgeKVRepository) Delete(edgeId string, key string) error {
	return r.DB.Where("edge_id = ? AND key = ?", edgeId, key).Delete(&EdgeKV{}).Error
}

func (r *EdgeKVRepository) Query(filters *EdgeKVFilter) ([]*EdgeKV, error) {
	var kvs []*EdgeKV
	query := r.DB.Model(&EdgeKV{})

	if filters.Key != "" {
		query = query.Where("key = ?", filters.Key)
	}
	if filters.TextContains != nil {
		query = query.Where("value_text LIKE ?", "%"+*filters.TextContains+"%")
	}
	if filters.NumberEquals != nil {
		query = query.Where("value_number = ?", *filters.NumberEquals)
	}
	if filters.IntEquals != nil {
		query = query.Where("value_int = ?", *filters.IntEquals)
	}
	if filters.Int64Equals != nil {
		query = query.Where("value_int64 = ?", *filters.Int64Equals)
	}
	if filters.BoolEquals != nil {
		query = query.Where("value_bool = ?", *filters.BoolEquals)
	}
	if filters.TimeFrom != nil {
		query = query.Where("value_time >= ?", *filters.TimeFrom)
	}
	if filters.TimeTo != nil {
		query = query.Where("value_time <= ?", *filters.TimeTo)
	}
	if filters.TextEquals != nil {
		query = query.Where("value_text = ?", *filters.TextEquals)
	}
	if filters.TextStartsWith != nil {
		query = query.Where("value_text LIKE ?", *filters.TextStartsWith+"%")
	}
	if filters.TextEndsWith != nil {
		query = query.Where("value_text LIKE ?", "%"+*filters.TextEndsWith)
	}
	if filters.NumberGreaterThan != nil {
		query = query.Where("value_number > ?", *filters.NumberGreaterThan)
	}
	if filters.NumberLessThan != nil {
		query = query.Where("value_number < ?", *filters.NumberLessThan)
	}
	if filters.IntGreaterThan != nil {
		query = query.Where("value_int > ?", *filters.IntGreaterThan)
	}
	if filters.IntLessThan != nil {
		query = query.Where("value_int < ?", *filters.IntLessThan)
	}

	if err := query.Find(&kvs).Error; err != nil {
		return nil, err
	}

	return kvs, nil
}
