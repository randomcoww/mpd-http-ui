//
// Add and remove to elasticsearch
//

package main

import (
	"fmt"
	"context"
	"time"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	Ready chan struct{}
	Down chan struct{}
	url string
	index string
	indexType string
	conn *elastic.Client
}

// new ES client
func NewEsClient(url, index, indexType string) (*EsClient) {
	c := &EsClient{
		Ready: make(chan struct{}, 1),
		Down: make(chan struct{}, 1),
		url: url,
		index: index,
		indexType: indexType,
	}

	c.setStatusDown()
	go c.processLoop()

	return c
}


func (c *EsClient) setStatusReady() {
	c.Ready <- struct{}{}
	fmt.Printf("Elasticsearch ready\n")
}

func (c *EsClient) setStatusDown() {
	c.Down <- struct{}{}
	fmt.Printf("Elasticsearch down\n")
}


func (c *EsClient) processLoop() {
	for {
		select {

		case <-c.Down:
			for {
				time.Sleep(1000 * time.Millisecond)

				fmt.Printf("Connecting to Elasticsearch...\n")
				conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

				if err == nil {
					c.conn = conn

					// get version
					err = c.testVersion()
					if err != nil {
						fmt.Printf("Checking Elasticsearch version...\n")
						continue
					}

					c.setStatusReady()
					break

				} else {
					fmt.Printf("Error connecting to Elasrticsearch\n")
				}
			}
		}
	}
}


func (c *EsClient) testVersion() (error) {
	// ping test
	info, code, err := c.conn.Ping(c.url).Do(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return nil
}


// search database - go through elasticsearch
func (c *EsClient) Get(file string) (*elastic.GetResult, error) {
	get, err := c.conn.Get().
		Index(c.index).
		Type(c.indexType).
		Id(file).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return get, nil
}


// search
func (c *EsClient) Search(query string) (*elastic.SearchResult, error) {
	search, err := c.conn.Search().
		Index(c.index).
		Type(c.indexType).
		Query(elastic.NewSimpleQueryStringQuery(query)).
		Pretty(true).
		Do(ctx)

	if err != nil {
		fmt.Printf("err %s\n", err)
		return nil, err
	}

	return search, nil
}
