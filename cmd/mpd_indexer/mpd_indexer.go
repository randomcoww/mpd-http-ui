//
// get add and remove item events by parsing the mpd log
//

package mpd_indexer

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"time"

	es "github.com/randomcoww/go-mpd-es/pkg/es_handler"
	mpd "github.com/randomcoww/go-mpd-es/pkg/mpd_handler"
	"github.com/sirupsen/logrus"
)

type LogEvents struct {
	MpdReady chan bool
	EsReady  chan bool
	added    chan string
	deleted  chan string
}

// elasticsearch stuff
type Song struct {
	File     string `json:"file"`
	Date     string `json:"date,omitempty"`
	Duration string `json:"duration,omitempty"`
	Composer string `json:"composer,omitempty"`
	Album    string `json:"album,omitempty"`
	Track    string `json:"track,omitempty"`
	Title    string `json:"title,omitempty"`
	Artist   string `json:"artist,omitempty"`
	Genre    string `json:"genre,omitempty"`
}

const esMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"song":{
			"properties":{
				"file":{
					"type":"keyword"
				},
				"date":{
					"type":"date"
				},
				"duration":{
					"type":"text"
				},
				"composer":{
					"type":"text"
				},
				"album":{
					"type":"text"
				},
				"track":{
					"type":"text"
				},
				"title":{
					"type":"text"
				},
				"artist":{
					"type":"text"
				},
				"genre":{
					"type":"text"
				}
			}
		}
	}
}`

var (
	addedString   = "update: added "
	deletedString = "update: removing "
	esIndex       = "songs"
	esDocument    = "song"
)

func newLogEvents() *LogEvents {
	return &LogEvents{
		added:   make(chan string),
		deleted: make(chan string),
	}
}

// parse logs and send items to add and remove channels
func (e *LogEvents) readLog(logFile string) {
	reader, _ := createNamesPipe(logFile)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		logrus.Infof("%s", line)

		if strings.Contains(line, addedString) {
			str := strings.Split(line, addedString)
			e.added <- strings.TrimSuffix(str[len(str)-1], "\n")

		} else if strings.Contains(line, deletedString) {
			str := strings.Split(line, deletedString)
			e.deleted <- strings.TrimSuffix(str[len(str)-1], "\n")
		}
	}
}

func (e *LogEvents) parseLog(mpdClient *mpd.MpdClient, esClient *es.EsClient) {
	for {
		select {
		case c := <-e.added:
			logrus.Infof("Add item event: %s", c)

			attr := mpdClient.GetDatabaseItem(c)
			addIndex := Song{
				File:     c,
				Date:     attr["date"],
				Duration: attr["duration"],
				Composer: attr["composer"],
				Album:    attr["album"],
				Track:    attr["track"],
				Title:    attr["title"],
				Artist:   attr["artist"],
				Genre:    attr["genre"],
			}
			esClient.IndexBulk(c, addIndex)

		case c := <-e.deleted:
			logrus.Infof("Delete item event: %s", c)
			esClient.DeleteBluk(c)

		case <-time.After(1000 * time.Millisecond):
		}
	}
}

func createNamesPipe(logFile string) (*bufio.Reader, error) {
	os.Remove(logFile)
	err := syscall.Mkfifo(logFile, 0600)
	if err != nil {
		logrus.Infof("Failed create named pipe: %v", err)
		return nil, err
	}

	f, err := os.OpenFile(logFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		logrus.Infof("Failed open named pipe: %v", err)
		return nil, err
	}

	return bufio.NewReader(f), nil
}
