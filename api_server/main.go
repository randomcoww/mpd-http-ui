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
	"strconv"
)

var (
	listenUrl = flag.String("listenurl", "localhost:3000", "Listen URL")
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

	r.HandleFunc("/playlist", querytPlaylistItems).
		Queries("start", "{start}").
		Queries("end", "{end}").
		Methods("GET")

	r.HandleFunc("/move", movePlaylistItems).
		Queries("start", "{start}").
		Queries("end", "{end}").
		Queries("moveto", "{moveto}").
		Methods("GET")

	r.HandleFunc("/add", addToPlaylist).
		Queries("path", "{path}").
		Queries("addto", "{addto}").
		Methods("GET")

	r.HandleFunc("/delete", deletePlaylistItems).
		Queries("start", "{start}").
		Queries("end", "{end}").
		Methods("GET")

	r.HandleFunc("/clear", clearPlaylist).Methods("GET")

	r.HandleFunc("/search", search).
		Queries("query", "{query}").
		Methods("GET")

	mpdClient = NewMpdClient("tcp", *mpdUrl)
	esClient = NewEsClient(*esUrl, esIndex, esDocument)

	<-mpdClient.Ready
	<-esClient.Ready

	fmt.Printf("start\n")
	log.Fatal(http.ListenAndServe(*listenUrl, r))
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get healthcheck\n")
	json.NewEncoder(w).Encode("ok")
}


func querytPlaylistItems(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get playlist\n")

	params := mux.Vars(r)
	start, err := strconv.Atoi(params["start"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		start = -1
	}

	end, err := strconv.Atoi(params["end"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		end = -1
	}

	attrs, err := mpdClient.QueryPlaylistItems(start, end)

	fmt.Printf("get playlist %s %d %d\n", attrs, start, end)

	if err == nil {
		json.NewEncoder(w).Encode(attrs)
	} else {
		json.NewEncoder(w).Encode(err)
	}
}


func movePlaylistItems(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("move playlist\n")

	params := mux.Vars(r)
	start, err := strconv.Atoi(params["start"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		start = -1
	}

	end, err := strconv.Atoi(params["end"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		end = -1
	}

	moveTo, err := strconv.Atoi(params["moveto"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		moveTo = -1
	}

	err = mpdClient.MovePlaylistItems(start, end, moveTo)

	if err == nil {
		json.NewEncoder(w).Encode("ok")
	} else {
		json.NewEncoder(w).Encode(err)
	}
}


func addToPlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("add to playlist\n")

	params := mux.Vars(r)
	path := params["path"]

	addTo, err := strconv.Atoi(params["addto"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		addTo = -1
	}

	id, err := mpdClient.AddToPlaylist(path, addTo)

	if err == nil {
		json.NewEncoder(w).Encode(strconv.Itoa(id))
	} else {
		json.NewEncoder(w).Encode(err)
	}
}


func deletePlaylistItems(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("delete playlist items\n")

	params := mux.Vars(r)
	start, err := strconv.Atoi(params["start"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		start = -1
	}

	end, err := strconv.Atoi(params["end"])
	if err != nil {
		fmt.Printf("err, %s\n", err)
		end = -1
	}

	err = mpdClient.DeletePlaylistItems(start, end)

	if err == nil {
		json.NewEncoder(w).Encode("ok")
	} else {
		json.NewEncoder(w).Encode(err)
	}
}


func clearPlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("clear playlist items\n")

	err := mpdClient.ClearPlaylist()

	if err == nil {
		json.NewEncoder(w).Encode("ok")
	} else {
		json.NewEncoder(w).Encode(err)
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
		json.NewEncoder(w).Encode(err)
		fmt.Printf("err %s\n", err)
	}
}
