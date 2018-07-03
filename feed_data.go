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
	AlbumArtist  string `json:"albumartist,omitempty"`
	Album        string `json:"album,omitempty"`
	Track        int    `json:"track,omitempty"`
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
				"albumartist":{
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
			fmt.Printf("Add_event: %s\n", c)

			attr := mpdClient.GetInfo(c)
			addIndex := Song{
				File: c,
				Title: attr["Title"],
			}
			esClient.Index(addIndex)

		case c := <- logParser.deleted:
			fmt.Printf("delete_event: %s\n", c)
			esClient.Delete(c)

		case <- time.After(1000 * time.Millisecond):
    }
  }

  return nil
}
