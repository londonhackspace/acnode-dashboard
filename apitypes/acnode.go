package apitypes

type PrinterStatus struct {
	MqttConnected bool `json:"mqttConnected"`
	OctoprintConnected bool `json:"octoprintConnected"`
	FirmwareVersion string `json:"firmwareVersion"`
	ZHeight float32 `json:"zHeight"`
	PiUndervoltage bool `json:"piUndervoltage"`
	PiOverheat bool `json:"piOverheat"`
	HotendTemperature float32 `json:"hotendTemperature"`
	BedTemperature float32 `json:"bedTemperature"`
}

type ACNode struct {
	Id int `json:"id"`
	Name string `json:"name"`
	MqttName string `json:"mqttName",omitempty`
	Type string `json:"nodeType"`
	InService bool `json:"inService"`
	InUse bool `json:"InUse"`

	// What was the timestamp this node was last seen at?
	LastSeen int64 `json:"LastSeen"`
	LastSeenMQTT int64 `json:"LastSeenMQTT"`
	LastSeenAPI int64 `json:"LastSeenAPI"`

	// When did we last see a START message for it?
	LastStarted int64 `json:"LastStarted"`

	MemFree int `json:"MemFree,omitempty"`
	MemUsed int `json:"MemUsed,omitempty""`
	StatusMessage string `json:"Status,omitempty"`
	Version string `json:"Version,omitempty"`

	SettingsVersion int `json:"SettingsVersion,omitempty"`
	EEPROMSettingsVersion int `json:"EEPROMSettingsVersion,omitempty"`
	ResetCause string `json:"ResetCause,omitempty"`

	CameraId *int `json:"CameraId,omitempty"`
	IsTransient bool `json:"IsTransient"`

	PrinterStatus *PrinterStatus `json:"PrinterStatus,omitempty"`
}

type NodeProps struct {
	CameraId *int `json:"CameraId,omitempty"`
	IsTransient *bool `json:"isTransient,omitempty"`
}