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
	up chan struct {}
	down chan struct{}

	conn *mpd.Client
	proto string
	addr string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
	c := &MpdClient{
		up: make(chan struct{}, 1),
		down: make(chan struct{}, 1),

		proto: proto,
		addr: addr,
	}

	c.down <- struct{}{}
	go c.reconnectLoop()

	return c
}


func (c *MpdClient) connect() (error) {
	if c.conn != nil {
		err := c.conn.Ping()

		if err != nil {
			fmt.Printf("Reconnecting MPD...\n")
			// c.conn.Close()

		} else {
			fmt.Printf("MPD connection still alive\n")
			return nil
		}
	}

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

		case <-c.down:
			for {
				err := c.connect()
				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}
				break
			}
			c.up <- struct{}{}
		}
	}
}


func (c *MpdClient) GetDatabaseItem(mpdPath string) (map[string]string) {
	for {
		attrs, err := c.conn.ListInfo(mpdPath)

		if err != nil {
			time.Sleep(2000 * time.Millisecond)

			c.down <- struct{}{}
			<-c.up

			continue
		}

		if len(attrs) > 0 {
			fmt.Printf("Got MPD attrs (%d) %s\n", len(attrs), attrs[0])
			return attrs[0]

		} else {
			return make(map[string]string)
		}
	}
}
