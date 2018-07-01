//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
  // elastic "gopkg.in/olivere/elastic.v5"
  mpd "github.com/fhs/gompd/mpd"
)

type MpdClient struct {
  Ready chan bool
  Conn  *mpd.Client
  proto string
  addr  string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
  m := &MpdClient{
    Ready: make(chan bool),
    proto: proto,
    addr:  addr,
  }

  go m.ConnectAndKeepalive()
  return m
}

// get or refresh mpd connection
func (m *MpdClient) ConnectAndKeepalive() {
  for {
    if m.Conn != nil {
      err := m.Conn.Ping()

      if err != nil {
        fmt.Printf("MPD ping")

        m.Ready <- true
        continue

      } else {
        fmt.Printf("MPD connection broke?")

        m.Ready <- false
        m.Conn.Close()
      }
    }

    c, err := mpd.Dial(m.proto, m.addr)

    if err != nil {
      m.Conn = c

      fmt.Println("MPD (re)connect")
      m.Ready <- false
      defer(m.Conn.Close())
      continue
    }

    time.Sleep(1000 * time.Millisecond)
	}
}
