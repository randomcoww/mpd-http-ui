//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
  // elastic "gopkg.in/olivere/elastic.v5"
)

// elasticsearch stuff
type Song struct {
	File         string `json:"file"`
	Date         string `json:"date,omitempty"`
	Duration     string `json:"duration,omitempty"`
	Composer     string `json:"composer,omitempty"`
	Album        string `json:"album,omitempty"`
	Track        string `json:"track,omitempty"`
	Title        string `json:"title,omitempty"`
	Artist       string `json:"artist,omitempty"`
	Genre        string `json:"genre,omitempty"`
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

// main
func NewDataFeeder() (error) {
	logParser, err := NewLogEventParser("env/mpd_mount/logs/log")

  if err != nil {
    return err
  }

	mpdClient := NewMpdClient("tcp", "localhost:6600")
	<-mpdClient.Ready

	esClient := NewEsClient("http://127.0.0.1:9200", "songs", "song", esMapping)
	<-esClient.Ready

  for {
    select {
		case c := <- logParser.added:
			fmt.Printf("Add item event: %s\n", c)

			// attrs, _ = mpdClient.Conn.PlaylistContents("dir1/test1.cue")
			// fmt.Printf("Got MPD playlist file %s\n", attrs)

			attr := mpdClient.GetInfo(c)
			addIndex := Song{
				File: c,
				Date: attr["date"],
				Duration: attr["duration"],
				Composer: attr["composer"],
				Album: attr["album"],
				Track: attr["track"],
				Title: attr["title"],
				Artist: attr["artist"],
				Genre: attr["genre"],
			}
			esClient.Index(addIndex)

		case c := <- logParser.deleted:
			fmt.Printf("Delete item event: %s\n", c)
			esClient.Delete(c)

		case e := <- mpdClient.Events:
			fmt.Printf("MPD event reader: %s\n", e)

		case <- time.After(1000 * time.Millisecond):
    }
  }

  return nil
}
