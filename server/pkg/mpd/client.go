//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package mpd

import (
	"errors"
	"time"

	mpd "github.com/fhs/gompd/mpd"
	"github.com/randomcoww/go-mpd-es/pkg/util"
	"github.com/sirupsen/logrus"
)

type MpdClient struct {
	eventHub *util.EventHub

	Conn  *mpd.Client
	proto string
	addr  string
}

// create new MPD client
func NewMpdClient(proto, addr string) *MpdClient {

	logrus.Infof("MpdClient: Start")

	c := &MpdClient{
		eventHub: util.NewEventHub(),
		proto:    proto,
		addr:     addr,
	}

	c.setReady()
	go c.runRecovery()

	return c
}

//
// get connection
//

func (c *MpdClient) pingTest() bool {
	err := c.Conn.Ping()
	return err == nil
}

func (c *MpdClient) waitConnect() {
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			conn, err := mpd.Dial(c.proto, c.addr)

			if err == nil {
				c.Conn = conn

				logrus.Infof("MpdClient: Connection ready")
				return
			}
		}
	}
}

func (c *MpdClient) waitPingState(state bool) {
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			if c.pingTest() == state {
				logrus.Infof("MpdClient: Ping state changed: ", state)
				return
			}
		}
	}
}

func (c *MpdClient) setReady() {
	c.waitConnect()
	c.waitPingState(true)
	c.eventHub.Send <- "api_ready"
}

func (c *MpdClient) setDown() {
	c.waitPingState(false)
	c.eventHub.Send <- "api_down"
}

func (c *MpdClient) runRecovery() {
	errClient := c.eventHub.NewClient([]string{"api_down"})
	readyClient := c.eventHub.NewClient([]string{"api_ready"})

	for {
		select {
		case <-errClient.Events:
			c.setReady()
			errClient.Drain()

		case <-readyClient.Events:
			c.setDown()
			readyClient.Drain()
		}
	}
}

// lookup song metadata for elasticsearch index
// loop with reconnect attempts to make sure this happens
func (c *MpdClient) GetDatabaseItem(mpdPath string) map[string]string {
	readyClient := c.eventHub.NewClient([]string{"api_ready"})

	for {
		if attrs, err := c.Conn.ListInfo(mpdPath); err == nil {
			if len(attrs) > 0 {
				return attrs[0]
			}
			return nil
		}

		c.eventHub.Send <- "api_down"
		readyClient.WaitEvent("api_ready")
	}

	return nil
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
		cmd = c.Conn.Command("plchangesposid %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.Conn.Command("plchangesposid %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.Conn.Command("plchangesposid %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, errors.New("negative start index")
	default:
		panic("unreachable")
	}
	return cmd.AttrsList("cpos")
}
