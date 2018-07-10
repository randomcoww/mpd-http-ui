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
	up chan struct{}
	down chan struct{}
	indexDown chan struct{}

	url string
	index string
	indexType string
	mapping string
	conn *elastic.Client
	bulk *elastic.BulkService
}


// new ES client
func NewEsClient(url, index, indexType, mapping string) (*EsClient) {
	c := &EsClient{
		up: make(chan struct{}, 1),
		down: make(chan struct{}, 1),
		indexDown: make(chan struct{}, 1),

		url: url,
		index: index,
		indexType: indexType,
		mapping: mapping,
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
				err := c.getOrCreateIndex()
				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}
				break
			}
			c.up <- struct{}{}

		case <-time.After(1000 * time.Millisecond):
			err := c.processBulk()

			if err != nil {
				c.down <- struct{}{}
			}
		}
	}
}


// run bulk processing job
func (c *EsClient) processBulk() (error) {
	if c.bulk.NumberOfActions() > 0 {
		_, err := c.bulk.Do(ctx)

		if err != nil {
			fmt.Printf("Error processsing Elasticsearch bulk\n")
			return err

		} else {
			fmt.Printf("Processsed Elasticsearch bulk\n")
		}
	}
	return nil
}


// get connection
func (c *EsClient) connect() (error) {
	if c.conn != nil {
		_, _, err := c.conn.Ping(c.url).Do(ctx)

		if err != nil {
      //
		} else {
			return nil
		}
	}

	fmt.Printf("Connecting to Elasticsearch...\n")
	conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

	if err != nil {
		return err
	}

	fmt.Printf("Connected to Elasticsearch\n")
	// defer conn.Close()
	c.conn = conn
	c.bulk = conn.Bulk()

	return nil
}


// create index with provided mapping
func (c *EsClient) getOrCreateIndex() (error) {
	fmt.Printf("Checking Elasticsearch index...\n")

	exists, err := c.conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := c.conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
		if err != nil {
			return err

		} else {
			fmt.Printf("Created Elasticsearch index\n")
			return nil
		}
	}

	fmt.Printf("Got Elasticsearch index\n")
	return nil
}


// Add index to next bulk update
func (c *EsClient) IndexBulk(s Song) {
	c.bulk.Add(elastic.NewBulkIndexRequest().
		Index(c.index).
		Type(c.indexType).
		Id(s.File).
		Doc(s))
}

// Add deletion to next bulk update
func (c *EsClient) DeleteBluk(id string) {
	c.bulk.Add(elastic.NewBulkDeleteRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id))
}
