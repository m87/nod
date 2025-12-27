package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Repository struct {
	Db *gorm.DB
	Node *NodeRepository
}

type TreeNode struct {
	Node     Node
	Tags		 []Tag
	Children []TreeNode
}

func (r *Repository) Transaction(fc func(txRepo *Repository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{
			Db:   tx,
			Node: &NodeRepository{DB: tx},
		}
		return fc(txRepo)
	})
}

func (r *Repository) Save(node *Node, tags []Tag) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(node).Error; err != nil {
			return err
		}

		if err := tx.Where("node_id = ?", node.Id).Delete(&NodeTag{}).Error; err != nil {
			return err
		}

		for _, tag := range tags {
			if err := tx.FirstOrCreate(&tag, Tag{Id: uuid.New().String(), Name: tag.Name}).Error; err != nil {
				return err
			}

			nodeTag := NodeTag{
				NodeId: node.Id,
				TagId:  tag.Id,
			}
			if err := tx.Create(&nodeTag).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) LoadTree(rootID string) (TreeNode, error) {
	var nodes []Node

	sql := `
WITH RECURSIVE tree AS (
  SELECT * FROM nodes WHERE id = ?
  UNION ALL
  SELECT n.* FROM nodes n
  JOIN tree t ON n.parent_id = t.id
)
SELECT * FROM tree;
`
	if err := r.Db.Raw(sql, rootID).Scan(&nodes).Error; err != nil {
		return TreeNode{}, err
	}
	if len(nodes) == 0 {
		return TreeNode{}, gorm.ErrRecordNotFound
	}

	nodeIDs := make([]string, 0, len(nodes))
	for _, n := range nodes {
		nodeIDs = append(nodeIDs, n.Id)
	}

	type row struct {
		NodeID string
		TagID  string
		Name   string
	}
	var rows []row
	if err := r.Db.Table("node_tags nt").
		Select("nt.node_id as node_id, t.id as tag_id, t.name as name").
		Joins("JOIN tags t ON t.id = nt.tag_id").
		Where("nt.node_id IN ?", nodeIDs).
		Scan(&rows).Error; err != nil {
		return TreeNode{}, err
	}

	tagsByNode := make(map[string][]Tag, len(nodeIDs))
	for _, r := range rows {
		tagsByNode[r.NodeID] = append(tagsByNode[r.NodeID], Tag{
			Id:   r.TagID,
			Name: r.Name,
		})
	}

	byID := make(map[string]*TreeNode, len(nodes))
	for _, n := range nodes {
		tn := &TreeNode{
			Node: n,
			Tags: tagsByNode[n.Id],
		}
		byID[n.Id] = tn
	}

	var root *TreeNode
	for _, n := range nodes {
		cur := byID[n.Id]
		if n.Id == rootID {
			root = cur
			continue
		}
		if n.ParentId == nil {
			continue
		}
		parent := byID[*n.ParentId]
		if parent == nil {
			continue
		}

		parent.Children = append(parent.Children, *cur)
	}

	if root == nil {
		return TreeNode{}, gorm.ErrRecordNotFound
	}
	return *root, nil
}
