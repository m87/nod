package nod

import "gorm.io/gorm"


type NodeQuery struct {
	db *gorm.DB
	nodeId string
	parentId *string
	namespaceId *string
	includeTags bool
	includeKV   bool
	limit int
	page int
	pageSize int
}

func Query(db *gorm.DB) *NodeQuery {
	return &NodeQuery{
		db: db,
	}
}

func (q *NodeQuery) Id(nodeId string) *NodeQuery {
	q.nodeId = nodeId
	return q
}

func (q *NodeQuery) ParentId(parentId string) *NodeQuery {
	q.parentId = &parentId
	return q
}

func (q *NodeQuery) NamespaceId(namespaceId string) *NodeQuery {
	q.namespaceId = &namespaceId
	return q
}

func (q *NodeQuery) Tags() *NodeQuery {
	q.includeTags = true
	return q
}

func (q *NodeQuery) KV() *NodeQuery {
	q.includeKV = true
	return q
}

func (q *NodeQuery) Limit(limit int) *NodeQuery {
	q.limit = limit
	return q
}

func (q *NodeQuery) Page(page int, pageSize int) *NodeQuery {
	q.page = page
	q.pageSize = pageSize
	return q
}

func (q *NodeQuery) Find() (Node , error) {
	var node Node
	db := q.db.Model(&Node{})

	if q.nodeId != "" {
		db = db.Where("id = ?", q.nodeId)
	}

	if q.parentId != nil {
		db = db.Where("parent_id = ?", *q.parentId)
	}

	if q.namespaceId != nil {
		db = db.Where("namespace_id = ?", *q.namespaceId)
	}

	if q.limit > 0 {
		db = db.Limit(q.limit)
	}

	if q.page > 0 && q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	if err := db.First(&node).Error; err != nil {
		return Node{}, err
	}

	if q.includeTags {
		var tags []Tag
		if err := q.db.Model(&Tag{}).
			Joins("JOIN node_tags ON node_tags.tag_id = tags.id").
			Where("node_tags.node_id = ?", node.Id).
			Find(&tags).Error; err != nil {
			return Node{}, err
		}
		node.Tags = tags
	}

	if q.includeKV {
		kvRepo := &KVRepository{DB: q.db}
		kvs, err := kvRepo.GetAll(node.Id)
		if err != nil {
			return Node{}, err
		}
		node.KV = kvs
	}

	return node, nil
}


func (q *NodeQuery) FindAll() ([]Node, error) {
	var nodes []Node
	db := q.db.Model(&Node{})

	if q.nodeId != "" {
		db = db.Where("id = ?", q.nodeId)
	}

	if q.parentId != nil {
		db = db.Where("parent_id = ?", *q.parentId)
	}

	if q.namespaceId != nil {
		db = db.Where("namespace_id = ?", *q.namespaceId)
	}

	if q.limit > 0 {
		db = db.Limit(q.limit)
	}

	if q.page > 0 && q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	if err := db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for i := range nodes {
			nodes[i].Tags = tagsByNode[nodes[i].Id]
		}
	}

	if q.includeKV {
		nodeIds := make([]string, len(nodes))
		for i, n := range nodes {
			nodeIds[i] = n.Id
		}
		kvRepo := &KVRepository{DB: q.db}
		kvsByNode, err := kvRepo.GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i := range nodes {
			nodes[i].KV = kvsByNode[nodes[i].Id]
		}
	}

	return nodes, nil
}

