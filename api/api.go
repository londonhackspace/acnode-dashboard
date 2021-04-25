package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

func (api *Api) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	ok, _ := auth.CheckAuthAPI(w, r)

	if !ok {
		w.WriteHeader(401)
	}

	return ok
}

func (api *Api) handleNodes(w http.ResponseWriter, r *http.Request) {
	if ! api.checkAuth(w, r) {
		return
	}
	nodes := api.acnodeHandler.GetNodes()

	apinodes := []string{}

	for _, n := range nodes {
		apinodes = append(apinodes, n.GetMqttName())
	}

	data,_ := json.Marshal(&apinodes)
	w.Write(data)
}

func (api *Api) handleNodeEntry(w http.ResponseWriter, r *http.Request) {
	if ! api.checkAuth(w, r) {
		return
	}

	nodes := api.acnodeHandler.GetNodes()

	for _, node := range nodes {
		if node.GetMqttName() == mux.Vars(r)["nodeName"] {
			noderec := node.GetAPIRecord()
			data,_ := json.Marshal(&noderec)
			w.Write(data)
			return
		}
	}

	w.WriteHeader(404)
}

type SetStatusBody struct {
	Version string `json:"version",omitempty`
	Timestamp int64 `json:"timestamp",omitempty`
}

func (api *Api) handleSetStatus(w http.ResponseWriter, r *http.Request) {
	if ! api.checkAuth(w, r) {
		return
	}

	id,err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		// not an ID -> not found
		w.WriteHeader(404)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var status SetStatusBody

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Err(err).Msg("Error unmarshalling request")
	}

	nodes := api.acnodeHandler.GetNodes()

	for i, _ := range nodes {
		n := &nodes[i]
		if id != n.GetId() {
			continue
		}

		if len(status.Version) != 0 {
			n.SetVersion(status.Version)
		}

		if status.Timestamp > 0 {
			d := time.Unix(status.Timestamp, 0)
			n.SetLastSeen(d)
		}

		var s = string(body)
		log.Info().Str("raw",s).Msg("")
		log.Info().Str("Node", n.GetMqttName()).
			Msg("Got update for node")

		w.WriteHeader(204)
		return
	}

	// node doesn't exist
	w.WriteHeader(404)
}

func (api *Api) GetRouter() http.Handler {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/nodes", api.handleNodes)
	rtr.HandleFunc("/nodes/{nodeName}", api.handleNodeEntry)
	rtr.HandleFunc("/nodes/setStatus/{id}", api.handleSetStatus)

	return rtr
}