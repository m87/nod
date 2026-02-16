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
	ValueBool   *bool      `gorm:"type:boolean"`
	ValueTime   *time.Time `gorm:"type:datetime"`
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
