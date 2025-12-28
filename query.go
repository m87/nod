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
	nodeIds      []string
	parentIds    []string
	namespaceIds []string
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

func (q *NodeQuery) NodeId(nodeId string) *NodeQuery {
	q.nodeIds = append(q.nodeIds, nodeId)
	return q
}

func (q *NodeQuery) ParentId(parentId string) *NodeQuery {
	q.parentIds = append(q.parentIds, parentId)
	return q
}

func (q *NodeQuery) NamespaceId(namespaceId string) *NodeQuery {
	q.namespaceIds = append(q.namespaceIds, namespaceId)
	return q
}

func (q *NodeQuery) NodeIds(nodeIds []string) *NodeQuery {
	q.nodeIds = append(q.nodeIds, nodeIds...)
	return q
}

func (q *NodeQuery) ParentIds(parentIds []string) *NodeQuery {
	q.parentIds = append(q.parentIds, parentIds...)
	return q
}

func (q *NodeQuery) NamespaceIds(namespaceIds []string) *NodeQuery {
	q.namespaceIds = append(q.namespaceIds, namespaceIds...)
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

func (q *NodeQuery) FindAll() ([]*Node, error) {
	db := q.db.Model(&Node{})

	if len(q.nodeIds) > 0 {
		db = db.Where("id IN ?", q.nodeIds)
	}
	if len(q.parentIds) > 0 {
		db = db.Where("parent_id IN ?", q.parentIds)
	}
	if len(q.namespaceIds) > 0 {
		db = db.Where("namespace_id IN ?", q.namespaceIds)
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
	}
	if q.page > 0 && q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	var nodes []*Node
	if err := db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Tags = tagsByNode[n.Id]
		}
	}

	if q.includeKV {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.KV = kvsByNode[n.Id]
		}
	}

	return nodes, nil
}

func (q *NodeQuery) Find() (*Node, error) {
	nodes, err := q.Limit(1).FindAll()
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return nodes[0], nil
}

func (q *NodeQuery) Count() (int64, error) {
	db := q.db.Model(&Node{})

	if len(q.nodeIds) > 0 {
		db = db.Where("id IN ?", q.nodeIds)
	}
	if len(q.parentIds) > 0 {
		db = db.Where("parent_id IN ?", q.parentIds)
	}
	if len(q.namespaceIds) > 0 {
		db = db.Where("namespace_id IN ?", q.namespaceIds)
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

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *NodeQuery) Decendants() ([]*TreeNode, error) {
	trees := make([]*TreeNode, 0)
	
	nodes, err := q.FindAll()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		tree, err := q.buildTree(n.Id)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func (q *NodeQuery) DecendantTree(rootID string) (*TreeNode, error) {
	return q.buildTree(rootID)
}

func (q *NodeQuery) buildTree(rootID string) (*TreeNode, error) {
	  db := q.db.Model(&Node{})
		var nodes []*Node

	sql := `
WITH RECURSIVE tree AS (
  SELECT * FROM nodes WHERE id = ?
  UNION ALL
  SELECT n.* FROM nodes n
  JOIN tree t ON n.parent_id = t.id
)
SELECT * FROM tree;
`
	if err := db.Raw(sql, rootID).Scan(&nodes).Error; err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}

		for _, n := range nodes {
			n.Tags = tagsByNode[n.Id]
		}
	}

	if q.includeKV {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i, n := range nodes {
			nodes[i].KV = kvsByNode[n.Id]
		}
	}
	
	byID := make(map[string]*TreeNode, len(nodes))
	for _, n := range nodes {
		byID[n.Id] = &TreeNode{
			Node: n,
		}
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
		parent.Children = append(parent.Children, cur)
	}

	if root == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil

}

