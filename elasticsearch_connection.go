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
	Conn *elastic.Client
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
	go c.reconnectLoop()

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


func (c *EsClient) reconnectLoop() {
	for {
		select {

		case <-c.Down:
			for {
				time.Sleep(1000 * time.Millisecond)

				fmt.Printf("Connecting to Elasticsearch...\n")
				conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

				if err == nil {
					c.Conn = conn

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
		}
	}
}


func (c *EsClient) testVersion() (error) {
	// ping test
	info, code, err := c.Conn.Ping(c.url).Do(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return nil
}


func (c *EsClient) createIndex() (error) {
	// index test
	exists, err := c.Conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := c.Conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}


func (c *EsClient) Index(s Song) (*elastic.IndexResponse) {
	for {
		create, err := c.Conn.Index().
			Index(c.index).
			Type(c.indexType).
			Id(s.File).
			BodyJson(s).
			Do(ctx)

		if err == nil {
			fmt.Printf("Created Elasticsearch entry %s\n", s)
			return create

		} else {
			c.setStatusDown()
			fmt.Printf("Start Elasticsearch ready wait\n")
			<-c.Ready
		}
	}
}


func (c *EsClient) Delete(id string) (*elastic.DeleteResponse) {
	for {
		delete, err := c.Conn.Delete().
			Index(c.index).
			Type(c.indexType).
			Id(id).
			Do(ctx)

		if err == nil {
			fmt.Printf("Deleted Elasticsearch entry %s\n", id)
			return delete

		} else {
			c.setStatusDown()
			fmt.Printf("Start Elasticsearch ready wait\n")
			<-c.Ready
		}
	}
}

// func (c *EsClient) Get(file string) (*elastic.GetResult, error) {
// 	get, err := c.Client.Get().
// 		Index(c.index).
// 		Type(c.indexType).
// 		Id(file).
// 		Do(ctx)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return get, nil
// }
