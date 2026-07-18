package nod

import "gorm.io/gorm"

func saveNodeKv(tx *gorm.DB, kv *NodeKV) error {
	if kv == nil {
		return NewNodeKVIsNilError()
	}
	return tx.Save(kv).Error
}

func saveNodeKvs(tx *gorm.DB, kvs []*NodeKV) error {
	for _, kv := range kvs {
		if err := saveNodeKv(tx, kv); err != nil {
			return err
		}
	}
	return nil
}

func deleteNodeKvs(tx *gorm.DB, nodeId string) error {
	return tx.Where("node_id = ?", nodeId).Delete(&NodeKV{}).Error
}

func (r *Repository) getNodeKvs(nodeId string) ([]*NodeKV, error) {
	var kvs []*NodeKV
	err := r.db.Where("node_id = ?", nodeId).Find(&kvs).Error
	if err != nil {
		return nil, err
	}
	return kvs, nil
}

func (r *Repository) getNodesKvs(nodeIds []string) (map[string][]*NodeKV, error) {
	var kvs []*NodeKV
	err := r.db.Where("node_id IN ?", nodeIds).Find(&kvs).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*NodeKV)
	for _, kv := range kvs {
		result[kv.NodeId] = append(result[kv.NodeId], kv)
	}
	return result, nil
}

func (r *Repository) getEdgesKvs(edgeIds []string) (map[string][]*EdgeKV, error) {
	var kvs []*EdgeKV
	err := r.db.Where("edge_id IN ?", edgeIds).Find(&kvs).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*EdgeKV)
	for _, kv := range kvs {
		result[kv.EdgeId] = append(result[kv.EdgeId], kv)
	}
	return result, nil
}

func (r *Repository) getEdgeKvs(edgeId string) ([]*EdgeKV, error) {
	var kvs []*EdgeKV
	err := r.db.Where("edge_id = ?", edgeId).Find(&kvs).Error
	if err != nil {
		return nil, err
	}
	return kvs, nil
}

func savneEdgeKv(tx *gorm.DB, kv *EdgeKV) error {
	if kv == nil {
		return NewEdgeKVIsNilError()
	}
	return tx.Save(kv).Error
}

func saveEdgeKvs(tx *gorm.DB, kvs []*EdgeKV) error {
	for _, kv := range kvs {
		if err := savneEdgeKv(tx, kv); err != nil {
			return err
		}
	}
	return nil
}

func deleteEdgeKvs(tx *gorm.DB, edgeId string) error {
	return tx.Where("edge_id = ?", edgeId).Delete(&EdgeKV{}).Error
}
