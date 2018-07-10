//
// Add and remove to elasticsearch
//

package main

import (
	"fmt"
	"context"
	"time"
	"errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	up chan struct{}
	down chan struct{}
	indexDown chan struct{}

	url string
	index string
	indexType string
	conn *elastic.Client
}


// new ES client
func NewEsClient(url, index, indexType string) (*EsClient) {
	c := &EsClient{
		up: make(chan struct{}, 1),
		down: make(chan struct{}, 1),
		indexDown: make(chan struct{}, 1),

		url: url,
		index: index,
		indexType: indexType,
	}

	c.down <- struct{}{}
	go c.processLoop()

	return c
}


func (c *EsClient) processLoop() {
	for {
		select {

		case <-c.down:
			for {
				err := c.connect()
				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}
				break
			}
			c.indexDown <- struct{}{}

		// test getting or creating index
		case <-c.indexDown:
			for {
				err := c.getIndex()
				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}
				break
			}
			c.up <- struct{}{}

    // ping
		case <-time.After(1000 * time.Millisecond):
			_, _, err := c.conn.Ping(c.url).Do(ctx)

			if err != nil {
				c.down <- struct{}{}
			}
		}
	}
}


// get connection
func (c *EsClient) connect() (error) {
	fmt.Printf("Connecting to Elasticsearch...\n")
	conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

	if err != nil {
		return err
	}

	fmt.Printf("Connected to Elasticsearch\n")

	// defer conn.Close()
	c.conn = conn

	return nil
}


// create index with provided mapping
func (c *EsClient) getIndex() (error) {
	fmt.Printf("Checking Elasticsearch index...\n")

	exists, err := c.conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("Elasticsearch index not found")
	}

	fmt.Printf("Got Elasticsearch index\n")
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
		return nil, err
	}

	return search, err
}
