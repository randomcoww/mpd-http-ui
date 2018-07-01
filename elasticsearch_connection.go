package main

import (
	"fmt"
	"context"
	"time"
	// "encoding/json"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ctx = context.Background()
)

type EsClient struct {
	url       string
	index     string
	indexType string
  mapping   string
	// Ctx       context.Context
	Client    *elastic.Client
}

// new ES client
func NewEsClient(url, index, indexType, mapping string) (*EsClient) {
  s := &EsClient{
    url: url,
    index: index,
    indexType: indexType,
    mapping: mapping,
  }

  return s
}

// get or wait for elasticsearch connection
func (c *EsClient) EsConn() (*EsClient, error) {

  if c.Client != nil {
    // ping := s.Client.Ping()
    //
    // ping.URL(url)
    // result, statusCode, err := ping.Do(s.Ctx)
    //
    // if err != nil {
    //   return s.Client, nil
    // }

    return c, nil
  }

  for {
    client, err := elastic.NewSimpleClient(elastic.SetURL(c.url))

    if err != nil {
      fmt.Println("cannot connect to Elasticsearch")
			time.Sleep(1000 * time.Millisecond)

    } else {
      c.Client = client
      // s.Ctx = context.Background()

      // ping test
      info, code, err := client.Ping(c.url).Do(ctx)
    	if err != nil {
    		// Handle error
    		return nil, err
    	}
			fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

      // index test
      exists, err := client.IndexExists(c.index).Do(ctx)
    	if err != nil {
    		// Handle error
    		return nil, err
    	}

      // create if it doesn't exist
    	if !exists {
    		// Create a new index.
    		createIndex, err := client.CreateIndex(c.index).BodyString(c.mapping).Do(ctx)
    		if err != nil {
    			// Handle error
    			return nil, err
    		}
    		if !createIndex.Acknowledged {
    			// Not acknowledged
    		}
    	}

      return c, nil
    }
  }
}


func (c *EsClient) Index(s Song) error {
	_, err := c.Client.Index().
		Index(c.index).
		Type(c.indexType).
		Id(s.File).
		BodyJson(s).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}


func (c *EsClient) Get(file string) (*elastic.GetResult, error) {
	get, err := c.Client.Get().
		Index(c.index).
		Type(c.indexType).
		Id(file).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return get, nil
}


func (c *EsClient) Delete(file string) (*elastic.DeleteResponse, error) {
	delete, err := c.Client.Delete().
		Index(c.index).
		Type(c.indexType).
		Id(file).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return delete, nil
}
