package acnode

type ACNodeHandler struct {
	nodes []ACNode
}

func CreateACNodeHandler() ACNodeHandler {
	return ACNodeHandler{}
}

func (h *ACNodeHandler) AddNode(node ACNode) {
	h.nodes = append(h.nodes, node)
}

func (h *ACNodeHandler) GetNodeByMqttName(name string) *ACNode {
	for i, _ := range h.nodes {
		if h.nodes[i].GetMqttName() == name {
			return &h.nodes[i]
		}
	}

	node := ACNode{
		nodeType: NodeTypeTool,
		name:     name,
		mqttName: name,
	}
	h.AddNode(node)

	// return a ref to the last entry we just added
	return &h.nodes[len(h.nodes)-1]
}

func (h *ACNodeHandler) GetNodes() []ACNode {
	return h.nodes
}