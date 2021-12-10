package acnode

import "errors"

type NodePersistence interface {
	GetNodeByMqttName(name string) (*ACNodeRec, error)
	StoreNode(node *ACNodeRec) (*ACNodeRec, error)
	GetAllNodes() ([]ACNodeRec, error)
}

var NodeNotFound = errors.New("node not found")
