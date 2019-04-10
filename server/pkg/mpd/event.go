//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package mpd

import (
	"time"

	mpd "github.com/fhs/gompd/mpd"
	"github.com/randomcoww/go-mpd-es/pkg/util"
	"github.com/sirupsen/logrus"
)

type MpdEvent struct {
	eventHub *util.EventHub

	conn  *mpd.Client
	proto string
	addr  string
	// Ready chan struct{}
	Events chan string
}

// create new MPD client
func NewMpdEvent(proto, addr string) *MpdEvent {

	logrus.Infof("MpdEvent: Start")

	c := &MpdEvent{
		eventHub: util.NewEventHub(),
		proto:    proto,
		addr:     addr,
		Events:   make(chan string),
	}

	c.setReady()
	go c.runRecovery()
	go c.runEventListener()

	return c
}

//
// get connection
//

func (c *MpdEvent) setReady() {
	c.waitConnect()
	c.eventHub.Send <- "api_ready"
}

func (c *MpdEvent) waitConnect() {
	for {
		select {
		case <-time.After(1000 * time.Millisecond):
			conn, err := mpd.Dial(c.proto, c.addr)

			if err == nil {
				c.conn = conn

				logrus.Infof("MpdEvent: Connection ready")
				return
			}
		}
	}
}

func (c *MpdEvent) runRecovery() {
	errClient := c.eventHub.NewClient([]string{"api_ready"})

	for {
		select {
		case <-errClient.Events:
			c.setReady()
			errClient.Drain()
		}
	}
}

func (c *MpdEvent) runEventListener() {
	readyClient := c.eventHub.NewClient([]string{"api_ready"})

	for {
		changed, err := c.conn.Command("idle").Strings("changed")
		if err == nil {
			for _, e := range changed {
				logrus.Infof("MpdEvent: Event: ", e)
				c.Events <- e
			}
		} else {
			c.eventHub.Send <- "api_down"
			readyClient.WaitEvent("api_ready")
		}
	}
}
