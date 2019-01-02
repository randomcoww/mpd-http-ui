//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package mpd_handler

import (
	"fmt"
	"time"

	mpd "github.com/fhs/gompd/mpd"
	"github.com/sirupsen/logrus"
)

type MpdClient struct {
	up       chan struct{}
	down     chan struct{}
	pingDown chan struct{}
	Conn     *mpd.Client
	proto    string
	addr     string
	Ready    chan struct{}
}

// create new MPD client
func NewMpdClient(proto, addr string) *MpdClient {
	c := &MpdClient{
		up:       make(chan struct{}, 1),
		down:     make(chan struct{}, 1),
		pingDown: make(chan struct{}, 1),

		proto: proto,
		addr:  addr,

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

func (c *MpdClient) connect() error {
	logrus.Infof("Connecting to MPD...")
	conn, err := mpd.Dial(c.proto, c.addr)
	if err != nil {
		return err
	}

	logrus.Infof("Connected to MPD")
	// defer conn.Close()
	c.Conn = conn

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
				err := c.Conn.Ping()
				if err != nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}
			c.setState(c.up)
			c.setState(c.Ready)

		case <-time.After(10000 * time.Millisecond):
			err := c.Conn.Ping()

			if err != nil {
				logrus.Errorf("MPD ping down %s", err)
				c.setState(c.down)

			} else {
				// logrus.Infof("MPD ping")
			}
		}
	}
}

// lookup song metadata for elasticsearch index
// loop with reconnect attempts to make sure this happens
func (c *MpdClient) GetDatabaseItem(mpdPath string) map[string]string {
	for {
		attrs, err := c.Conn.ListInfo(mpdPath)

		if err != nil {
			c.drainState(c.up)
			c.setState(c.down)
			<-c.up
			continue
		}

		if len(attrs) > 0 {
			logrus.Infof("Got MPD attrs (%d) %s", len(attrs), attrs[0])
			return attrs[0]

		} else {
			return make(map[string]string)
		}
	}
}

// implement plchanges in same way as playlistinfo
func (c *MpdClient) PlChanges(version, start, end int) ([]mpd.Attrs, error) {
	var cmd *mpd.Command
	switch {
	case start < 0 && end < 0:
		// Request all playlist items.
		cmd = c.Conn.Command("plchanges %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.Conn.Command("plchanges %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.Conn.Command("plchanges %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, fmt.Errorf("negative start index")
	default:
		return nil, fmt.Errorf("unreachable")
	}
	return cmd.AttrsList("file")
}

func (c *MpdClient) PlChangePosId(version, start, end int) ([]mpd.Attrs, error) {
	var cmd *mpd.Command
	switch {
	case start < 0 && end < 0:
		// Request all playlist items.
		cmd = c.Conn.Command("plchangesposid %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.Conn.Command("plchangesposid %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.Conn.Command("plchangesposid %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, fmt.Errorf("negative start index")
	default:
		return nil, fmt.Errorf("unreachable")
	}
	return cmd.AttrsList("cpos")
}
