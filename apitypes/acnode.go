package apitypes

type ACNode struct {
	Id int `json:"id"`
	Name string `json:"name"`
	MqttName string `json:"mqttName",omitempty`
	Type string `json:"nodeType"`

	// What was the timestamp this node was last seen at?
	LastSeen int64 `json:"LastSeen"`
	LastSeenMQTT int64 `json:"LastSeenMQTT"`
	LastSeenAPI int64 `json:"LastSeenAPI"`

	// When did we last see a START message for it?
	LastStarted int64 `json:"LastStarted"`

	MemFree int `json:"MemFree",omitempty`
	MemUsed int `json:"MemUsed",omitempty"`
	StatusMessage string `json:"Status", omitempty`
	Version string `json:"Version"",omitempty`

	SettingsVersion int `json:"SettingsVersion,omitempty"`
	EEPROMSettingsVersion int `json:"EEPROMSettingsVersion,omitempty"`
	ResetCause string `json:"ResetCause,omitempty"`
}
