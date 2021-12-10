package slackapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type responseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type Impl struct {
	token string

	channelMap map[string]string
}

func CreateSlackAPI(token string) SlackAPI {
	return &Impl{
		token:      token,
		channelMap: map[string]string{},
	}
}

func (slack *Impl) makeGetRequest(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "https://slack.com/api/"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+slack.token)

	return req, nil
}

func (slack *Impl) makePostRequest(endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/"+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+slack.token)

	return req, nil
}

type getChannelsResponse struct {
	Ok       bool             `json:"ok"`
	Channels []Channel        `json:"channels"`
	Extra    responseMetadata `json:"response_metadata"`
}

func (slack *Impl) getChannels(cursor string) (getChannelsResponse, error) {
	req, err := slack.makeGetRequest("conversations.list")
	if err != nil {
		return getChannelsResponse{}, err
	}

	if len(cursor) > 0 {
		req.URL.RawQuery = req.URL.RawQuery + "cursor=" + cursor
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return getChannelsResponse{}, err
	}

	if resp.StatusCode > 299 {
		err = errors.New("Unexpected status code")
		return getChannelsResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	var res getChannelsResponse
	json.Unmarshal(body, &res)

	return res, nil
}

func (slack *Impl) GetChannels() ([]Channel, error) {
	var channels []Channel
	cursor := ""
	for {
		res, err := slack.getChannels(cursor)
		if err != nil {
			return []Channel{}, err
		}

		channels = append(channels, res.Channels...)

		cursor = res.Extra.NextCursor
		if len(cursor) == 0 {
			break
		}
	}

	return channels, nil
}

func (slack *Impl) getChannelKey(channel string) string {
	key, ok := slack.channelMap[channel]
	if ok {
		return key
	}

	channels, err := slack.GetChannels()
	if err != nil {
		log.Err(err).Msg("Error looking up channel keys")
		return ""
	}

	for _, c := range channels {
		slack.channelMap[c.Name] = c.Id
	}

	key, ok = slack.channelMap[channel]
	if ok {
		return key
	}

	log.Error().Str("channelName", channel).Msg("Unknown slack channel")
	return ""
}

type postMessageBody struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (slack *Impl) PostMessage(message string, channel string) error {
	key := slack.getChannelKey(channel)
	if len(key) == 0 {
		err := errors.New("unknown channel")
		log.Err(err).Msg("Empty Channel Key")
		return err
	}

	bodyStruct := postMessageBody{
		Channel: key,
		Text:    message,
	}
	data, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	req, err := slack.makePostRequest("chat.postMessage", data)
	if err != nil {
		return err
	}

	// We're going to be posting a JSON block
	req.Header.Add("Content-type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
