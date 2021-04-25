package acserver_api

type ToolStatusResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	StatusMessage string `json:"status_message"`
	InUse string `json:"in_use"`
	Type string `json:"type"`
	MqttName string `json:"mqtt_name"`
}

type UserToolSummaryResponse struct {
	Name string `json:"name"`
	Status string `json:"status"`
	StatusMessage string `json:"status_message"`
	InUse string `json:"in_use"`
	Permission string `json:"permission"`
	Type string `json:"type"`
}

type UserCardResponse struct {
	UserName string `json:"user_name"`
	GlaDOSFile string `json:"gladosfile"`
	Subscribed bool `json:"subscribed"`
}
