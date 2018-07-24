package main

import (
	"fmt"
	"time"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	mpd_handler "github.com/randomcoww/go-mpd-es/mpd_handler"
	mpd_event "github.com/randomcoww/go-mpd-es/mpd_event"
	es_handler "github.com/randomcoww/go-mpd-es/es_handler"
	// es_handler "local/es_handler"
)

var (
	esIndex, esDocument = "songs", "song"
	mpdClient *mpd_handler.MpdClient
	esClient *es_handler.EsClient
	mpdEvent *mpd_event.MpdEvent

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)


// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}


const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second
)


type response struct {
	Message string
}

type socketMessage struct {
	Name string `json:"mutation"`
	Data interface{} `json:"value"`
}


func NewServer(listenUrl, mpdUrl, esUrl string) {

	// backend stuff
	mpdClient = mpd_handler.NewMpdClient("tcp", mpdUrl)
	esClient = es_handler.NewEsClient(esUrl, esIndex, esDocument, "")
	mpdEvent = mpd_event.NewEventWatcher("tcp", mpdUrl)

	// websocket hub
	hub := newHub()
	go hub.run()
	go hub.eventBroadcaster()

	// mux routes
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
		serveWs(hub, w, r)
	})

	// wait backend start
	<-mpdClient.Ready
	<-esClient.Ready

  // set mpd repeat by default
	mpdClient.Conn.Repeat(true)

	// serve http
	fmt.Printf("API server start on %s\n", listenUrl)
	log.Fatal(http.ListenAndServe(listenUrl, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r)))
}


func parseNum(input string) (int) {
	v, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Error parsing param %s: %s\n", input, err)
		v = -1
	}

	return v
}


func (h *Hub) eventBroadcaster() {
	for {
		select {
		case e := <-mpdEvent.Event:
			fmt.Printf("Got MPD event: %s\n", e)
			h.broadcast <-[]byte(e)
		}
	}
}


func (c *Client) sendStatusMessage() {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return
	}

	err = c.conn.WriteJSON(&socketMessage{Data: attrs, Name: "status"})
	if err != nil {
		fmt.Printf("Failed to write status message %s\n", err)
		return
	}
}

func (c *Client) sendCurrentSongMessage() {
	attrs, err := mpdClient.Conn.CurrentSong()
	if err != nil {
		return
	}

	err = c.conn.WriteJSON(&socketMessage{Data: attrs, Name: "currentsong"})
	if err != nil {
		fmt.Printf("Failed to write currentsong message %s\n", err)
		return
	}
}

// respond to MPD playist event - return current playlist length
func (c *Client) sendPlaylistMessage() {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return
	}

	message := make([]int, 2)

	playlist, err := strconv.Atoi(attrs["playlist"])
	if err != nil {
		return
	}

	playlistlength, err := strconv.Atoi(attrs["playlistlength"])
	if err != nil {
		return
	}

	message[0] = playlist
	message[1] = playlistlength

	// p, err := mpdClient.Conn.PlaylistInfo(-1, -1)
	// if err != nil {
	// 	return
	// }
	// for i, j := range p {
	// 	fmt.Printf("%s %s %s\n", i, j)
	// }

  // send playlist and playlistlength
	err = c.conn.WriteJSON(&socketMessage{Data: message, Name: "playlist"})
	if err != nil {
		fmt.Printf("Failed to write playlist message %s\n", err)
		return
	}
}

// send playlist info
func (c *Client) sendPlaylistUpdateMessage(start, end int) {
	attrs, err := mpdClient.Conn.PlaylistInfo(start, end)
	if err != nil {
		return
	}

	err = c.conn.WriteJSON(&socketMessage{Data: attrs, Name: "playlistupdate"})
	if err != nil {
		fmt.Printf("Failed to write playlistupdate message %s\n", err)
		return
	}
}

func (c *Client) sendSeekMessage() {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return
	}

	switch attrs["state"] {
	case "play":
		message := make([]float64, 2)

		elapsed, err := strconv.ParseFloat(attrs["elapsed"], 32)
		if err != nil {
			return
		}

		duration, err := strconv.ParseFloat(attrs["duration"], 32)
		if err != nil {
			return
		}

		message[0] = elapsed
		message[1] = duration

		c.conn.WriteJSON(&socketMessage{Data: message, Name: "seek"})
	}
}

func (c *Client) sendSearchMessage(query string, start, size int) {
	search, err := esClient.Search(query, start, size)
	if err != nil {
		return
	}

	message := make([]interface{}, 2)

	var	result []*json.RawMessage
	for _, hits := range search.Hits.Hits {
		result = append(result, hits.Source)
	}

	message[0] = result
	message[1] = start

	err = c.conn.WriteJSON(&socketMessage{Data: message, Name: "search"})
	if err != nil {
		fmt.Printf("Failed to write search message %s\n", err)
		return
	}
}


func (c *Client) readSocketEvents() {

	defer func() {
		fmt.Printf("Close reader\n")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	defer func() {
		r := recover()
		if r != nil {
			time.Sleep(1000 * time.Millisecond)
			fmt.Printf("Recovered %s\n", r)
		}
	}()

	for {
		v := &socketMessage{}

		err := c.conn.ReadJSON(v)
		if err != nil {
			fmt.Printf("Error reading socket %s\n", err)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch v.Name {
		case "seek":
			t := int64(v.Data.(float64) * 1000000000)
			mpdClient.Conn.SeekCur(time.Duration(t), false)
			c.sendSeekMessage()

		case "playlistupdate":
			d := v.Data.([]interface{})
			start := int(d[0].(float64))
			end := int(d[1].(float64))
      // respond with playlist
			c.sendPlaylistUpdateMessage(start, end)

		case "currentsong":
      // respond with current song
			c.sendCurrentSongMessage()

		case "playlistmove":
			d := v.Data.([]interface{})
			// fmt.Printf("Move %s\n", d)
			start := int(d[0].(float64))
			end := int(d[1].(float64))
			position := int(d[2].(float64))
			// fmt.Printf("Move %s %s %s\n")
			if (start != position) {
				mpdClient.Conn.Move(start, end, position)
			}

		// case "play":
    //   // play current
		// 	mpdClient.Conn.PlayID(-1)

		case "playid":
			// -1 for play current
			d := int(v.Data.(float64))
			mpdClient.Conn.PlayID(d)

		case "stop":
			mpdClient.Conn.Stop()

		case "pause":
			mpdClient.Conn.Pause(true)

		case "playnext":
			mpdClient.Conn.Next()

		case "playprev":
			mpdClient.Conn.Previous()

		case "removeid":
			d := int(v.Data.(float64))
			mpdClient.Conn.DeleteID(d)

		case "addpath":
			d := v.Data.([]interface{})
			// fmt.Printf("Move %s\n", d)
			path := d[0].(string)
			position := int(d[1].(float64))
			mpdClient.Conn.AddID(path, position)

		case "search":
			d := v.Data.([]interface{})
			query := d[0].(string)
			start := int(d[1].(float64))
			size := int(d[2].(float64))
			// fmt.Printf("Search %s %s %s\n", query, start, size)
			c.sendSearchMessage(query, start, size)
		}
	}
}


func (c *Client) writeSocketEvents() {

	defer func() {
		fmt.Printf("Close writer\n")
		c.conn.Close()
	}()

	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered %s\n", r)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				fmt.Printf("Hub closed the channel\n")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			switch string(message) {
			case "player":
				c.sendStatusMessage()
				c.sendCurrentSongMessage()

			case "mixer":
				c.sendStatusMessage()

			case "options":
				c.sendStatusMessage()

			case "outputs":
				c.sendStatusMessage()

			case "playlist":
				c.sendPlaylistMessage()
			}

		case <- time.After(1000 * time.Millisecond):
			c.sendSeekMessage()
		}
	}
}


//
// web socket feeder
// based on example https://github.com/gorilla/websocket/blob/master/examples/filewatch/main.go
//
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Serving WS\n")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Printf("%s\n", err)
		}

		fmt.Printf("Serving WS err %s\n", err)
		return
	}

	client := &Client{
		hub: hub,
		conn: ws,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client
	go client.writeSocketEvents()
	go client.readSocketEvents()
}


//
// handle funcs
//
func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Healthcheck\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{"ok"})
}


func search(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Printf("Search database %s\n", params)

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
