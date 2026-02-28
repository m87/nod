package nod

import (
	"time"

	"gorm.io/gorm"
)

type Blob struct {
	NodeId    string    `gorm:"type:char(36);primaryKey;index:idx_blob_node_id,priority:1"`
	Key       string    `gorm:"type:text;primaryKey;index:idx_blob_key,priority:2"`
	BlobId    string    `gorm:"type:char(36);not null"`
	Name      string    `gorm:"type:text;not null"`
	MimeType  string    `gorm:"type:text;not null"`
	Size      int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:datetime;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;autoUpdateTime"`
}

type BlobRepository struct {
	DB *gorm.DB
}

func (r *BlobRepository) GetAllForNodes(nodeIds []string) (map[string]map[string]*Blob, error) {
	var blobs []Blob
	result := make(map[string]map[string]*Blob)

	if err := r.DB.Find(&blobs, "node_id IN ?", nodeIds).Error; err != nil {
		return nil, err
	}

	for _, blob := range blobs {
		if _, exists := result[blob.NodeId]; !exists {
			result[blob.NodeId] = make(map[string]*Blob)
		}
		blobCopy := blob
		result[blob.NodeId][blob.Key] = &blobCopy
	}

	return result, nil
}

func (r *BlobRepository) Save(blob *Blob, data []byte) error {
	return r.DB.Save(blob).Error
}

func (r *BlobRepository) DeleteAll(nodeId string) error {
	return r.DB.Delete(&Blob{}, "node_id = ?", nodeId).Error
}
