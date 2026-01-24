package nod

import (
	"fmt"
	"log/slog"
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
	log            *slog.Logger
	db             *gorm.DB
	nodeIds        []string
	parentIds      []string
	namespaceIds   []string
	name           *StringFilter
	type_          *StringFilter
	kind           *StringFilter
	status         *StringFilter
	createdDate    *TimeFilter
	updatedDate    *TimeFilter
	includeTags    bool
	includeKV      bool
	includeContent bool
	excludeRoot    bool
	onlyRoots      bool
	limit          int
	page           int
	pageSize       int
	mappers        *MapperRegistry
}

type TreeNode struct {
	Node     NodeModel
	Children []*TreeNode
}

func NewNodeQuery(db *gorm.DB, log *slog.Logger, mappers *MapperRegistry) *NodeQuery {
	return &NodeQuery{
		db:      db,
		log:     log,
		mappers: mappers,
	}
}

func (q *NodeQuery) Clone() *NodeQuery {
	clone := *q
	clone.nodeIds = append([]string{}, q.nodeIds...)
	clone.parentIds = append([]string{}, q.parentIds...)
	clone.namespaceIds = append([]string{}, q.namespaceIds...)
	return &clone
}

func StringEquals(value string) *StringFilter {
	return &StringFilter{Equals: &value}
}

func StringContains(value string) *StringFilter {
	return &StringFilter{Contains: &value}
}

func StringStartsWith(value string) *StringFilter {
	return &StringFilter{StartsWith: &value}
}

func StringEndsWith(value string) *StringFilter {
	return &StringFilter{EndsWith: &value}
}

func TimeFrom(value time.Time) *TimeFilter {
	return &TimeFilter{From: &value}
}

func TimeTo(value time.Time) *TimeFilter {
	return &TimeFilter{To: &value}
}

func TimeBetween(from, to time.Time) *TimeFilter {
	return &TimeFilter{From: &from, To: &to}
}

func NewStringFilter(equals, contains, startsWith, endsWith *string) *StringFilter {
	return &StringFilter{
		Equals:     equals,
		Contains:   contains,
		StartsWith: startsWith,
		EndsWith:   endsWith,
	}
}

func NewTimeFilter(from, to *time.Time) *TimeFilter {
	return &TimeFilter{
		From: from,
		To:   to,
	}
}

func (q *NodeQuery) Roots() *NodeQuery {
	q.onlyRoots = true
	return q
}

func (q *NodeQuery) ExcludeRoot() *NodeQuery {
	q.excludeRoot = true
	return q
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

func (q *NodeQuery) Content() *NodeQuery {
	q.includeContent = true
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

func (q *NodeQuery) NameEquals(value string) *NodeQuery {
	q.name = &StringFilter{Equals: &value}
	return q
}

func (q *NodeQuery) NameContains(value string) *NodeQuery {
	q.name = &StringFilter{Contains: &value}
	return q
}

func (q *NodeQuery) NameStartsWith(value string) *NodeQuery {
	q.name = &StringFilter{StartsWith: &value}
	return q
}

func (q *NodeQuery) NameEndsWith(value string) *NodeQuery {
	q.name = &StringFilter{EndsWith: &value}
	return q
}

func (q *NodeQuery) Type(filter *StringFilter) *NodeQuery {
	q.type_ = filter
	return q
}

func (q *NodeQuery) TypeEquals(value string) *NodeQuery {
	q.type_ = &StringFilter{Equals: &value}
	return q
}

func (q *NodeQuery) TypeContains(value string) *NodeQuery {
	q.type_ = &StringFilter{Contains: &value}
	return q
}

func (q *NodeQuery) TypeStartsWith(value string) *NodeQuery {
	q.type_ = &StringFilter{StartsWith: &value}
	return q
}

func (q *NodeQuery) TypeEndsWith(value string) *NodeQuery {
	q.type_ = &StringFilter{EndsWith: &value}
	return q
}

func (q *NodeQuery) Kind(filter *StringFilter) *NodeQuery {
	q.kind = filter
	return q
}

func (q *NodeQuery) KindEquals(value string) *NodeQuery {
	q.kind = &StringFilter{Equals: &value}
	return q
}

func (q *NodeQuery) KindContains(value string) *NodeQuery {
	q.kind = &StringFilter{Contains: &value}
	return q
}

func (q *NodeQuery) KindStartsWith(value string) *NodeQuery {
	q.kind = &StringFilter{StartsWith: &value}
	return q
}

func (q *NodeQuery) KindEndsWith(value string) *NodeQuery {
	q.kind = &StringFilter{EndsWith: &value}
	return q
}

func (q *NodeQuery) Status(filter *StringFilter) *NodeQuery {
	q.status = filter
	return q
}

func (q *NodeQuery) StatusEquals(value string) *NodeQuery {
	q.status = &StringFilter{Equals: &value}
	return q
}

func (q *NodeQuery) StatusContains(value string) *NodeQuery {
	q.status = &StringFilter{Contains: &value}
	return q
}

func (q *NodeQuery) StatusStartsWith(value string) *NodeQuery {
	q.status = &StringFilter{StartsWith: &value}
	return q
}

func (q *NodeQuery) StatusEndsWith(value string) *NodeQuery {
	q.status = &StringFilter{EndsWith: &value}
	return q
}

func (q *NodeQuery) CreatedDate(filter *TimeFilter) *NodeQuery {
	q.createdDate = filter
	return q
}

func (q *NodeQuery) CreatedDateFrom(from time.Time) *NodeQuery {
	if q.createdDate == nil {
		q.createdDate = &TimeFilter{}
	}
	q.createdDate.From = &from
	return q
}

func (q *NodeQuery) CreatedDateTo(to time.Time) *NodeQuery {
	if q.createdDate == nil {
		q.createdDate = &TimeFilter{}
	}
	q.createdDate.To = &to
	return q
}

func (q *NodeQuery) CreatedBetween(from, to time.Time) *NodeQuery {
	q.createdDate = &TimeFilter{From: &from, To: &to}
	return q
}

func (q *NodeQuery) UpdatedDate(filter *TimeFilter) *NodeQuery {
	q.updatedDate = filter
	return q
}

func (q *NodeQuery) UpdatedDateFrom(from time.Time) *NodeQuery {
	if q.updatedDate == nil {
		q.updatedDate = &TimeFilter{}
	}
	q.updatedDate.From = &from
	return q
}

func (q *NodeQuery) UpdatedDateTo(to time.Time) *NodeQuery {
	if q.updatedDate == nil {
		q.updatedDate = &TimeFilter{}
	}
	q.updatedDate.To = &to
	return q
}

func (q *NodeQuery) UpdatedBetween(from, to time.Time) *NodeQuery {
	q.updatedDate = &TimeFilter{From: &from, To: &to}
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

func ApplyCommonFilters(db *gorm.DB, t *NodeQuery) *gorm.DB {
	if len(t.nodeIds) > 0 {
		db = db.Where("id IN ?", t.nodeIds)
	}
	if t.onlyRoots {
		db = db.Where("parent_id IS NULL or parent_id = \"\"")
	}
	if t.excludeRoot {
		db = db.Where("parent_id IS NOT NULL")
	}
	if len(t.parentIds) > 0 {
		db = db.Where("parent_id IN ?", t.parentIds)
	}
	if len(t.namespaceIds) > 0 {
		db = db.Where("namespace_id IN ?", t.namespaceIds)
	}
	if t.name != nil {
		db = ApplyStringFilter(db, "name", t.name)
	}
	if t.type_ != nil {
		db = ApplyStringFilter(db, "type", t.type_)
	}
	if t.kind != nil {
		db = ApplyStringFilter(db, "kind", t.kind)
	}
	if t.status != nil {
		db = ApplyStringFilter(db, "status", t.status)
	}
	if t.createdDate != nil {
		db = ApplyTimeFilter(db, "created_at", t.createdDate)
	}
	if t.updatedDate != nil {
		db = ApplyTimeFilter(db, "updated_at", t.updatedDate)
	}
	return db
}

func (q *NodeQuery) ApplyConditions(db *gorm.DB) *gorm.DB {
	q.log.Debug(fmt.Sprintf("TypedQuery current filters: nodeIds=%v, parentIds=%v, namespaceIds=%v, name=%v, type_=%v, kind=%v, status=%v, createdDate=%v, updatedDate=%v, onlyRoots=%v, excludeRoot=%v",
		q.nodeIds, q.parentIds, q.namespaceIds, q.name, q.type_, q.kind, q.status, q.createdDate, q.updatedDate, q.onlyRoots, q.excludeRoot))

	db = ApplyCommonFilters(db, q)

	if q.limit > 0 {
		db = db.Limit(q.limit)
	}
	if q.page > 0 && q.pageSize > 0 {
		offset := (q.page - 1) * q.pageSize
		db = db.Offset(offset).Limit(q.pageSize)
	}

	return db
}

func (q *NodeQuery) fetchNodes() ([]*Node, error) {
	db := q.db.Model(&NodeCore{})
	q.log.Debug("NodeQuery FindAll: starting query")
	db = q.ApplyConditions(db)

	var nodeCores []NodeCore
	if err := db.Find(&nodeCores).Error; err != nil {
		return nil, err
	}

	q.log.Debug("NodeQuery FindAll: retrieved node cores", slog.Int("count", len(nodeCores)))
	var nodes []*Node
	for _, nc := range nodeCores {
		nodes = append(nodes, &Node{
			Core: nc})
	}
	q.log.Debug("NodeQuery FindAll: constructed nodes", slog.Int("count", len(nodes)))

	if q.includeTags {
		q.log.Debug("NodeQuery FindAll: loading tags for nodes")
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Tags = tagsByNode[n.Core.Id]
		}
		q.log.Debug("NodeQuery FindAll: loaded tags for nodes")
	}

	if q.includeKV {
		q.log.Debug("NodeQuery FindAll: loading KV for nodes")
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.KV = kvsByNode[n.Core.Id]
		}
		q.log.Debug("NodeQuery FindAll: loaded KV for nodes")
	}

	if q.includeContent {
		q.log.Debug("NodeQuery FindAll: loading Content for nodes")
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		contentsByNode, err := (&ContentRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			n.Content = contentsByNode[n.Core.Id]
		}
		q.log.Debug("NodeQuery FindAll: loaded Content for nodes")
	}
	return nodes, nil
}

func (q *NodeQuery) List() ([]NodeModel, error) {
	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	var results []NodeModel
	for _, n := range nodes {
		mapper, err := q.mappers.ForNode(n)
		if err != nil {
			return nil, err
		}
		mapped, err := mapper.FromNode(n)
		if err != nil {
			return nil, err
		}
		results = append(results, mapped)
	}
	return results, nil
}

func (q *NodeQuery) First() (NodeModel, error) {
	nodes, err := q.Limit(1).List()
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return nodes[0], nil
}

func (q *NodeQuery) Count() (int64, error) {
	db := q.db.Model(&NodeCore{})

	db = ApplyCommonFilters(db, q)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *NodeQuery) Descendants(onlyRoots bool) ([]*TreeNode, error) {
	trees := make([]*TreeNode, 0)

	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		q.log.Debug("Debug: Processing node ID =%v, with Parent ID =%v", n.Core.Id, n.Core.ParentId)
		if !onlyRoots && (n.Core.ParentId == nil || *n.Core.ParentId == "") {
			tree, err := q.buildTree(n.Core.Id)
			if err != nil {
				return nil, err
			}

			mappedTree := &TreeNode{
				Node:     tree.Node,
				Children: tree.Children,
			}

			trees = append(trees, mappedTree)
		}
	}

	return trees, nil
}

func (q *NodeQuery) DescendantTree(rootID string) (*TreeNode, error) {
	return q.buildTree(rootID)
}

func (q *NodeQuery) buildTree(rootID string) (*TreeNode, error) {
	db := q.db.Model(&NodeCore{})
	var nodeCores []NodeCore
	var nodes []*Node

	sql := `
WITH RECURSIVE tree AS (
  SELECT * FROM node_cores WHERE id = ?
  UNION ALL
  SELECT n.* FROM node_cores n
  JOIN tree t ON n.parent_id = t.id
)
SELECT * FROM tree;
`
	if err := db.Raw(sql, rootID).Scan(&nodeCores).Error; err != nil {
		return nil, err
	}
	if len(nodeCores) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, nc := range nodeCores {
		nodes = append(nodes, &Node{
			Core: nc,
		})
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}

		for _, n := range nodes {
			n.Tags = tagsByNode[n.Core.Id]
		}
	}

	if q.includeKV {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i, n := range nodes {
			nodes[i].KV = kvsByNode[n.Core.Id]
		}
	}

	if q.includeContent {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		contentsByNode, err := (&ContentRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i, n := range nodes {
			nodes[i].Content = contentsByNode[n.Core.Id]
		}
	}

	byID := make(map[string]*TreeNode, len(nodes))
	for _, n := range nodes {
		mapper, err := q.mappers.ForNode(n)
		if err != nil {
			return nil, err
		}
		mappedNode, err := mapper.FromNode(n)
		if err != nil {
			return nil, err
		}

		byID[n.Core.Id] = &TreeNode{
			Node: mappedNode,
		}
	}

	var root *TreeNode
	for _, n := range nodes {
		cur := byID[n.Core.Id]
		if n.Core.Id == rootID {
			root = cur
			continue
		}
		if n.Core.ParentId == nil {
			continue
		}
		parent := byID[*n.Core.ParentId]
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

func (q *NodeQuery) buildAncestorTree(childID string) (*TreeNode, error) {
	db := q.db.Model(&NodeCore{})
	var nodeCores []NodeCore
	var nodes []*Node

	sql := `
WITH RECURSIVE path AS (
  SELECT * FROM node_cores WHERE id = ?
  UNION ALL
  SELECT p.* FROM node_cores p
  JOIN path c ON p.id = c.parent_id
)
SELECT * FROM path;
`
	if err := db.Raw(sql, childID).Scan(&nodeCores).Error; err != nil {
		return nil, err
	}
	if len(nodeCores) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, nc := range nodeCores {
		nodes = append(nodes, &Node{
			Core: nc,
		})
	}

	if q.includeTags {
		tagsByNode, err := loadTagsByNode(q.db, nodes)
		if err != nil {
			return nil, err
		}

		for _, n := range nodes {
			n.Tags = tagsByNode[n.Core.Id]
		}
	}

	if q.includeKV {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		kvsByNode, err := (&KVRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i, n := range nodes {
			nodes[i].KV = kvsByNode[n.Core.Id]
		}
	}

	if q.includeContent {
		nodeIds := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodeIds = append(nodeIds, n.Core.Id)
		}
		contentsByNode, err := (&ContentRepository{DB: q.db}).GetAllForNodes(nodeIds)
		if err != nil {
			return nil, err
		}
		for i, n := range nodes {
			nodes[i].Content = contentsByNode[n.Core.Id]
		}
	}

	q.log.Debug("Debug: Building ancestor tree for childID = %v", childID)
	q.log.Debug("Debug: Retrieved nodes:")
	for _, n := range nodes {
		q.log.Debug("  Node ID: %v, Parent ID: %v\n", n.Core.Id, *n.Core.ParentId)
	}

	byID := make(map[string]*TreeNode, len(nodes))
	for _, n := range nodes {
		mapper, err := q.mappers.ForNode(n)
		if err != nil {
			return nil, err
		}
		mappedNode, err := mapper.FromNode(n)
		if err != nil {
			return nil, err
		}
		byID[n.Core.Id] = &TreeNode{
			Node: mappedNode,
		}
	}

	var root *TreeNode
	for _, n := range nodes {
		cur := byID[n.Core.Id]
		if n.Core.ParentId == nil || *n.Core.ParentId == "" {
			root = cur
			continue
		}
		parent := byID[*n.Core.ParentId]

		if parent == nil {
			continue
		}

		if parent.Children == nil {
			parent.Children = make([]*TreeNode, 0)
		}
		parent.Children = append(parent.Children, cur)
	}

	if root == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil
}

func (q *NodeQuery) Ancestors() ([]*TreeNode, error) {
	trees := make([]*TreeNode, 0)

	nodes, err := q.fetchNodes()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		tree, err := q.buildAncestorTree(n.Core.Id)
		if err != nil {
			return nil, err
		}

		trees = append(trees, tree)
	}
	return trees, nil
}

func (q *NodeQuery) AncestorTree(childID string) (*TreeNode, error) {
	return q.buildAncestorTree(childID)
}

func (q *NodeQuery) Exists() (bool, error) {
	db := q.db.Model(&NodeCore{})
	db = q.ApplyConditions(db)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (q *NodeQuery) HasChildren() bool {
	db := q.db.Model(&NodeCore{})
	db = q.ApplyConditions(db)
	var parents []NodeCore
	if err := db.Find(&parents).Error; err != nil {
		return false
	}
	if len(parents) == 0 {
		return false
	}

	db = q.db.Model(&NodeCore{})
	db = db.Where("parent_id IN ?", func() []string {
		ids := make([]string, 0, len(parents))
		for _, p := range parents {
			ids = append(ids, p.Id)
		}
		return ids
	}())

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (q *NodeQuery) Delete() error {
	db := q.db.Model(&NodeCore{})
	q.log.Debug("NodeQuery Delete: starting query")
	db = q.ApplyConditions(db)

	if q.HasChildren() {
		return fmt.Errorf("cannot delete nodes that have children")
	}

	var nodeCores []NodeCore
	if err := db.Delete(&nodeCores).Error; err != nil {
		return err
	}

	q.log.Debug("NodeQuery Delete: nodes deleted")
	return nil
}
