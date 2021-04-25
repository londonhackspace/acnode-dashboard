package acnode

import (
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"sync"
	"time"
)

const (
	NodeTypeTool    = iota
	NodeTypeDoor    = iota
	NodeTypePrinter = iota
)

func NodeTypeToString(t int) string {
	switch t {
	case NodeTypeTool:
		return "Tool"
	case NodeTypeDoor:
		return "Door"
	case NodeTypePrinter:
		return "Printer"
	}

	return "Unknown"
}

type ACNode struct {
	mtx sync.Mutex

	id       int
	name     string
	mqttName string
	nodeType int

	// last known status
	lastSeen time.Time
	memFree int
	memUsed int
	statusMessage string
	version string
}

func (node *ACNode) GetId() int {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.id
}

func (node *ACNode) SetId(id int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.id = id
}

func (node *ACNode) GetType() int {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.nodeType
}

func (node *ACNode) GetName() string {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.name
}

func (node *ACNode) SetName(name string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.name = name
}

func (node *ACNode) GetMqttName() string {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.mqttName
}

func (node *ACNode) SetType(t int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.nodeType = t
}

func (node *ACNode) SetMemoryStats(free int, used int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.memFree = free
	node.memUsed = used
}

func (node *ACNode) SetVersion(ver string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.version = ver
}

func (node *ACNode) GetLastSeen() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.lastSeen
}

func (node *ACNode) SetLastSeen(t time.Time) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if t.After(node.lastSeen) {
		node.lastSeen = t
	}
}

func (node *ACNode) SetStatusMessage(m string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.statusMessage = m
}

func (node *ACNode) GetAPIRecord() apitypes.ACNode {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	// if we've never seen it, return -1 in this field
	lastSeen := int(time.Now().Sub(node.lastSeen).Seconds())
	if node.lastSeen.IsZero() {
		lastSeen = -1
	}

	return apitypes.ACNode{
		Id:            node.id,
		Name:          node.name,
		MqttName: 	   node.mqttName,
		Type:          NodeTypeToString(node.nodeType),
		LastSeen:      lastSeen,
		MemFree:       node.memFree,
		MemUsed:       node.memUsed,
		StatusMessage: node.statusMessage,
		Version:       node.version,
	}
}