package acnode

type MemoryNodePersistence struct {
	nodes []ACNodeRec
}

func CreateMemoryNodePersistence() NodePersistence {
	return &MemoryNodePersistence{}
}

func (np *MemoryNodePersistence) GetNodeByMqttName(name string) (*ACNodeRec,error) {
	for i,_ := range np.nodes {
		if np.nodes[i].MqttName == name {
			return &np.nodes[i],nil
		}
	}
	return nil, NodeNotFound
}

func (np *MemoryNodePersistence) StoreNode(node *ACNodeRec) (*ACNodeRec, error) {
	for i,_ := range np.nodes {
		if np.nodes[i].MqttName == node.MqttName {
			// short path if the references are identical
			if &np.nodes[i] == node {
				return node, nil
			}
			np.nodes[i] = *node
			return node, nil
		}
	}

	np.nodes = append(np.nodes, *node)
	return node, nil
}

func (np *MemoryNodePersistence) GetAllNodes() ([]ACNodeRec, error) {
	return np.nodes, nil
}