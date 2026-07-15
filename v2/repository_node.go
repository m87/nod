package nod

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NodeScope is a generic struct that provides methods for managing nodes of type T within a repository.
type NodeScope[T any] struct {
	repository *Repository
}

// Nodes returns a NodeScope for the given repository, allowing for operations on nodes.
func (repository *Repository) Nodes() *NodeScope[Node] {
	return &NodeScope[Node]{
		repository: repository,
	}
}

// Nodes is a generic function that returns a NodeScope for the given repository, allowing for operations on nodes of type T.
func Nodes[T any](repository *Repository) *NodeScope[T] {
	return &NodeScope[T]{
		repository: repository,
	}
}

// SaveNode saves the given node to the repository.
func (scope *NodeScope[T]) SaveNode(model *T) (string, error) {
	if model == nil {
		return "", NewNodeIsNilError()
	}

	node, err := nodeFromModel(scope.repository.adapters, model)
	if err != nil {
		return "", err
	}

	id := ensureNodeID(node)

	err = scope.repository.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&node.Core).Error; err != nil {
			return err
		}

		if err := deleteNodeContents(tx, node.Core.Id); err != nil {
			return err
		}

		content := []*NodeContent{}

		for _, value := range node.Content {
			value.NodeId = id
			content = append(content, value)
		}

		if err := saveNodeContents(tx, content); err != nil {
			return err
		}

		if err := unbindNodeTagsFromNode(tx, node.Core.Id); err != nil {
			return err
		}

		for _, tag := range node.Tags {
			savedTag, err := saveNodeTagIfNotExists(tx, node.Core.NamespaceId, tag.Name)
			if err != nil {
				return err
			}
			if err := bindNodeTagToNode(tx, id, savedTag.Id); err != nil {
				return err
			}
		}

		return nil
	})

	return id, err
}

// DeleteNode deletes the given node from the repository.
func (scope *NodeScope[T]) DeleteNode(model *T) error {
	if model == nil {
		return NewNodeIsNilError()
	}

	node, err := nodeFromModel(scope.repository.adapters, model)
	if err != nil {
		return err
	}
	return scope.repository.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&node.Core).Error
	})
}

func (scope *NodeScope[T]) GetNode(id string) (*T, error) {
	node := &Node{}
	err := scope.repository.db.First(&node.Core, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	content, err := scope.repository.getNodeContents(id)

	contentMap := make(map[string]*NodeContent)
	for _, c := range content {
		contentMap[c.Key] = c
	}
	node.Content = contentMap

	tags, err := scope.repository.getNodeTags(id)
	if err != nil {
		return nil, err
	}
	node.Tags = tags

	return modelFromNode[T](scope.repository.adapters, node)
}

func ensureNodeID(node *Node) string {
	if node.Core.Id == "" {
		node.Core.Id = uuid.New().String()
	}
	return node.Core.Id
}
