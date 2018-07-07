// https://pragmacoders.com/building-a-json-api-in-golang/
// https://gowebexamples.com/routes-using-gorilla-mux/
// https://hakaselogs.me/2017-06-23/rest-api-with-golang

package main

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"flag"
)

var (
	mpdUrl = flag.String("mpdurl", "localhost:6600", "MPD URL")
	esUrl = flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
	esIndex, esDocument = "songs", "song"
	mpdClient *MpdClient
	esClient *EsClient
)


func main() {
	fmt.Printf("Start")
	flag.Parse()

  r := mux.NewRouter()

	r.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	r.HandleFunc("/playlist", querytPlaylistItems).Methods("GET")
  r.HandleFunc("/search/{query}", search).Methods("GET")

	mpdClient = NewMpdClient("tcp", *mpdUrl)
	esClient = NewEsClient(*esUrl, esIndex, esDocument)

	<-mpdClient.Ready
	<-esClient.Ready

	log.Fatal(http.ListenAndServe(":3000", r))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get healthcheck\n")
	json.NewEncoder(w).Encode("ok")
}

func querytPlaylistItems(w http.ResponseWriter, r *http.Request) {
	attrs, err := mpdClient.QueryPlaylistItems(-1, -1)

	fmt.Printf("get playlist %s\n", attrs)

	if err == nil {
		json.NewEncoder(w).Encode(attrs)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := params["query"]

	searchResult, err := esClient.Search(query)

	if err == nil {
		fmt.Printf("searchresult %s\n", searchResult)

		json.NewEncoder(w).Encode(searchResult.Hits.Hits)
	} else {
		fmt.Printf("err %s\n", err)
	}
}
