package nod

import "gorm.io/gorm"

// NodeQuery represents a query for nodes in the repository, allowing for filtering based on various criteria.
type NodeQuery struct {
	repository *Repository

	nodeIds      []string
	parentIds    []string
	namespaceIds []string
	nameFilters []*Filter[string]
}

// NewNodeQuery creates a new NodeQuery for the given repository.
func NewNodeQuery(repository *Repository) *NodeQuery {
	return &NodeQuery{
		repository: repository,
	}
}

// FindAll retrieves all nodes matching the query criteria, with pagination support.
func (query *NodeQuery) FindAll(page int, pageSize int) ([]*Node, error) {
	db := query.repository.db.Model(&NodeCore{})
	db = query.applyCoreFitlers(db)

	var nodeCores []NodeCore
	result := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&nodeCores)
	if result.Error != nil {
		return nil, result.Error
	}

	nodes := make([]*Node, len(nodeCores))
	for i, nodeCore := range nodeCores {
		nodes[i] = &Node{
			Core: nodeCore,
		}
	}
	return nodes, nil
}

// FindOne retrieves a single node matching the query criteria. If multiple nodes match, it returns error.
func (query *NodeQuery) FindOne() (*Node, error) {
	db := query.repository.db.Model(&NodeCore{})
	db = query.applyCoreFitlers(db)

	var nodeCore NodeCore
	result := db.First(&nodeCore)
	if result.Error != nil {
		return nil, result.Error
	}

	node := &Node{
		Core: nodeCore,
	}
	return node, nil
}

// Count returns the total number of nodes matching the query criteria.
func (query *NodeQuery) Count() (int64, error) {
	db := query.repository.db.Model(&NodeCore{})
	db = query.applyCoreFitlers(db)

	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// Edges retrieves all edges connected to the nodes matching the query criteria.
func (query *NodeQuery) Edges() ([]*Edge, error) {
	return nil, nil
}

// OutEdges retrieves all outgoing edges from the nodes matching the query criteria.
func (query *NodeQuery) OutEdges() ([]*Edge, error) {
	return nil, nil
}

// InEdges retrieves all incoming edges to the nodes matching the query criteria.
func (query *NodeQuery) InEdges() ([]*Edge, error) {
	return nil, nil
}

// NodeIds adds multiple node IDs to the query for filtering.
func (query *NodeQuery) NodeIds(nodeIds []string) *NodeQuery {
	query.nodeIds = append(query.nodeIds, nodeIds...)
	return query
}

// NodeId adds a single node ID to the query for filtering.
func (query *NodeQuery) NodeId(nodeId string) *NodeQuery {
	query.nodeIds = append(query.nodeIds, nodeId)
	return query
}

// ParentIds adds multiple parent IDs to the query for filtering.
func (query *NodeQuery) ParentIds(parentIds []string) *NodeQuery {
	query.parentIds = append(query.parentIds, parentIds...)
	return query
}

// ParentId adds a single parent ID to the query for filtering.
func (query *NodeQuery) ParentId(parentId string) *NodeQuery {
	query.parentIds = append(query.parentIds, parentId)
	return query
}

// NamespaceIds adds multiple namespace IDs to the query for filtering.
func (query *NodeQuery) NamespaceIds(namespaceIds []string) *NodeQuery {
	query.namespaceIds = append(query.namespaceIds, namespaceIds...)
	return query
}

// NamespaceId adds a single namespace ID to the query for filtering.
func (query *NodeQuery) NamespaceId(namespaceId string) *NodeQuery {
	query.namespaceIds = append(query.namespaceIds, namespaceId)
	return query
}

// NameFilter adds a name filter to the query for filtering nodes based on their names.
func (query *NodeQuery) NameFilter(filter *Filter[string]) *NodeQuery {
	query.nameFilters = append(query.nameFilters, filter)
	return query
}

// NameFilters adds multiple name filters to the query for filtering nodes based on their names.
func (query *NodeQuery) NameFilters(filters []*Filter[string]) *NodeQuery {
	query.nameFilters = append(query.nameFilters, filters...)
	return query
}

// Name adds a name filter to the query for filtering nodes based on their names, using an equality comparison.
func (query *NodeQuery) Name(value string) *NodeQuery {
	query.nameFilters = append(query.nameFilters, Equals(value))
	return query
}

func (query *NodeQuery) applyCoreFitlers(db *gorm.DB) *gorm.DB {
	if len(query.nodeIds) > 0 {
		db = db.Where("id IN ?", query.nodeIds)
	}
	if len(query.parentIds) > 0 {
		db = db.Where("parent_id IN ?", query.parentIds)
	}
	if len(query.namespaceIds) > 0 {
		db = db.Where("namespace_id IN ?", query.namespaceIds)
	}
	if len(query.nameFilters) > 0 {
		for _, filter := range query.nameFilters {
			applyStringFilter(db, "name", filter)
		}
	}
	return db
}



