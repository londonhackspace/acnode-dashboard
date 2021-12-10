package slackapi

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SlackAPI interface {
	GetChannels() ([]Channel, error)
	PostMessage(message string, channel string) error
}
