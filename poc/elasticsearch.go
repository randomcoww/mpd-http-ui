package main

import (
	"fmt"
	"context"
	// "encoding/json"
	elastic "gopkg.in/olivere/elastic.v5"
)

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

const mapping = `
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


type SongStore struct {
	url       string
	index     string
	indexType string
	ctx       context.Context
	Client    *elastic.Client
}


func SetupClient(url, index, indexType, mapping string) (*SongStore, error) {
	// create elasticsearch client based on example:
	// https://olivere.github.io/elastic/

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		// Handle error
		return nil, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(url)
	if err != nil {
		// Handle error
		return nil, err
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			return nil, err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	c := &SongStore{
		url:    url,
		ctx:    ctx,
		Client: client,
		index:  index,
		indexType:  indexType,
	}

	return c, nil
}


func (c *SongStore) IndexSong(s Song) error {
	_, err := c.Client.Index().
		Index(c.index).
		Type(c.indexType).
		Id(s.File).
		BodyJson(s).
		Do(c.ctx)

	if err != nil {
		// Handle error
		return err
	}
	return nil
}


func (c *SongStore) GetSong(file string) (*elastic.GetResult, error) {
	get, err := c.Client.Get().
		Index(c.index).
		Type(c.indexType).
		Id(file).
		Do(c.ctx)

	if err != nil {
		// Handle error
		return nil, err
	}

	return get, nil
}


func (c *SongStore) DeleteSong(file string) (*elastic.DeleteResponse, error) {
	delete, err := c.Client.Delete().
		Index(c.index).
		Type(c.indexType).
		Id(file).
		Do(c.ctx)

	if err != nil {
		// Handle error
		return nil, err
	}

	return delete, nil
}


func main() {
	c, err := SetupClient("http://127.0.0.1:9200", "songs", "song", mapping)

	if err != nil {
		fmt.Println(err)
	}

	testIndex := Song{
		File: "testfile",
		Title: "testtitle",
	}

	err = c.IndexSong(testIndex)

	if err != nil {
		fmt.Println(err)
	}

	_, err = c.GetSong("testfile")

	if err != nil {
		fmt.Println(err)
	}

	deleted, err := c.DeleteSong("testfile")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(deleted)
}
