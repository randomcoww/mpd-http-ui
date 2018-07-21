//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package mpd_handler

import (
	"fmt"
	"time"
	"errors"
	mpd "github.com/fhs/gompd/mpd"
)

type MpdClient struct {
	up chan struct {}
	down chan struct{}
	pingDown chan struct{}
	conn *mpd.Client
	proto string
	addr string
	Ready chan struct{}
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
	c := &MpdClient{
		up: make(chan struct{}, 1),
		down: make(chan struct{}, 1),
		pingDown: make(chan struct{}, 1),

		proto: proto,
		addr: addr,

		// for external use
		Ready: make(chan struct{}, 1),
	}

	c.setState(c.down)
	go c.reconnectLoop()

	return c
}


func (c *MpdClient) setState(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (c *MpdClient) drainState(ch chan struct{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
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

		case <-c.down:
			for {
				err := c.connect()
				if err != nil {
					time.Sleep(2000 * time.Millisecond)
					continue
				}
				break
			}
			c.setState(c.pingDown)

		case <-c.pingDown:
			for {
				err := c.conn.Ping()
				if err != nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}
			c.setState(c.up)
			c.setState(c.Ready)

		case <-time.After(10000 * time.Millisecond):
			err := c.conn.Ping()

			if err != nil {
				fmt.Printf("MPD ping down %s\n", err)
				c.setState(c.down)

			} else {
				// fmt.Printf("MPD ping\n")
			}
		}
	}
}


// lookup song metadata for elasticsearch index
// loop with reconnect attempts to make sure this happens
func (c *MpdClient) GetDatabaseItem(mpdPath string) (map[string]string) {
	for {
		attrs, err := c.conn.ListInfo(mpdPath)

		if err != nil {
			c.drainState(c.up)
			c.setState(c.down)
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


// manipulate playlist
// query current playlist items between position start and end
func (c *MpdClient) QueryPlaylistItems(start, end int) ([]mpd.Attrs, error) {
	attrs, err := c.conn.PlaylistInfo(start, end)
	return attrs, err
}

//
func (c *MpdClient) CurrentSong() (mpd.Attrs, error) {
	attrs, err := c.conn.CurrentSong()
	return attrs, err
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

func (c *MpdClient) Status() (mpd.Attrs, error) {
	attrs, err := c.conn.Status()
	return attrs, err
}

// implement plchanges in same way as playlistinfo
func (c *MpdClient) PlChanges(version, start, end int) ([]mpd.Attrs, error) {
	var cmd *mpd.Command
	switch {
	case start < 0 && end < 0:
		// Request all playlist items.
		cmd = c.conn.Command("plchanges %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.conn.Command("plchanges %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.conn.Command("plchanges %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, errors.New("negative start index")
	default:
		panic("unreachable")
	}
	return cmd.AttrsList("file")
}

func (c *MpdClient) PlChangePosId(version, start, end int) ([]mpd.Attrs, error) {
	var cmd *mpd.Command
	switch {
	case start < 0 && end < 0:
		// Request all playlist items.
		cmd = c.conn.Command("plchangesposid %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.conn.Command("plchangesposid %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.conn.Command("plchangesposid %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, errors.New("negative start index")
	default:
		panic("unreachable")
	}
	return cmd.AttrsList("cpos")
}
