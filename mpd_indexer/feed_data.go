//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
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

var (
	esIndex, esDocument = "songs", "song"
)


// main
func NewDataFeeder(logFile, mpdUrl, esUrl string) (error) {
	logParser, err := NewLogEventParser(logFile)

	if err != nil {
		return err
	}

	mpdClient := NewMpdClient("tcp", mpdUrl)
	esClient := NewEsClient(esUrl, esIndex, esDocument, esMapping)

	<-mpdClient.up
	<-esClient.up

	for {
		select {
		case c := <- logParser.added:
			fmt.Printf("Add item event: %s\n", c)

			attr := mpdClient.GetDatabaseItem(c)
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
			esClient.IndexBulk(addIndex)

		case c := <- logParser.deleted:
			fmt.Printf("Delete item event: %s\n", c)
			esClient.DeleteBluk(c)

		case <- time.After(1000 * time.Millisecond):
		}
	}

	return nil
}
