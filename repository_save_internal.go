package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ensureNodeID(node *Node) string {
	if node.Core.Id == "" {
		node.Core.Id = uuid.New().String()
	}
	return node.Core.Id
}

func saveNodeGraph(tx *gorm.DB, node *Node) error {
	if err := tx.Save(&node.Core).Error; err != nil {
		return err
	}

	if err := tx.Model(&NodeTag{}).Where("node_id = ?", node.Core.Id).Delete(&NodeTag{}).Error; err != nil {
		return err
	}

	for _, tag := range node.Tags {
		if tag.Id == "" {
			tag.Id = uuid.New().String()
		}
		if err := tx.FirstOrCreate(tag, Tag{Id: tag.Id}).Error; err != nil {
			return err
		}
		nodeTag := &NodeTag{
			NodeId: node.Core.Id,
			TagId:  tag.Id,
		}
		if err := tx.Create(nodeTag).Error; err != nil {
			return err
		}
	}

	kvRepo := &KVRepository{DB: tx}
	if err := kvRepo.DeleteAll(node.Core.Id); err != nil {
		return err
	}
	for _, kv := range node.KV {
		kv.NodeId = node.Core.Id
		if err := kvRepo.Set(kv); err != nil {
			return err
		}
	}

	contentRepo := &ContentRepository{DB: tx}
	if err := contentRepo.DeleteAll(node.Core.Id); err != nil {
		return err
	}
	for _, content := range node.Content {
		content.NodeId = node.Core.Id
		if err := contentRepo.Save(content); err != nil {
			return err
		}
	}

	return nil
}
