//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package server

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

const (
	esSongMapping = `
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
	esSongIndex    = "songs"
	esSongDocument = "song"
)
