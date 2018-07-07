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
	mapping string
	conn *elastic.Client
	bulk *elastic.BulkService
}

// new ES client
func NewEsClient(url, index, indexType, mapping string) (*EsClient) {
	c := &EsClient{
		Ready: make(chan struct{}, 1),
		Down: make(chan struct{}, 1),
		url: url,
		index: index,
		indexType: indexType,
		mapping: mapping,
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
					c.bulk = conn.Bulk()

					// get version
					err = c.testVersion()
					if err != nil {
						fmt.Printf("Checking Elasticsearch version...\n")
						continue
					}

					// detect or create index
					err = c.createIndex()
					if err != nil {
						fmt.Printf("Checking or creating Elasticsearch index...\n")
						continue
					}

					c.setStatusReady()
					break

				} else {
					fmt.Printf("Error connecting to Elasrticsearch\n")
				}
			}

		case <-time.After(1000 * time.Millisecond):
			if c.bulk.NumberOfActions() > 0 {

				_, err := c.bulk.Do(ctx)

				if err == nil {
					fmt.Printf("Processsed elasticsearch bulk operation\n")

				} else {
					c.setStatusDown()
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


func (c *EsClient) createIndex() (error) {
	// index test
	exists, err := c.conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := c.conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}


func (c *EsClient) Index(s Song) {
  // add to bulk opetation
	bulk := elastic.NewBulkIndexRequest().
		Index(c.index).
		Type(c.indexType).
		Id(s.File).
		Doc(s)

	c.bulk.Add(bulk)
}


func (c *EsClient) Delete(id string) {
  // add to bulk operation
	bulk := elastic.NewBulkDeleteRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id)

	c.bulk.Add(bulk)
}
