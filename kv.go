package nod

import (
	"time"

	"gorm.io/gorm"
)

type KV struct {
	NodeId      string     `gorm:"type:char(36);primaryKey;index:idx_kv_node_id,priority:1"`
	Key         string     `gorm:"type:text;primaryKey;index:idx_kv_key,priority:2"`
	ValueText   *string    `gorm:"type:text"`
	ValueNumber *float64   `gorm:"type:real"`
	ValueInt    *int       `gorm:"type:integer"`
	ValueInt64  *int64     `gorm:"type:bigint"`
	ValueBool   *bool      `gorm:"type:boolean"`
	ValueTime   *time.Time `gorm:"type:datetime"`
}

type KVFilter struct {
	Key               *string
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

type KVRepository struct {
	DB *gorm.DB
}

func (r *KVRepository) Set(kv *KV) error {
	return r.DB.Save(kv).Error
}

func (r *KVRepository) Get(nodeId string, key string) (*KV, error) {
	var kv KV
	if err := r.DB.First(&kv, "node_id = ? AND key = ?", nodeId, key).Error; err != nil {
		return nil, err
	}
	return &kv, nil
}

func (r *KVRepository) GetAll(nodeId string) (map[string]*KV, error) {
	var kvs []KV
	result := make(map[string]*KV)

	if err := r.DB.Find(&kvs, "node_id = ?", nodeId).Error; err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		kvCopy := kv
		result[kv.Key] = &kvCopy
	}

	return result, nil
}

func (r *KVRepository) GetAllForNodes(nodeIds []string) (map[string]map[string]*KV, error) {
	var kvs []KV
	result := make(map[string]map[string]*KV)

	if err := r.DB.Find(&kvs, "node_id IN ?", nodeIds).Error; err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		kvCopy := kv
		if _, exists := result[kv.NodeId]; !exists {
			result[kv.NodeId] = make(map[string]*KV)
		}
		result[kv.NodeId][kv.Key] = &kvCopy
	}

	return result, nil
}

func (r *KVRepository) DeleteAll(nodeId string) error {
	return r.DB.Delete(&KV{}, "node_id = ?", nodeId).Error
}

func (r *KVRepository) Delete(nodeId string, key string) error {
	return r.DB.Delete(&KV{}, "node_id = ? AND key = ?", nodeId, key).Error
}

func (r *KVRepository) Query(filters []*KVFilter) ([]*KV, error) {
	var kvs []KV
	db := r.DB.Model(&KV{})

	for _, filter := range filters {
		if filter.Key != nil {
			db = db.Where("key = ?", *filter.Key)
		}
		if filter.TextContains != nil {
			db = db.Where("value_text LIKE ?", "%"+*filter.TextContains+"%")
		}
		if filter.TextEquals != nil {
			db = db.Where("value_text = ?", *filter.TextEquals)
		}
		if filter.TextStartsWith != nil {
			db = db.Where("value_text LIKE ?", *filter.TextStartsWith+"%")
		}
		if filter.TextEndsWith != nil {
			db = db.Where("value_text LIKE ?", "%"+*filter.TextEndsWith)
		}
		if filter.NumberEquals != nil {
			db = db.Where("value_number = ?", *filter.NumberEquals)
		}
		if filter.IntEquals != nil {
			db = db.Where("value_int = ?", *filter.IntEquals)
		}
		if filter.BoolEquals != nil {
			db = db.Where("value_bool = ?", *filter.BoolEquals)
		}
		if filter.TimeFrom != nil {
			db = db.Where("value_time >= ?", *filter.TimeFrom)
		}
		if filter.TimeTo != nil {
			db = db.Where("value_time <= ?", *filter.TimeTo)
		}
		if filter.NumberGreaterThan != nil {
			db = db.Where("value_number > ?", *filter.NumberGreaterThan)
		}
		if filter.NumberLessThan != nil {
			db = db.Where("value_number < ?", *filter.NumberLessThan)
		}
		if filter.IntGreaterThan != nil {
			db = db.Where("value_int > ?", *filter.IntGreaterThan)
		}
		if filter.IntLessThan != nil {
			db = db.Where("value_int < ?", *filter.IntLessThan)
		}
	}

	if err := db.Find(&kvs).Error; err != nil {
		return nil, err
	}

	result := make([]*KV, len(kvs))
	for i, kv := range kvs {
		kvCopy := kv
		result[i] = &kvCopy
	}

	return result, nil
}

func ConvertKVToStringMap(kvs map[string]*KV) map[string]string {
	result := make(map[string]string)
	for key, kv := range kvs {
		if kv.ValueText != nil {
			result[key] = *kv.ValueText
		}
	}
	return result
}

func ConvertStringMapToKV(data map[string]string) map[string]*KV {
	result := make(map[string]*KV)
	for key, value := range data {
		result[key] = &KV{
			Key:       key,
			ValueText: &value,
		}
	}
	return result
}

func ConvertKVToIntMap(kvs map[string]*KV) map[string]int {
	result := make(map[string]int)
	for key, kv := range kvs {
		if kv.ValueInt != nil {
			result[key] = *kv.ValueInt
		}
	}
	return result
}

func ConvertInt64MapToKV(data map[string]int64) map[string]*KV {
	result := make(map[string]*KV)
	for key, value := range data {
		result[key] = &KV{
			Key:        key,
			ValueInt64: &value,
		}
	}
	return result
}

func ConvertKVToInt64Map(kvs map[string]*KV) map[string]int64 {
	result := make(map[string]int64)
	for key, kv := range kvs {
		if kv.ValueInt64 != nil {
			result[key] = *kv.ValueInt64
		}
	}
	return result
}

func ConvertIntMapToKV(data map[string]int) map[string]*KV {
	result := make(map[string]*KV)
	for key, value := range data {
		result[key] = &KV{
			Key:      key,
			ValueInt: &value,
		}
	}
	return result
}

func SafeString(data map[string]*KV, key string) string {
	if kv, exists := data[key]; exists && kv.ValueText != nil {
		return *kv.ValueText
	}
	return ""
}

func SafeInt(data map[string]*KV, key string) int {
	if kv, exists := data[key]; exists && kv.ValueInt != nil {
		return *kv.ValueInt
	}
	return 0
}

func SafeInt64(data map[string]*KV, key string) int64 {
	if kv, exists := data[key]; exists && kv.ValueInt64 != nil {
		return *kv.ValueInt64
	}
	return 0
}

func SafeBool(data map[string]*KV, key string) bool {
	if kv, exists := data[key]; exists && kv.ValueBool != nil {
		return *kv.ValueBool
	}
	return false
}

func SafeTime(data map[string]*KV, key string) time.Time {
	if kv, exists := data[key]; exists && kv.ValueTime != nil {
		return *kv.ValueTime
	}
	return time.Time{}
}

func SafeFloat64(data map[string]*KV, key string) float64 {
	if kv, exists := data[key]; exists && kv.ValueNumber != nil {
		return *kv.ValueNumber
	}
	return 0
}

func SafeStringWithDefault(data map[string]*KV, key string, defaultValue string) string {
	if kv, exists := data[key]; exists && kv.ValueText != nil {
		return *kv.ValueText
	}
	return defaultValue
}

func SafeIntWithDefault(data map[string]*KV, key string, defaultValue int) int {
	if kv, exists := data[key]; exists && kv.ValueInt != nil {
		return *kv.ValueInt
	}
	return defaultValue
}

func SafeInt64WithDefault(data map[string]*KV, key string, defaultValue int64) int64 {
	if kv, exists := data[key]; exists && kv.ValueInt64 != nil {
		return *kv.ValueInt64
	}
	return defaultValue
}

func SafeBoolWithDefault(data map[string]*KV, key string, defaultValue bool) bool {
	if kv, exists := data[key]; exists && kv.ValueBool != nil {
		return *kv.ValueBool
	}
	return defaultValue
}

func SafeTimeWithDefault(data map[string]*KV, key string, defaultValue time.Time) time.Time {
	if kv, exists := data[key]; exists && kv.ValueTime != nil {
		return *kv.ValueTime
	}
	return defaultValue
}

func SafeFloat64WithDefault(data map[string]*KV, key string, defaultValue float64) float64 {
	if kv, exists := data[key]; exists && kv.ValueNumber != nil {
		return *kv.ValueNumber
	}
	return defaultValue
}
