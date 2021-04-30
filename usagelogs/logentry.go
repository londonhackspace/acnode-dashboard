package usagelogs

import "time"

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`

	Card string `json:"user_card"`
	Name string `json:"user_name"`

	Node string `json:"node_mqttname"`
	Success bool `json:"success"`
}
