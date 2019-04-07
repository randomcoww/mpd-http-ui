//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package mpd_client

import (
	"errors"
	"time"

	mpd "github.com/fhs/gompd/mpd"
	"github.com/randomcoww/go-mpd-es/pkg/util"
	"github.com/sirupsen/logrus"
)

type MpdClient struct {
	eventHub *util.EventHub

	ApiClient   *mpd.Client
	eventClient *mpd.Client
	proto       string
	addr        string

	// Ready chan struct{}
	Events chan string
}

// create new MPD client
func NewMpdClient(proto, addr string) *MpdClient {
	c := &MpdClient{
		eventHub: util.NewEventHub(),

		proto: proto,
		addr:  addr,

		// for external use
		Events: make(chan string),
	}

	go c.run()
	go c.runEventClient()
	go c.runEventListener()

	// initial state is down
	c.eventHub.Send <- "api_down"
	c.eventHub.Send <- "event_down"

	return c
}

//
// API client
//

func (c *MpdClient) run() {
	errClient := c.eventHub.NewClient([]string{
		"api_down", "ping_down",
	})

	for {
		select {
		case event := <-errClient.Events:
			switch event {
			case "api_down":
				c.eventHub.Send <- "ping_down"

				conn, err := mpd.Dial(c.proto, c.addr)
				if err == nil {
					c.ApiClient = conn
					logrus.Infof("API ready")
					c.eventHub.Send <- "api_ready"
				} else {
					logrus.Error("Service down: %v", err)
				}

			case "ping_down":
				err := c.ApiClient.Ping()
				if err == nil {
					logrus.Infof("Ping OK")
					c.eventHub.Send <- "ping_ready"
				} else {
					logrus.Error("Ping failed")
				}
			}

		// healthcheck
		case <-time.After(2000 * time.Millisecond):
			err := c.ApiClient.Ping()

			if err != nil {
				logrus.Errorf("Ping down: %v", err)
				c.eventHub.Send <- "api_down"
			} else {
				// logrus.Infof("Ping")
			}
		}
	}
}

// lookup song metadata for elasticsearch index
// loop with reconnect attempts to make sure this happens
func (c *MpdClient) GetDatabaseItem(mpdPath string) map[string]string {
	readyClient := c.eventHub.NewClient([]string{
		"ping_ready",
	})

	for {
		attrs, err := c.ApiClient.ListInfo(mpdPath)
		if err != nil {
			c.eventHub.Send <- "api_down"
			logrus.Error("Get database item: Wait for ping")
			readyClient.WaitEvent("ping_ready")
		} else {
			if len(attrs) > 0 {
				logrus.Infof("Got MPD attrs (%d) %s\n", len(attrs), attrs[0])
				return attrs[0]
			}
			return nil
		}
	}

	return nil
}

// implement plchanges in same way as playlistinfo
func (c *MpdClient) PlChanges(version, start, end int) ([]mpd.Attrs, error) {
	var cmd *mpd.Command
	switch {
	case start < 0 && end < 0:
		// Request all playlist items.
		cmd = c.ApiClient.Command("plchanges %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.ApiClient.Command("plchanges %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.ApiClient.Command("plchanges %d %d", version, start)
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
		cmd = c.ApiClient.Command("plchangesposid %d", version)
	case start >= 0 && end >= 0:
		// Request this range of playlist items.
		cmd = c.ApiClient.Command("plchangesposid %d %d:%d", version, start, end)
	case start >= 0 && end < 0:
		// Request the single playlist item at this position.
		cmd = c.ApiClient.Command("plchangesposid %d %d", version, start)
	case start < 0 && end >= 0:
		return nil, errors.New("negative start index")
	default:
		panic("unreachable")
	}
	return cmd.AttrsList("cpos")
}

//
// Event capture client. This needs to be a separate connection from API client
//

func (c *MpdClient) runEventClient() {
	errClient := c.eventHub.NewClient([]string{
		"event_down",
	})

	for {
		select {
		case event := <-errClient.Events:
			switch event {
			case "event_down":
				conn, err := mpd.Dial(c.proto, c.addr)
				if err == nil {
					c.eventClient = conn
					logrus.Infof("API ready")
					c.eventHub.Send <- "event_ready"
				} else {
					logrus.Error("Service down")
				}
			}
		case <-time.After(2000 * time.Millisecond):
		}
	}
}

func (c *MpdClient) runEventListener() {
	readyClient := c.eventHub.NewClient([]string{
		"event_ready",
	})

	for {
		// idle will block until there is an event
		changed, err := c.eventClient.Command("idle").Strings("changed")
		if err != nil {
			c.eventHub.Send <- "event_down"
			readyClient.WaitEvent("event_ready")
		} else {
			// Got event names from MPD
			for _, e := range changed {
				c.Events <- e
			}
		}
	}
}
