//
// Add and remove to elasticsearch
//

package es_handler

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	up          chan struct{}
	down        chan struct{}
	indexDown   chan struct{}
	bulkRequest chan struct{}

	url       string
	index     string
	indexType string
	mapping   string

	conn *elastic.Client
	bulk *elastic.BulkService

	Ready chan struct{}
}

// new ES client
func NewEsClient(url, index, indexType, mapping string) *EsClient {
	c := &EsClient{
		up:          make(chan struct{}, 1),
		down:        make(chan struct{}, 1),
		indexDown:   make(chan struct{}, 1),
		bulkRequest: make(chan struct{}, 1),

		url:       url,
		index:     index,
		indexType: indexType,
		mapping:   mapping,

		Ready: make(chan struct{}, 1),
	}

	c.setState(c.down)
	go c.processLoop()
	go c.processBulk()

	return c
}

func (c *EsClient) setState(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (c *EsClient) drainState(ch chan struct{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
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
			c.setState(c.indexDown)

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
			c.setState(c.up)
			c.setState(c.Ready)

		// ping and reconnect
		case <-time.After(10000 * time.Millisecond):
			_, _, err := c.conn.Ping(c.url).Do(ctx)
			if err != nil {
				logrus.Infof("ES ping down %s", err)
				c.setState(c.down)

			} else {
				// logrus.Infof("ES ping")
			}
		}
	}
}

// run bulk processing job
func (c *EsClient) processBulk() {
	for {
		select {
		case <-c.bulkRequest:

			if c.bulk.NumberOfActions() > 0 {

				for {
					time.Sleep(2000 * time.Millisecond)

					// bulk process seems to lose data on failure
					// make sure index is accessible before trying bulk write
					exists, err := c.conn.IndexExists(c.index).Do(ctx)
					if err != nil {
						logrus.Errorf("ES index test - %s", err)
						continue
					}

					if !exists {
						logrus.Infof("ES index test - not found")
						continue
					}

					logrus.Infof("ES index test - success")
					break
				}

				c.drainState(c.bulkRequest)
				_, err := c.bulk.Do(ctx)
				if err != nil {
					logrus.Errorf("Error processsing ES bulk %s", err)

				} else {
					logrus.Infof("Processsed ES bulk")
				}
			}
		}
	}
}

// get connection
func (c *EsClient) connect() error {
	logrus.Infof("Connecting to ES...")
	conn, err := elastic.NewSimpleClient(elastic.SetURL(c.url))
	if err != nil {
		return err
	}

	logrus.Infof("Connected to ES")
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
		logrus.Infof("ES index exists")
		return nil
	}

	if len(c.mapping) == 0 {
		return fmt.Errorf("ES mapping not provided")
	}

	_, err = c.conn.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
	if err != nil {
		return err
	}

	logrus.Infof("Created ES index")
	return nil
}

// Add index to next bulk update
func (c *EsClient) IndexBulk(id string, s interface{}) {
	c.bulk.Add(elastic.NewBulkIndexRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id).
		Doc(s))

	c.setState(c.bulkRequest)
}

// Add deletion to next bulk update
func (c *EsClient) DeleteBluk(id string) {
	c.bulk.Add(elastic.NewBulkDeleteRequest().
		Index(c.index).
		Type(c.indexType).
		Id(id))

	c.setState(c.bulkRequest)
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
