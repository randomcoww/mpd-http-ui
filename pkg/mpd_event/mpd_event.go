// capture and send event details

package mpd_event

import (
	"time"

	mpd "github.com/fhs/gompd/mpd"
	"github.com/sirupsen/logrus"
)

type MpdEvent struct {
	up    chan struct{}
	down  chan struct{}
	conn  *mpd.Client
	proto string
	addr  string

	// Ready chan struct{}
	Event chan string
}

// create new MPD client
func NewEventWatcher(proto, addr string) *MpdEvent {
	c := &MpdEvent{
		up:   make(chan struct{}, 1),
		down: make(chan struct{}, 1),

		proto: proto,
		addr:  addr,

		// for external use
		Event: make(chan string),
	}

	c.setState(c.down)
	go c.reconnectLoop()

	return c
}

func (c *MpdEvent) setState(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (c *MpdEvent) drainState(ch chan struct{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

func (c *MpdEvent) connect() error {
	logrus.Infof("Connecting to MPD...")
	conn, err := mpd.Dial(c.proto, c.addr)

	if err != nil {
		return err
	}

	logrus.Infof("Connected to MPD")
	// defer conn.Close()
	c.conn = conn

	return nil
}

func (c *MpdEvent) reconnectLoop() {
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
			c.setState(c.up)

		case <-c.up:
			changed, err := c.conn.Command("idle").Strings("changed")
			if err != nil {
				logrus.Errorf("MPD event watcher error %s", err)
				c.setState(c.down)

			} else {
				c.setState(c.up)
				for _, e := range changed {
					c.Event <- e
				}
			}
		}
	}
}
