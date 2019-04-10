//
// Add and remove to elasticsearch
//

package elasticsearch

import (
	"context"
	"time"

	"github.com/randomcoww/go-mpd-es/pkg/util"
	"github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	eventHub *util.EventHub

	url       string
	index     string
	indexType string
	mapping   string

	conn *elastic.Client
	bulk *elastic.BulkService
}

// new ES client
func NewEsClient(url, index, indexType, mapping string) *EsClient {

	logrus.Infof("EsClient: Start")

	c := &EsClient{
		eventHub: util.NewEventHub(),

		url:       url,
		index:     index,
		indexType: indexType,
		mapping:   mapping,
	}

	c.setReady()
	go c.runRecovery()
	go c.runIndexUpdate()

	return c
}

//
// get connection
//

func (c EsClient) pingTest() bool {
	_, _, err := c.conn.Ping(c.url).Do(ctx)
	return err == nil
}

func (c *EsClient) waitConnect() {
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

			if err == nil {
				c.conn = conn
				c.bulk = conn.Bulk()

				logrus.Infof("EsClient: Connection ready")
				return
			}
		}
	}
}

func (c *EsClient) waitPingState(state bool) {
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			if c.pingTest() == state {
				logrus.Infof("EsClient: Ping state changed: ", state)
				return
			}
		}
	}
}

func (c *EsClient) setReady() {
	c.waitConnect()
	c.waitPingState(true)
	c.eventHub.Send <- "api_ready"
}

func (c *EsClient) setDown() {
	c.waitPingState(false)
	c.eventHub.Send <- "api_down"
}

func (c *EsClient) runRecovery() {
	apiErrClient := c.eventHub.NewClient([]string{"api_down"})
	apiReadyClient := c.eventHub.NewClient([]string{"api_ready"})

	for {
		select {
		case <-apiErrClient.Events:
			c.setReady()
			apiErrClient.Drain()

		case <-apiReadyClient.Events:
			c.setDown()
			apiReadyClient.Drain()
		}
	}
}

//
// Update index
//

func (c *EsClient) getOrCreateIndex() error {
	exists, err := c.conn.IndexExists(c.index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		logrus.Info("EsClient: Index exists: %s", c.index)
		return nil
	}

	_, err = c.conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
	if err != nil {
		return err
	}

	logrus.Info("EsClient: Index created: %s", c.index)
	return nil
}

func (c *EsClient) waitIndex() {
	apiReadyClient := c.eventHub.NewClient([]string{"api_ready"})
	for {
		if err := c.getOrCreateIndex(); err == nil {
			return
		}

		c.eventHub.Send <- "api_down"
		apiReadyClient.WaitEvent("api_ready")
	}
}

// run bulk processing job
func (c *EsClient) runIndexUpdate() {
	updateClient := c.eventHub.NewClient([]string{"index_update"})

	for {
		select {
		case <-updateClient.Events:
			c.waitIndex()
			updateClient.Drain()

			// Ok to try updating
			_, err := c.bulk.Do(ctx)
			if err != nil {
				logrus.Errorf("EsClient: Bulk update: Failed: %v", err)
			} else {
				logrus.Info("EsClient: Bulk update: Success")
			}
		// Add some throtting for bulk update
		case <-time.After(2000 * time.Millisecond):
		}
	}
}

// Add index to next bulk update
func (c *EsClient) IndexBulk(id string, s interface{}) {
	c.bulk.Add(elastic.NewBulkIndexRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id).
		Doc(s))

	c.eventHub.Send <- "index_update"
}

// Add deletion to next bulk update
func (c *EsClient) DeleteBluk(id string) {
	c.bulk.Add(elastic.NewBulkDeleteRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id))

	c.eventHub.Send <- "index_update"
}

//
// Search
//

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
