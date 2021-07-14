package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/londonhackspace/acnode-dashboard/usagelogs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_requests_served",
		Help: "API Request counter",
	}, []string{"endpoint","method"})
	inflightCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "api_inflight_requests",
		Help: "Currently in-flight API requests",
	})
	requestTimer = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "api_processing_time",
		Help: "Time taken to process requests, in milliseconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path"})
)

type Api struct {
	conf *config.Config

	acnodeHandler *acnode.ACNodeHandler
	usageLogger usagelogs.UsageLogger
}

func CreateApi(conf *config.Config, acnodeHandler *acnode.ACNodeHandler, usageLogger usagelogs.UsageLogger) Api {
	return Api{
		conf: conf,
		acnodeHandler: acnodeHandler,
		usageLogger: usageLogger,
	}
}

func (api *Api) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	ok, _ := auth.CheckAuthAPI(w, r)

	if !ok {
		w.WriteHeader(401)
	}

	return ok
}

func (api *Api) checkAuthAdmin(w http.ResponseWriter, r *http.Request) bool {
	ok, user := auth.CheckAuthAPI(w, r)

	if !ok || !user.IsAdmin(api.conf) {
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

	var status apitypes.SetStatusBody

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Err(err).Msg("Error unmarshalling request")
	}

	nodes := api.acnodeHandler.GetNodes()

	for i, _ := range nodes {
		n := nodes[i]
		if id != n.GetId() {
			continue
		}

		if len(status.Version) != 0 {
			n.SetVersion(status.Version)
		}

		if status.Timestamp > 0 {
			d := time.Unix(status.Timestamp, 0)
			n.SetLastSeenAPI(d)
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

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginError struct {
	Error string `json:"error"`
}

func (api *Api) handleLogin(w http.ResponseWriter, r *http.Request) {
	// first see if the user is already logged in
	if ok,_ := auth.CheckAuthAPI(w, r); ok {
		w.WriteHeader(204)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	payload := loginBody{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
	}
	if auth.AuthenticateUser(w, payload.Username, payload.Password) {
		w.WriteHeader(204)
		return
	}

	resp := loginError{Error: "Bad Credentials"}
	data,_ := json.Marshal(&resp)
	w.WriteHeader(401)
	w.Write(data)
}

func (api *Api) handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r)
	w.WriteHeader(204)
}

func (api *Api) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
	ok, user := auth.CheckAuthUser(w, r)

	if !ok {
		w.WriteHeader(401)
		return
	}

	ret := apitypes.User{
		Username: user.UserName,
		Name: user.Name,
		Admin: user.IsAdmin(api.conf),
	}

	data,_ := json.Marshal(ret)

	w.Write(data)
}

type promInterceptor struct {
	next http.Handler
}

func (i promInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCounter.WithLabelValues(r.URL.Path, r.Method).Inc()
	inflightCounter.Inc()
	start := time.Now()
	i.next.ServeHTTP(w, r)
	end := time.Since(start)
	requestTimer.WithLabelValues(r.URL.Path).Observe(float64(end.Nanoseconds())/1000000)
	inflightCounter.Dec()
}

func (api *Api) GetRouter() http.Handler {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/nodes", api.handleNodes)
	rtr.HandleFunc("/nodes/{nodeName}", api.handleNodeEntry)
	rtr.HandleFunc("/nodes/setStatus/{id}", api.handleSetStatus)

	rtr.HandleFunc("/accesslogs", api.handleAccessLogs)
	rtr.HandleFunc("/accesslogs/{node}", api.handleAccessLogsNode)

	rtr.Methods("POST").Path("/auth/login").HandlerFunc(api.handleLogin)
	rtr.HandleFunc("/auth/logout", api.handleLogout)
	rtr.HandleFunc("/auth/currentuser", api.handleCurrentUser)
	return promInterceptor{next: rtr}
}