package apitypes

type ACNode struct {
	Id int `json:"id"`
	Name string `json:"name"`
	MqttName string `json:"mqttName",omitempty`
	Type string `json:"nodeType"`

	// how many seconds ago was the node last seen?
	LastSeen int `json:"LastSeen"`

	MemFree int `json:"MemFree",omitempty`
	MemUsed int `json:"MemUsed",omitempty"`
	StatusMessage string `json:"Status", omitempty`
	Version string `json:"Version"",omitempty`
}
