package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"math"
	"net/http"
	"strconv"
)

func (api *Api) handleAccessLogs(w http.ResponseWriter, r *http.Request) {
	if !api.checkAuthAdmin(w, r) {
		return
	}

	if api.usageLogger != nil {
		data,_ := json.Marshal(api.usageLogger.GetUsageLogNodes())
		w.Write(data)
	} else {
		w.Write([]byte("[]"))
	}
}

func (api *Api) handleAccessLogsNode(w http.ResponseWriter, r *http.Request) {
	if !api.checkAuthAdmin(w, r) {
		return
	}

	node := mux.Vars(r)["node"]
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	entriesPerPage := int64(50)

	entryCount := api.usageLogger.GetUsageLogCountForNode(node)

	entries := api.usageLogger.GetUsageLogsForNode(node, int64(page-1)*entriesPerPage, int64(page)*entriesPerPage)

	entriesOut := []apitypes.AccessLogEntry{}

	for _,i := range entries {
		e := apitypes.AccessLogEntry{
			Timestamp: i.Timestamp.Unix(),
			UserName: i.Name,
			UserCard: i.Card,
			Success: i.Success,
		}

		entriesOut = append(entriesOut, e)
	}

	result := apitypes.AccessLogsResponse{
		Count: entryCount,
		Page: page,
		PageCount: int64(math.Ceil(float64(entryCount)/float64(entriesPerPage))),
		LogEntries: entriesOut,
	}

	data,_ := json.Marshal(result)
	w.Write(data)
}