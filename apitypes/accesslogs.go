package apitypes

type AccessLogEntry struct {
	Timestamp int64 `json:"timestamp"`

	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
	UserCard string `json:"user_card"`

	Success bool `json:"success"`
}

type AccessLogsResponse struct {
	Count int64 `json:"count"`
	Page int `json:"page"`
	PageCount int64 `json:"pageCount"`

	LogEntries []AccessLogEntry `json:"entries"`
}
