package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	es_handler "github.com/randomcoww/go-mpd-es/es_handler"
	mpd_event "github.com/randomcoww/go-mpd-es/mpd_event"
	mpd_handler "github.com/randomcoww/go-mpd-es/mpd_handler"
	// es_handler "local/es_handler"
)

var (
	esIndex, esDocument = "songs", "song"
	mpdClient           *mpd_handler.MpdClient
	esClient            *es_handler.EsClient
	mpdEvent            *mpd_event.MpdEvent
	playlistVersion     int
	playlistLength      int

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
	send chan *socketMessage
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second
)

type response struct {
	Message string
}

type socketMessage struct {
	Name string      `json:"mutation"`
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

	// playlistVersion, playlistLength, _ = getPlaylistStatus()
	updatePlaylistStatus()

	// set mpd repeat by default
	mpdClient.Conn.Repeat(true)

	// serve http
	fmt.Printf("API server start on %s\n", listenUrl)
	log.Fatal(http.ListenAndServe(listenUrl, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r)))
}

//
// broadcast events
//
func createStatusMessage() (*socketMessage, error) {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return nil, err
	}
	return &socketMessage{Data: attrs, Name: "status"}, nil
}

func createCurrentSongMessage() (*socketMessage, error) {
	attrs, err := mpdClient.Conn.CurrentSong()
	if err != nil {
		return nil, err
	}
	return &socketMessage{Data: attrs, Name: "currentsong"}, nil
}

func createSeekMessage() (*socketMessage, error) {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return nil, err
	}
	switch attrs["state"] {
	case "play":
		message := make([]float64, 2)

		elapsed, err := strconv.ParseFloat(attrs["elapsed"], 32)
		if err != nil {
			return nil, err
		}

		duration, err := strconv.ParseFloat(attrs["duration"], 32)
		if err != nil {
			return nil, err
		}

		message[0] = elapsed
		message[1] = duration

		return &socketMessage{Data: message, Name: "seek"}, nil
	}
	return nil, nil
}

func createUpdateDatabaseMessage() *socketMessage {
	return &socketMessage{Name: "updatedb"}
}

func createPlaylistChangedMessage() (*socketMessage, error) {
	curPlaylistLength := playlistLength
	curPlaylistVersion := playlistVersion

	// set new playlistVersion and playlistLength
	updatePlaylistStatus()

	fmt.Printf("MPD playlist update length: %v -> %v\n", curPlaylistLength, playlistLength)
	fmt.Printf("MPD playlist update version: %v -> %v\n", curPlaylistVersion, playlistVersion)

	message := make([]int, 2)

	if playlistLength > curPlaylistLength {
		// Behavior for add to playlist
		// 0. song1
		// 1. song2 <-- added
		// 2. song3 <-- added
		// 3. song4
		// 4. song5
		// Receives: start: 0, end: 3 (new length of playlist)
		addCount := playlistLength - curPlaylistLength
		changeStartPos, _, err := getPlaylistChangePos(curPlaylistVersion)

		if err != nil {
			return nil, err
		}
		// send socket event
		// changeStartPos, changeStartPos + addCount
		fmt.Printf("MPD playlist add positions at: %v count: %v\n", changeStartPos, addCount)

		message[0] = changeStartPos
		message[1] = addCount

		return &socketMessage{Data: message, Name: "playlistadd"}, nil

	} else if playlistLength < curPlaylistLength {
		// Behavior for removed from playlist
		// 0. song1
		// 1. song2 <-- deleting
		// 2. song3 <-- deleting
		// 3. song4
		// 4. song5
		// Receives: start: 1, end: 2 (new length of playlist)
		removeCount := curPlaylistLength - playlistLength
		changeStartPos, _, err := getPlaylistChangePos(curPlaylistVersion)

		if err != nil {
			return nil, err
		}

		// send socket event
		// changeStartPos, changeStartPos + removeCount
		fmt.Printf("MPD playlist delete at: %v count: %v\n", changeStartPos, removeCount)

		message[0] = changeStartPos
		message[1] = removeCount

		return &socketMessage{Data: message, Name: "playlistdelete"}, nil

	} else {
		// Fallback for generic playlist changes (move, shuffle, etc)
		changeStartPos, changeEndPos, err := getPlaylistChangePos(curPlaylistVersion)
		changeCount := changeEndPos - changeStartPos + 1

		if err != nil {
			return nil, err
		}

		fmt.Printf("MPD playlist moved positions at: %v count: %v\n", changeStartPos, changeCount)

		message[0] = changeStartPos
		message[1] = changeCount

		return &socketMessage{Data: message, Name: "playlistmove"}, nil
	}
}

//
// client specific events
//
func createPlaylistQueryMessage(start, end int) (*socketMessage, error) {
	attrs, err := mpdClient.Conn.PlaylistInfo(start, end)
	if err != nil {
		return nil, err
	}
	return &socketMessage{Data: attrs, Name: "playlistquery"}, nil
}

func createSearchMessage(query string, start, size int) (*socketMessage, error) {
	search, err := esClient.Search(query, start, size)
	if err != nil {
		return nil, err
	}

	var result []*json.RawMessage
	for _, hits := range search.Hits.Hits {
		result = append(result, hits.Source)
	}

	message := make([]interface{}, 2)
	message[0] = result
	message[1] = start
	return &socketMessage{Data: message, Name: "search"}, nil
}

//
// broadcaster
//
func (h *Hub) eventBroadcaster() {
	for {
		select {
		case e := <-mpdEvent.Event:
			fmt.Printf("Got MPD event: %s\n", e)

			switch e {
			case "player":
				msg, err := createStatusMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

				msg, err = createCurrentSongMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

				msg, err = createSeekMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "playlist":
				msg, err := createPlaylistChangedMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "mixer", "options", "outputs":
				msg, err := createStatusMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "update":
				msg := createUpdateDatabaseMessage()
				h.broadcast <- msg
			}

		case <-time.After(1000 * time.Millisecond):
			msg, err := createSeekMessage()
			if err != nil {
				break
			}
			h.broadcast <- msg
		}
	}
}

//
// send broadcast events to each client
//
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
		case msg, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				fmt.Printf("Hub closed the channel\n")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(*msg)
		}
	}
}

//
// read messages from client
//
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

			// client specific playlist query
		case "playlistquery":
			d := v.Data.([]interface{})
			start := int(d[0].(float64))
			end := int(d[1].(float64))
			msg, err := createPlaylistQueryMessage(start, end)
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

			// client specific current song query
		case "currentsong":
			msg, err := createCurrentSongMessage()
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

			// global playlist items moved
			// send only and allow server to emit event
		case "playlistmove":
			d := v.Data.([]interface{})
			// fmt.Printf("Move %s\n", d)
			start := int(d[0].(float64))
			end := int(d[1].(float64))
			position := int(d[2].(float64))
			// fmt.Printf("Move %s %s %s\n")
			if start != position {
				mpdClient.Conn.Move(start, end, position)
			}

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
			// fmt.Printf("remove %v\n", d)
			mpdClient.Conn.DeleteID(d)

		case "addpath":
			d := v.Data.([]interface{})
			// fmt.Printf("Move %s\n", d)
			path := d[0].(string)
			position := int(d[1].(float64))
			mpdClient.Conn.AddID(path, position)

			// client specific database search
		case "search":
			d := v.Data.([]interface{})
			query := d[0].(string)
			start := int(d[1].(float64))
			size := int(d[2].(float64))
			// fmt.Printf("Search %s %s %s\n", query, start, size)
			msg, err := createSearchMessage(query, start, size)
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

		case "clear":
			mpdClient.Conn.Clear()

		case "updatedb":
			mpdClient.Conn.Update("")
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
		hub:  hub,
		conn: ws,
		send: make(chan *socketMessage, 256),
	}

	client.hub.register <- client
	go client.writeSocketEvents()
	go client.readSocketEvents()
}

//
// http handle funcs
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

//
// helpers
//
func parseNum(input string) int {
	v, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Error parsing param %s: %s\n", input, err)
		v = -1
	}

	return v
}

// update global playlistVersion and playlistLength
func updatePlaylistStatus() error {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return err
	}

	playlistVersion, err = strconv.Atoi(attrs["playlist"])
	if err != nil {
		return err
	}

	playlistLength, err = strconv.Atoi(attrs["playlistlength"])
	if err != nil {
		return err
	}

	return nil
}

// when playlist changes, get the start and end index of change
func getPlaylistChangePos(playlistVersion int) (int, int, error) {
	attrs, err := mpdClient.PlChangePosId(playlistVersion, -1, -1)

	if err != nil {
		return 0, 0, err
	}

	var (
		startPos = 0
		endPos   = 0
	)

	if len(attrs) > 0 {
		v, ok := attrs[0]["cpos"]
		if ok {
			i, err := strconv.Atoi(v)
			if err != nil {
				return 0, 0, err
			}
			startPos = i
		}

		v, ok = attrs[len(attrs)-1]["cpos"]
		if ok {
			i, err := strconv.Atoi(v)
			if err != nil {
				return 0, 0, err
			}
			endPos = i
		}

		return startPos, endPos, nil
	}

	// if no result, last N items were removed
	// ignore endPos
	return playlistLength, -1, nil
}
