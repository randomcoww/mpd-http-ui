//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
	mpd "github.com/fhs/gompd/mpd"
)

type MpdClient struct {
	Ready chan struct{}
	Down chan struct{}
	conn *mpd.Client
	proto string
	addr string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
	c := &MpdClient{
		Ready: make(chan struct{}, 1),
		Down: make(chan struct{}, 1),
		proto: proto,
		addr: addr,
	}

	c.Down <- struct{}{}
	go c.reconnectLoop()

	return c
}


func (c *MpdClient) connect() (error) {
	fmt.Printf("Connecting to MPD...\n")
	conn, err := mpd.Dial(c.proto, c.addr)

	if err != nil {
		return err
	}

	fmt.Printf("Connected to MPD\n")

	// defer conn.Close()
	c.conn = conn

	return nil
}


func (c *MpdClient) reconnectLoop() {
	for {
		select {

		case <-c.Down:
			for {
				err := c.connect()

				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}

				c.Ready <- struct{}{}
				break
			}
		}
	}
}


func (c *MpdClient) GetDatabaseItem(mpdPath string) (map[string]string) {
	for {
		attrs, err := c.conn.ListInfo(mpdPath)

		if err == nil {
			if len(attrs) > 0 {
				fmt.Printf("Got MPD attrs (%d) %s\n", len(attrs), attrs[0])
				return attrs[0]

			} else {
				fmt.Printf("Got MPD empty attrs\n")
				return make(map[string]string)
			}

		} else {
			c.Down <- struct{}{}
			fmt.Printf("Start MPD ready wait\n")
			<-c.Ready
		}
	}
}
