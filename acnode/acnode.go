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

type ACNodeRec struct {
	mtx sync.Mutex

	Id       int
	Name     string
	MqttName string
	NodeType int

	// last known status
	LastSeen time.Time
	LastStarted time.Time
	MemFree int
	MemUsed int
	StatusMessage string
	Version string

	SettingsVersion int
	EEPROMSettingsVersion int
	ResetCause string
}

type ACNode interface {
	GetId() int
	SetId(id int)
	GetType() int
	GetName() string
	SetName(name string)
	GetMqttName() string
	SetType(t int)
	SetMemoryStats(free int, used int)
	SetVersion(ver string)
	GetLastSeen() time.Time
	SetLastSeen(t time.Time)
	GetLastStarted() time.Time
	SetLastStarted(t time.Time)
	SetStatusMessage(m string)
	GetAPIRecord() apitypes.ACNode
	SetSettingsVersion(ver int)
	SetEepromSettingsVersion(ver int)
	SetResetCause(rstc string)
}

func (node *ACNodeRec) GetId() int {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.Id
}

func (node *ACNodeRec) SetId(id int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.Id = id
}

func (node *ACNodeRec) GetType() int {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.NodeType
}

func (node *ACNodeRec) GetName() string {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.Name
}

func (node *ACNodeRec) SetName(name string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.Name = name
}

func (node *ACNodeRec) GetMqttName() string {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.MqttName
}

func (node *ACNodeRec) SetType(t int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.NodeType = t
}

func (node *ACNodeRec) SetMemoryStats(free int, used int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.MemFree = free
	node.MemUsed = used
}

func (node *ACNodeRec) SetVersion(ver string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.Version = ver
}

func (node *ACNodeRec) GetLastSeen() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.LastSeen
}

func (node *ACNodeRec) SetLastSeen(t time.Time) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if t.After(node.LastSeen) {
		node.LastSeen = t
	}
}

func (node *ACNodeRec) GetLastStarted() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.LastStarted
}

func (node *ACNodeRec) SetLastStarted(t time.Time) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if t.After(node.LastStarted) {
		node.LastStarted = t
	}
}

func (node *ACNodeRec) SetStatusMessage(m string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.StatusMessage = m
}

func (node *ACNodeRec) SetSettingsVersion(ver int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.SettingsVersion = ver
}
func (node *ACNodeRec) SetEepromSettingsVersion(ver int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.EEPROMSettingsVersion = ver
}

func (node *ACNodeRec) SetResetCause(rstc string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.ResetCause = rstc
}

func (node *ACNodeRec) GetAPIRecord() apitypes.ACNode {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	// if we've never seen it, return -1 in this field
	lastSeen := int(time.Now().Sub(node.LastSeen).Seconds())
	if node.LastSeen.IsZero() {
		lastSeen = -1
	}

	lastStarted := int(time.Now().Sub(node.LastStarted).Seconds())
	if node.LastStarted.IsZero() {
		lastStarted = -1
	}

	return apitypes.ACNode{
		Id:            node.Id,
		Name:          node.Name,
		MqttName: 	   node.MqttName,
		Type:          NodeTypeToString(node.NodeType),
		LastSeen:      lastSeen,
		LastStarted:   lastStarted,
		MemFree:       node.MemFree,
		MemUsed:       node.MemUsed,
		StatusMessage: node.StatusMessage,
		Version:       node.Version,
		SettingsVersion: node.SettingsVersion,
		EEPROMSettingsVersion: node.EEPROMSettingsVersion,
		ResetCause: node.ResetCause,
	}
}