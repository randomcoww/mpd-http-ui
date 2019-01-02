package api_server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type response struct {
	Message string
}

//
// http handle funcs
//
func healthCheck(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Healthcheck")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{"ok"})
}

func search(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	logrus.Infof("Search database %s", params)

	search, err := esClient.Search(
		params["query"],
		parseNum(params["start"]),
		parseNum(params["size"]))
	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		w.WriteHeader(http.StatusOK)

		// var	result []*json.RawMessage
		// for _, hits := range search.Hits.Hits {
		// 	result = append(result, hits.Source)
		// }

		json.NewEncoder(w).Encode(search)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response{err.Error()})
	}
}

func parseNum(input string) int {
	v, err := strconv.Atoi(input)
	if err != nil {
		logrus.Errorf("Error parsing param %s: %s", input, err)
		v = -1
	}

	return v
}
