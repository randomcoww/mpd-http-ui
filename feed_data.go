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
  mpdClient := NewMpdClient("tcp", "localhost:6600")
  logParser, err := NewLogEventParser("env/mpd_mount/logs/log")
	esClient := NewEsClient("http://127.0.0.1:9200", "songs", "song", esMapping)

  if err != nil {
    return err
  }

  for {
    select {
		case c := <- logParser.added:
			fmt.Println("add_event:", c)

			_, err := esClient.EsConn()
			if err != nil {
				fmt.Printf("error connecting to elasticsearch %s \n", err)
			}

      // _, err = mpdClient.MpdConn()

      if err == nil {
				attrs, err := mpdClient.Conn.ListAllInfo(c)

				if err == nil {
					if len(attrs) > 0 {
						attr := attrs[0]

						fmt.Printf("%s \n", attr)

						addIndex := Song{
							File: c,
							Title: attr["Title"],
						}

						err = esClient.Index(addIndex)
						if err != nil {
							fmt.Printf("error indexing %s \n", err)
						}

						fmt.Printf("es add %s %s \n", c, addIndex)
					}

				} else {
					fmt.Printf("error parsing %s %s \n", c, err)
				}
      }

		case c := <- logParser.deleted:
			fmt.Println("delete_event:", c)

			_, err := esClient.EsConn()
			if err != nil {
				fmt.Printf("error connecting to elasticsearch %s \n", err)
			}

			_, err = esClient.Delete(c)
			if err != nil {
				fmt.Printf("error deleting %s \n", err)
			}

			fmt.Printf("es delete %s \n", c)

		case <- time.After(1000 * time.Millisecond):
    }
  }

  return nil
}
