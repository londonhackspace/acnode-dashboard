package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type APIClientImpl struct {
	apikey string
	server string
}

func MakeAPIClient(apikey string, server string) APIClient {
	return &APIClientImpl{
		apikey: apikey,
		server: server,
	}
}

func (apiclient *APIClientImpl) makeGetRequest(endpoint string) (*http.Request,error) {
	req, err := http.NewRequest(http.MethodGet, apiclient.server + endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("API-KEY", apiclient.apikey)

	return req, nil
}

func (apiclient *APIClientImpl) makePostRequest(endpoint string, body []byte) (*http.Request,error) {
	req, err := http.NewRequest(http.MethodPost, apiclient.server + endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("API-KEY", apiclient.apikey)

	return req, nil
}

func (apiclient *APIClientImpl) GetNodes() ([]string, error) {
	req, err := apiclient.makeGetRequest("nodes")
	if err != nil {
		log.Err(err).Msg("Error creating Get Nodes request")
		return []string{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Err(err).Msg("Error getting nodes list")
		return []string{}, err
	}

	if resp.StatusCode != 200 {
		err = errors.New("bad response to GET request")
		log.Err(err).Int("responseCode", resp.StatusCode).Send()
		return []string{}, err
	}

	var res []string

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("Error getting nodes list")
		return []string{}, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Err(err).Msg("Error getting nodes list")
		return []string{}, err
	}

	// filter some badly named nodes out
	var filtered []string

	for _,n := range res {
		if len(n) == 0 {
			continue
		}
		if n[0] == '/' {
			continue
		}
		filtered = append(filtered, n)
	}

	return filtered,nil
}

func (apiclient *APIClientImpl) GetNode(name string) (*apitypes.ACNode, error) {
	if len(name) == 0 {
		return nil, errors.New("Invalid node name")
	}

	if name[0] == '/' {
		return nil, errors.New("Invalid node name")
	}

	req, err := apiclient.makeGetRequest("nodes/"+name)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New("unexpected return code")
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var node apitypes.ACNode
	err = json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}