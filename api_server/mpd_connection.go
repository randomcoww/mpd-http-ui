//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
	"strings"
	mpd "github.com/fhs/gompd/mpd"
)

type MpdClient struct {
	Ready chan struct{}
	Down chan struct{}
	Watch chan struct{}
	Conn *mpd.Client
	proto string
	addr string
	Events chan string
}

var (
	watchEvents = []string{
		"database",
		"update",
		"stored_playlist",
		"playlist",
		"mixer",
		"output",
		"options",
		// "partition",
		"sticker",
		"subscription",
		"message",
	}
)

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
	c := &MpdClient{
		Ready: make(chan struct{}, 1),
		Down: make(chan struct{}, 1),
		Watch: make(chan struct{}, 1),
		proto: proto,
		addr: addr,
		Events: make(chan string),
	}

	c.setStatusDown()
	go c.reconnectLoop()

	return c
}


func (c *MpdClient) setStatusReady() {
	c.Ready <- struct{}{}
	c.Watch <- struct{}{}
	fmt.Printf("MPD ready\n")
}

func (c *MpdClient) setStatusDown() {
	c.Down <- struct{}{}
	fmt.Printf("MPD down\n")
}


func (c *MpdClient) reconnectLoop() {
	for {
		select {

		case <-c.Down:
			for {
				time.Sleep(1000 * time.Millisecond)

				fmt.Printf("Connecting to MPD...\n")
				conn, err := mpd.Dial(c.proto, c.addr)
				defer conn.Close()

				if err == nil {
					c.Conn = conn

					c.setStatusReady()
					break

				} else {
					fmt.Printf("Error connecting to MPD\n")
				}
			}
		}
	}
}


// reimplement watch included in log watch
func (c *MpdClient) setupWatcher() {
	<-c.Ready

	for {
		changed, err := c.Conn.
			Command("idle %s", mpd.Quoted(strings.Join(watchEvents, " "))).
			Strings("changed")

		if err == nil {
			fmt.Printf("MPD event add: %s\n", changed)

			for _, e := range changed {
				c.Events <-e
			}

		} else {
			c.setStatusDown()
			fmt.Printf("Start MPD ready wait\n")
			<-c.Watch
		}
	}
}

// query current playlist items between position start and end
func  (c *MpdClient) QueryPlaylistItems(start, end int) ([]mpd.Attrs, error) {
	attrs, err := c.Conn.PlaylistInfo(start, end)
	return attrs, err
}

// add database item to current playlist at position
func (c *MpdClient) AddToPlaylist(mpdPath string, position int) (int, error) {
	id, err := c.Conn.AddID(mpdPath, position)
	return id, err
}

// moves songs in current playlist between positions start and end to new position position
func (c *MpdClient) MovePlaylistItems(start, end, newPosition int) (error) {
	err := c.Conn.Move(start, end, newPosition)
	return err
}

// deletes playlist items between positions start and end
func (c *MpdClient) DeletePlaylistItem(start, end int) (error) {
	err := c.Conn.Delete(start, end)
	return err
}

// clear current playlist
func (c *MpdClient) ClearPlaylist() (error) {
	err := c.Conn.Clear()
	return err
}
