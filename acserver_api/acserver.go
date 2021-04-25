package acserver_api

import (
	"encoding/json"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type ACServer struct {
	config *config.Config
}

func CreateACServer(config *config.Config) ACServer {
	return ACServer{
		config: config,
	}
}

func (svr *ACServer) makeRequest(path string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", svr.config.AcserverUrl + path, nil)
	if err != nil {
		log.Err(err).Str("path", path).Msg("Error making request")
		return []byte{}
	}
	req.Header.Add("API-KEY", svr.config.AcserverApiKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Err(err).Str("path", path).Msg("Error requesting data")
		return []byte{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Str("path", path).Msg("Error reading from server")
		return []byte{}
	}
	return body
}

func (svr *ACServer) GetTools() []ToolStatusResponse {
	response := make([]ToolStatusResponse, 0)

	json.Unmarshal(svr.makeRequest("/api/get_tools_status"), &response)

	return response
}
