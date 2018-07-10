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
	watch chan struct{}
	up chan struct {}
	down chan struct{}

	conn *mpd.Client
	proto string
	addr string
	eventChanges chan []string
	events chan string
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
		watch: make(chan struct{}, 1),
		up: make(chan struct{}, 1),
		down: make(chan struct{}, 1),

		proto: proto,
		addr: addr,
		eventChanges: make(chan []string),
		events: make(chan string),
	}

	c.down <- struct{}{}
	go c.reconnectLoop()
	go c.setupWatcher()

	return c
}


func (c *MpdClient) connect() (error) {
	if c.conn != nil {
		err := c.conn.Ping()

		if err != nil {
			fmt.Printf("Ping MPD failed\n")
			// c.conn.Close()

		} else {
			fmt.Printf("Ping MPD\n")
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

				c.up <- struct{}{}
				break
			}

		case <-time.After(1000 * time.Millisecond):
			err := c.conn.Ping()

			if err != nil {
				c.down <- struct{}{}
			}
		}
	}
}


// reimplement watch included in log watch
func (c *MpdClient) setupWatcher() {
	<-c.watch

	for {
		changed, err := c.conn.
			Command("idle %s", mpd.Quoted(strings.Join(watchEvents, " "))).
			Strings("changed")

		if err != nil {
			time.Sleep(2000 * time.Millisecond)

			c.down <- struct{}{}
			<-c.up

			continue
		}

		select {
		case c.eventChanges <- changed:
		default:
		}

		select {
		case changed := <-c.eventChanges:
			for _, e := range changed {
				c.events <-e
			}
		default:
		}
	}
}

// manipulate playlist

// query current playlist items between position start and end
func  (c *MpdClient) QueryPlaylistItems(start, end int) ([]mpd.Attrs, error) {
	attrs, err := c.conn.PlaylistInfo(start, end)

	if err != nil {
		fmt.Printf("MPD error getting playlist: %s\n", err)
		return nil, err
	}

	fmt.Printf("MPD got playlist: %s\n", attrs)
	return attrs, nil
}

// add database item to current playlist
func (c *MpdClient) AddToPlaylist(mpdPath string) (error) {
	err := c.conn.Add(mpdPath)
	return err
}

// moves songs in current playlist between positions start and end to new position position
func (c *MpdClient) MovePlaylistItems(start, end, newPosition int) (error) {
	err := c.conn.Move(start, end, newPosition)
	return err
}

// deletes playlist items between positions start and end
func (c *MpdClient) DeletePlaylistItems(start, end int) (error) {
	err := c.conn.Delete(start, end)
	return err
}

// clear current playlist
func (c *MpdClient) ClearPlaylist() (error) {
	err := c.conn.Clear()
	return err
}

// play/pause/stop

// start playing
func (c *MpdClient) PlayItem(position int) (error) {
	err := c.conn.Play(position)
	return err
}

// stop playing
func (c *MpdClient) Stop() (error) {
	err := c.conn.Stop()
	return err
}

// pause playing
func (c *MpdClient) Pause() (error) {
	err := c.conn.Pause(true)
	return err
}
