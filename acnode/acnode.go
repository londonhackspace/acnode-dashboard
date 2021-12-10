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

func NodeTypeIsTool(t int) bool {
	return t == NodeTypeTool || t == NodeTypePrinter
}

type ACNodeRec struct {
	mtx sync.Mutex

	updateTrigger acnodeUpdateTrigger

	Id        int
	Name      string
	MqttName  string
	NodeType  int
	InService bool
	InUse     bool

	// last known status
	LastSeen     time.Time // old LastSeen field
	LastSeenMQTT time.Time
	LastSeenAPI  time.Time

	LastStarted   time.Time
	MemFree       int
	MemUsed       int
	StatusMessage string
	Version       string

	SettingsVersion       int
	EEPROMSettingsVersion int
	ResetCause            string

	Transient bool

	PrinterStatus *PrinterStatus
}

type ACNode interface {
	GetId() int
	SetId(id int)
	GetType() int
	GetInService() bool
	SetInService(inService bool)
	GetName() string
	SetName(name string)
	GetMqttName() string
	SetType(t int)
	SetMemoryStats(free int, used int)
	SetVersion(ver string)
	GetLastSeen() time.Time
	GetLastSeenAPI() time.Time
	SetLastSeenAPI(t time.Time)
	GetLastSeenMQTT() time.Time
	SetLastSeenMQTT(t time.Time)
	GetLastStarted() time.Time
	SetLastStarted(t time.Time)
	SetStatusMessage(m string)
	GetAPIRecord() apitypes.ACNode
	SetSettingsVersion(ver int)
	SetEepromSettingsVersion(ver int)
	SetResetCause(rstc string)
	GetPrinterStatus() *PrinterStatus
	GetIsTransient() bool
	SetIsTransient(transient bool)
	GetInUse() bool
	SetInUse(inuse bool)
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
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) GetType() int {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.NodeType
}

func (node *ACNodeRec) GetInService() bool {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.InService
}

func (node *ACNodeRec) SetInService(inService bool) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.InService = inService
	node.updateTrigger.OnNodeUpdate(node)
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
	node.updateTrigger.OnNodeUpdate(node)
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
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) SetMemoryStats(free int, used int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.MemFree = free
	node.MemUsed = used
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) SetVersion(ver string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.Version = ver
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) GetLastSeen() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if node.LastSeen.After(node.LastSeenAPI) && node.LastSeen.After(node.LastSeenMQTT) {
		return node.LastSeen
	}

	if node.LastSeenAPI.After(node.LastSeenMQTT) {
		return node.LastSeenAPI
	}

	return node.LastSeenMQTT
}

func (node *ACNodeRec) GetLastSeenAPI() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.LastSeenAPI
}

func (node *ACNodeRec) SetLastSeenAPI(t time.Time) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if t.After(node.LastSeenAPI) {
		node.LastSeenAPI = t
		node.updateTrigger.OnNodeUpdate(node)
	}
}

func (node *ACNodeRec) GetLastSeenMQTT() time.Time {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.LastSeenMQTT
}
func (node *ACNodeRec) SetLastSeenMQTT(t time.Time) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if t.After(node.LastSeenMQTT) {
		node.LastSeenMQTT = t
		node.updateTrigger.OnNodeUpdate(node)
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
		node.updateTrigger.OnNodeUpdate(node)
	}
}

func (node *ACNodeRec) SetStatusMessage(m string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.StatusMessage = m
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) SetSettingsVersion(ver int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.SettingsVersion = ver
	node.updateTrigger.OnNodeUpdate(node)
}
func (node *ACNodeRec) SetEepromSettingsVersion(ver int) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.EEPROMSettingsVersion = ver
	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) SetResetCause(rstc string) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.ResetCause = rstc
	node.updateTrigger.OnNodeUpdate(node)
}

func makeApiTimestamp(t time.Time) int64 {
	if t.IsZero() {
		return -1
	}
	return t.Unix()
}

func (node *ACNodeRec) GetPrinterStatus() *PrinterStatus {
	if node.NodeType != NodeTypePrinter {
		return nil
	}
	if node.PrinterStatus == nil {
		node.PrinterStatus = GetDefaultPrinterStatus()
	}
	return node.PrinterStatus
}

func (node *ACNodeRec) GetIsTransient() bool {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.Transient
}

func (node *ACNodeRec) SetIsTransient(transient bool) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	node.Transient = transient

	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) GetInUse() bool {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	return node.InUse
}

func (node *ACNodeRec) SetInUse(inuse bool) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if inuse == node.InUse {
		return
	}
	node.InUse = inuse

	node.updateTrigger.OnNodeUpdate(node)
}

func (node *ACNodeRec) GetAPIRecord() apitypes.ACNode {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	// figure out which LastSeen value to include
	lastSeen := makeApiTimestamp(node.LastSeen)
	if node.LastSeenMQTT.After(node.LastSeenAPI) && node.LastSeenMQTT.After(node.LastSeen) {
		lastSeen = makeApiTimestamp(node.LastSeenMQTT)
	} else if node.LastSeenAPI.After(node.LastSeenMQTT) && node.LastSeenAPI.After(node.LastSeen) {
		lastSeen = makeApiTimestamp(node.LastSeenAPI)
	}

	var printerStatus *apitypes.PrinterStatus = nil

	if node.NodeType == NodeTypePrinter && node.PrinterStatus != nil {
		printerStatus = &apitypes.PrinterStatus{
			MqttConnected:      node.PrinterStatus.MqttConneced,
			OctoprintConnected: node.PrinterStatus.OctoprintConnected,
			FirmwareVersion:    node.PrinterStatus.FirmwareVersion,
			ZHeight:            node.PrinterStatus.ZHeight,
			PiUndervoltage:     node.PrinterStatus.PiUndervoltage,
			PiOverheat:         node.PrinterStatus.PiOverheat,
			HotendTemperature:  node.PrinterStatus.HotendTemperature,
			BedTemperature:     node.PrinterStatus.BedTemperature,
		}
	}

	return apitypes.ACNode{
		Id:                    node.Id,
		Name:                  node.Name,
		MqttName:              node.MqttName,
		Type:                  NodeTypeToString(node.NodeType),
		InService:             node.InService,
		LastSeen:              lastSeen,
		LastSeenAPI:           makeApiTimestamp(node.LastSeenAPI),
		LastSeenMQTT:          makeApiTimestamp(node.LastSeenMQTT),
		LastStarted:           makeApiTimestamp(node.LastStarted),
		MemFree:               node.MemFree,
		MemUsed:               node.MemUsed,
		StatusMessage:         node.StatusMessage,
		Version:               node.Version,
		SettingsVersion:       node.SettingsVersion,
		EEPROMSettingsVersion: node.EEPROMSettingsVersion,
		ResetCause:            node.ResetCause,
		PrinterStatus:         printerStatus,
		IsTransient:           node.Transient,
		InUse:                 node.InUse,
	}
}
