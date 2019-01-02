// https://pragmacoders.com/building-a-json-api-in-golang/
// https://gowebexamples.com/routes-using-gorilla-mux/
// https://hakaselogs.me/2017-06-23/rest-api-with-golang

package api_server

import (
	"flag"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	es "github.com/randomcoww/go-mpd-es/pkg/es_handler"
	event "github.com/randomcoww/go-mpd-es/pkg/mpd_event"
	mpd "github.com/randomcoww/go-mpd-es/pkg/mpd_handler"
	"github.com/sirupsen/logrus"
)

func Main() {
	listenURL := flag.String("listenurl", "localhost:3000", "Listen URL")
	mpdURL := flag.String("mpdurl", "localhost:6600", "MPD URL")
	esURL := flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
	flag.Parse()

	newServer(*listenURL, *mpdURL, *esURL)
}

func newServer(listenURL, mpdURL, esURL string) {
	// backend stuff
	mpdClient = mpd.NewMpdClient("tcp", mpdURL)
	esClient = es.NewEsClient(esURL, esIndex, esDocument, "")
	mpdEvent = event.NewEventWatcher("tcp", mpdURL)

	// websocket hub
	hub := newHub()

	// mux http server
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", healthCheck).
		Methods("GET")

	r.HandleFunc("/database/search", search).
		Queries("q", "{query}").
		Queries("start", "{start}").
		Queries("size", "{size}").
		Methods("GET")

	// websocket handler
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	// serve http
	logrus.Infof("API server start on %s", listenURL)
	err := http.ListenAndServe(listenURL, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r))
	if err != nil {
		logrus.Errorf("%v", err)
	}

	// wait backend start
	<-mpdClient.Ready
	<-esClient.Ready

	go hub.run()
	go hub.eventBroadcaster()

	// playlistVersion, playlistLength, _ = getPlaylistStatus()
	updatePlaylistStatus()
	// set mpd repeat by default
	mpdClient.Conn.Repeat(true)
}
