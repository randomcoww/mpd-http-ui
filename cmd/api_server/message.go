package api_server

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
	es "github.com/randomcoww/go-mpd-es/pkg/es_handler"
	event "github.com/randomcoww/go-mpd-es/pkg/mpd_event"
	mpd "github.com/randomcoww/go-mpd-es/pkg/mpd_handler"
	"github.com/sirupsen/logrus"
)

type socketMessage struct {
	Name string      `json:"mutation"`
	Data interface{} `json:"value"`
}

var (
	esIndex, esDocument = "songs", "song"
	mpdClient           *mpd.MpdClient
	esClient            *es.EsClient
	mpdEvent            *event.MpdEvent
	playlistVersion     int
	playlistLength      int

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

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
	prevPlaylistVersion := playlistVersion
	prevPlaylistLength := playlistLength

	updatePlaylistStatus()

	logrus.Infof("MPD playlist update length: %v -> %v\n", prevPlaylistLength, playlistLength)
	logrus.Infof("MPD playlist update version: %v -> %v\n", prevPlaylistVersion, playlistVersion)

	// message := make([]int, 2)

	if playlistLength > prevPlaylistLength {
		// Behavior for add to playlist
		// 0. song1
		// 1. song2 <-- added
		// 2. song3 <-- added
		// 3. song4
		// 4. song5
		// Receives: start: 0, end: 3 (new length of playlist)
		addCount := playlistLength - prevPlaylistLength
		changeStartPos, _, err := getPlaylistChangePos(prevPlaylistVersion)
		if err != nil {
			return nil, err
		}
		// send socket event
		// changeStartPos, changeStartPos + addCount
		logrus.Infof("MPD playlist add positions at: %v count: %v\n", changeStartPos, addCount)

		message := make([]int, 2)
		message[0] = changeStartPos
		message[1] = addCount

		return &socketMessage{Data: message, Name: "playlistadd"}, nil

	} else {
		// Fallback for generic playlist changes (move, shuffle, etc)
		changeStartPos, changeEndPos, err := getPlaylistChangePos(prevPlaylistVersion)
		changeCount := changeEndPos - changeStartPos + 1
		// Items were removed
		removeCount := prevPlaylistLength - playlistLength

		if err != nil {
			return nil, err
		}

		logrus.Infof("MPD playlist changed positions at: %v count: %v length: %v\n", changeStartPos, changeCount, playlistLength)

		message := make([]int, 4)
		// Report changes from API
		message[0] = changeStartPos
		message[1] = changeCount
		// If playlist is smaller, report items to remove
		message[2] = prevPlaylistLength
		message[3] = removeCount

		return &socketMessage{Data: message, Name: "playlistchange"}, nil
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

// when playlist changes, get the start and end index of change
func getPlaylistChangePos(version int) (int, int, error) {
	attrs, err := mpdClient.PlChangePosId(version, -1, -1)

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
	return -1, -1, nil
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
