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
  Watch chan struct{}
  Down chan struct{}
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
    Watch: make(chan struct{}, 1),
    Down: make(chan struct{}, 1),
    proto: proto,
    addr: addr,
    Events: make(chan string),
  }

  c.setStatusDown()
  go c.reconnectLoop()
  go c.setupWatcher()

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
  <-c.Watch

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


func (c *MpdClient) GetInfo(mpdPath string) (map[string]string) {
  for {
    attrs, err := c.Conn.ListInfo(mpdPath)

    if err == nil {
      if len(attrs) > 0 {
        fmt.Printf("Got MPD attrs (%d) %s\n", len(attrs), attrs[0])
        return attrs[0]

      } else {
        fmt.Printf("Got MPD empty attrs\n")
        return make(map[string]string)
      }

    } else {
      c.setStatusDown()
      fmt.Printf("Start MPD ready wait\n")
      <-c.Ready
    }
  }
}
