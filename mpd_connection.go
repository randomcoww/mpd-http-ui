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
  Ready chan struct{}
  Down chan struct{}
  Conn *mpd.Client
  proto string
  addr string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
  c := &MpdClient{
    Ready: make(chan struct{}, 1),
    Down: make(chan struct{}, 1),
    proto: proto,
    addr: addr,
  }

  c.setStatusDown()
  go c.reconnectLoop()

  return c
}


func (c *MpdClient) setStatusReady() {
  c.Ready <- struct{}{}
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


func (c *MpdClient) GetInfo(mpdPath string) (map[string]string) {
  for {
    attrs, err := c.Conn.ListAllInfo(mpdPath)

    if err == nil {
      if len(attrs) > 0 {
  fmt.Printf("Get MPD attrs %s\n", attrs[0])

        return attrs[0]
      } else {
        return make(map[string]string)
      }

    } else {
      c.setStatusDown()
      fmt.Printf("Start MPD ready wait\n")
      <-c.Ready
    }
  }
}
