package nod

import (
	"time"

	"gorm.io/gorm"
)

type TimeFilter struct {
	From *time.Time
	To   *time.Time
}

type StringFilter struct {
	Equals     *string
	Contains   *string
	StartsWith *string
	EndsWith   *string
}

type NodeQuery struct {
	db          *gorm.DB
	nodeId      *string
	parentId    *string
	namespaceId *string
	name        *StringFilter
	type_       *StringFilter
	kind        *StringFilter
	status      *StringFilter
	createdDate *TimeFilter
	updatedDate *TimeFilter
	includeTags bool
	includeKV   bool
	limit       int
	page        int
	pageSize    int
}

func Query(db *gorm.DB) *NodeQuery {
	return &NodeQuery{
		db: db,
	}
}

func (q *NodeQuery) Id(nodeId string) *NodeQuery {
	q.nodeId = &nodeId
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

func (q *NodeQuery) Name(filter *StringFilter) *NodeQuery {
	q.name = filter
	return q
}

func (q *NodeQuery) Type(filter *StringFilter) *NodeQuery {
	q.type_ = filter
	return q
}

func (q *NodeQuery) Kind(filter *StringFilter) *NodeQuery {
	q.kind = filter
	return q
}

func (q *NodeQuery) Status(filter *StringFilter) *NodeQuery {
	q.status = filter
	return q
}

func (q *NodeQuery) CreatedDate(filter *TimeFilter) *NodeQuery {
	q.createdDate = filter
	return q
}

func (q *NodeQuery) UpdatedDate(filter *TimeFilter) *NodeQuery {
	q.updatedDate = filter
	return q
}

func ApplyStringFilter(db *gorm.DB, field string, filter *StringFilter) *gorm.DB {
	if filter.Equals != nil {
		db = db.Where(field+" = ?", *filter.Equals)
	}
	if filter.Contains != nil {
		db = db.Where(field+" LIKE ?", "%"+*filter.Contains+"%")
	}
	if filter.StartsWith != nil {
		db = db.Where(field+" LIKE ?", *filter.StartsWith+"%")
	}
	if filter.EndsWith != nil {
		db = db.Where(field+" LIKE ?", "%"+*filter.EndsWith)
	}
	return db
}

func ApplyTimeFilter(db *gorm.DB, field string, filter *TimeFilter) *gorm.DB { 
	if filter.From != nil {
		db = db.Where(field+" >= ?", *filter.From)
	}
	if filter.To != nil {
		db = db.Where(field+" <= ?", *filter.To)
	}
	return db
}

func (q *NodeQuery) Find() (*Node, error) {
	var node Node
	db := q.db.Model(&Node{})

	if q.nodeId != nil {
		db = db.Where("id = ?", *q.nodeId)
	}
	if q.parentId != nil {
		db = db.Where("parent_id = ?", *q.parentId)
	}
	if q.namespaceId != nil {
		db = db.Where("namespace_id = ?", *q.namespaceId)
	}
	if q.name != nil {
		db = ApplyStringFilter(db, "name", q.name)
	}
	if q.type_ != nil {
		db = ApplyStringFilter(db, "type", q.type_)
	}
	if q.kind != nil {
		db = ApplyStringFilter(db, "kind", q.kind)
	}
	if q.status != nil {
		db = ApplyStringFilter(db, "status", q.status)
	}
	if q.createdDate != nil {
		db = ApplyTimeFilter(db, "created_at", q.createdDate)
	}
	if q.updatedDate != nil {
		db = ApplyTimeFilter(db, "updated_at", q.updatedDate)
	}

	if q.limit > 0 {
		db = db.Limit(q.limit)
	} else if q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	
	if err := db.First(&node).Error; err != nil {
		return nil, err
	}

	if q.includeTags {
		var tags []Tag
		if err := q.db.Model(&Tag{}).
			Joins("JOIN node_tags ON node_tags.tag_id = tags.id").
			Where("node_tags.node_id = ?", node.Id).
			Find(&tags).Error; err != nil {
			return nil, err
		}
		node.Tags = tags
	}

	if q.includeKV {
		kvRepo := &KVRepository{DB: q.db}
		kvs, err := kvRepo.GetAll(node.Id)
		if err != nil {
			return nil, err
		}
		node.KV = kvs
	}

	return &node, nil
}

func (q *NodeQuery) FindAll() ([]*Node, error) {
	var nodes []Node
	db := q.db.Model(&Node{})

	if q.nodeId != nil {
		db = db.Where("id = ?", *q.nodeId)
	}
	if q.parentId != nil {
		db = db.Where("parent_id = ?", *q.parentId)
	}
	if q.namespaceId != nil {
		db = db.Where("namespace_id = ?", *q.namespaceId)
	}
	if q.name != nil {
		db = ApplyStringFilter(db, "name", q.name)
	}
	if q.type_ != nil {
		db = ApplyStringFilter(db, "type", q.type_)
	}	
	if q.kind != nil {
		db = ApplyStringFilter(db, "kind", q.kind)
	}
	if q.status != nil {
		db = ApplyStringFilter(db, "status", q.status)
	}
	if q.createdDate != nil {
		db = ApplyTimeFilter(db, "created_at", q.createdDate)
	}
	if q.updatedDate != nil {
		db = ApplyTimeFilter(db, "updated_at", q.updatedDate)
	}

	if q.limit > 0 {
		db = db.Limit(q.limit)
	} else if q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	if err := db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	result := make([]*Node, len(nodes))
	for i, n := range nodes {
		result[i] = &n
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range result {
			n.Tags = tagsByNode[n.Id]
		}
	}

	if q.includeKV {
		kvRepo := &KVRepository{DB: q.db}
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Id)
		}
		kvsByNode, err := kvRepo.GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for _, n := range result {
			n.KV = kvsByNode[n.Id]
		}
	}

	return result, nil
}

