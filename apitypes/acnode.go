package apitypes

type ACNode struct {
	Id int `json:"id"`
	Name string `json:"name"`
	MqttName string `json:"mqttName",omitempty`
	Type string `json:"nodeType"`

	// how many seconds ago was the node last seen?
	LastSeen int `json:"LastSeen"`

	// When did we last see a START message for it?
	LastStarted int `json:"LastStarted"`

	MemFree int `json:"MemFree",omitempty`
	MemUsed int `json:"MemUsed",omitempty"`
	StatusMessage string `json:"Status", omitempty`
	Version string `json:"Version"",omitempty`

	SettingsVersion int `json:"SettingsVersion,omitempty"`
	EEPROMSettingsVersion int `json:"EEPROMSettingsVersion,omitempty"`
	ResetCause string `json:"ResetCause,omitempty"`
}
