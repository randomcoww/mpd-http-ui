//
// Add and remove to elasticsearch
//

package es_handler

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	eventHub *eventHub

	url       string
	index     string
	indexType string
	mapping   string

	conn *elastic.Client
	bulk *elastic.BulkService
}

// new ES client
func NewEsClient(url, index, indexType, mapping string) *EsClient {
	c := &EsClient{
		eventHub: newEventHub(),

		url:       url,
		index:     index,
		indexType: indexType,
		mapping:   mapping,
	}

	go c.processLoop()
	go c.processBulk()

	// initial state is down
	c.eventHub.send <- "api_down"

	return c
}

func (c *EsClient) processLoop() {
	errClient := c.eventHub.newClient([]string{
		"api_down", "index_down",
	})

	for {
		select {
		case event := <-errClient.events:
			switch event {
			case "api_down":
				c.eventHub.send <- "index_down"

				err := c.connect()
				if err != nil {
					logrus.Infof("API ready")
					c.eventHub.send <- "api_ready"
				} else {
					logrus.Error("Service down")
				}

			case "index_down":
				err := c.getOrCreateIndex()
				if err == nil {
					logrus.Infof("Index ready")
					c.eventHub.send <- "index_ready"
				} else {
					logrus.Error("Index not found")
				}
			}

		// healthcheck
		case <-time.After(2000 * time.Millisecond):
			_, _, err := c.conn.Ping(c.url).Do(ctx)

			if err != nil {
				logrus.Errorf("Ping down: %v", err)
				c.eventHub.send <- "api_down"
			} else {
				// logrus.Infof("Ping")
			}
		}
	}
}

// run bulk processing job
func (c *EsClient) processBulk() {
	errClient := c.eventHub.newClient([]string{
		"api_down", "index_down",
	})

	readyClient := c.eventHub.newClient([]string{
		"index_ready",
	})

	updateClient := c.eventHub.newClient([]string{
		"index_update",
	})

	for {
		select {
		case event := <-errClient.events:
			switch event {
			case "api_down", "index_down":
				errClient.drain()
				logrus.Error("Bulk update - wait for index to become ready")
				readyClient.waitEvent("index_ready")
			}
		}

		select {
		case event := <-updateClient.events:
			switch event {
			case "index_update":
				if c.bulk.NumberOfActions() > 0 {

					// bulk process seems to lose data on failure
					// make sure index is accessible before trying bulk write
					exists, err := c.conn.IndexExists(c.index).Do(ctx)

					if err != nil {
						// API down?
						logrus.Errorf("Bulk update: Test API failed: %v", err)
						c.eventHub.send <- "api_down"
					} else if !exists {
						// Index not found
						logrus.Errorf("Bulk update: Index not found: %s", c.index)
						c.eventHub.send <- "index_down"
					} else {
						// Ok to try updating
						_, err := c.bulk.Do(ctx)
						if err != nil {
							logrus.Errorf("Bulk update: Failed: %v", err)
						} else {
							logrus.Info("Bulk update: Success")
						}
					}
				}
			}
		// Add some throtting for bulk update
		case <-time.After(2000 * time.Millisecond):
		}
	}
}

// get connection
func (c *EsClient) connect() error {
	conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

	if err != nil {
		return err
	}

	// defer conn.Close()
	c.conn = conn
	c.bulk = conn.Bulk()

	return nil
}

// create index with provided mapping
func (c *EsClient) getOrCreateIndex() error {
	exists, err := c.conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		logrus.Info("Index exists: %s", c.index)
		return nil
	}

	if len(c.mapping) == 0 {
		return errors.New("Mapping not provided")
	}

	_, err = c.conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
	if err != nil {
		return err
	}

	logrus.Info("Index created: %s", c.index)
	return nil
}

// Add index to next bulk update
func (c *EsClient) IndexBulk(id string, s interface{}) {
	c.bulk.Add(elastic.NewBulkIndexRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id).
		Doc(s))

	c.eventHub.send <- "index_update"
}

// Add deletion to next bulk update
func (c *EsClient) DeleteBluk(id string) {
	c.bulk.Add(elastic.NewBulkDeleteRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id))

	c.eventHub.send <- "index_update"
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
func (c *EsClient) Search(query string, start, size int) (*elastic.SearchResult, error) {
	search, err := c.conn.Search().
		Index(c.index).
		Type(c.indexType).
		Query(elastic.NewSimpleQueryStringQuery(query)).
		Pretty(true).
		From(start).
		Size(size).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return search, err
}
