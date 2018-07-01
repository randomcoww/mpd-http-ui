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
  Conn  *mpd.Client
  proto string
  addr  string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
  m := &MpdClient{
    proto: proto,
    addr:  addr,
  }

  return m
}

// get or refresh mpd connection
func (m *MpdClient) mpdConn() (*mpd.Client, error) {

  if m.Conn != nil {
    err := m.Conn.Ping()

    if err != nil {
      return m.Conn, nil
    }
  }

	for {
  	c, err := mpd.Dial(m.proto, m.addr)

	  if err != nil {
	    fmt.Println("cannot connect to MPD")
			time.Sleep(1000 * time.Millisecond)

		} else {
			m.Conn = c
		  return m.Conn, nil
		}
	}
}

// main
func NewDataFeeder() (error) {
  mpdClient := NewMpdClient("tcp", "localhost:6600")
  logParser, err := NewLogEventParser("env/mpd_mount/logs/log")

  if err != nil {
    return err
  }

  for {
    select {
		case c := <- logParser.added:
			fmt.Println("add_event:", c)

      conn, err := mpdClient.mpdConn()

      if err == nil {
				attrs, err := conn.ListAllInfo(c)

				if err == nil {
					fmt.Printf("%s \n", attrs)
				} else {
					fmt.Printf("error parsing %s %s \n", c, err)
				}
      }

		case c := <- logParser.deleted:
			fmt.Println("delete_event:", c)

		case <- time.After(1000 * time.Millisecond):
    }
  }

  return nil
}
