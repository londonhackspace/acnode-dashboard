package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/config"
	"net/http"
)

type Api struct {
	conf *config.Config

	acnodeHandler *acnode.ACNodeHandler
}

func CreateApi(conf *config.Config, acnodeHandler *acnode.ACNodeHandler) Api {
	return Api{
		conf: conf,
		acnodeHandler: acnodeHandler,
	}
}

func (api *Api) handleNodes(w http.ResponseWriter, r *http.Request) {
	nodes := api.acnodeHandler.GetNodes()

	apinodes := []string{}

	for _, n := range nodes {
		apinodes = append(apinodes, n.GetMqttName())
	}

	data,_ := json.Marshal(&apinodes)
	w.Write(data)
}

func (api *Api) handleNodeEntry(w http.ResponseWriter, r *http.Request) {
	nodes := api.acnodeHandler.GetNodes()

	for _, node := range nodes {
		if node.GetMqttName() == mux.Vars(r)["nodeId"] {
			noderec := node.GetAPIRecord()
			data,_ := json.Marshal(&noderec)
			w.Write(data)
			return
		}
	}

	w.WriteHeader(404)
}

func (api *Api) GetRouter() http.Handler {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/nodes", api.handleNodes)
	rtr.HandleFunc("/nodes/{nodeId}", api.handleNodeEntry)

	return rtr
}