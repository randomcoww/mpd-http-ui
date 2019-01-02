// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan *socketMessage
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan *socketMessage
}

//
// web socket feeder
// based on example https://github.com/gorilla/websocket/blob/master/examples/filewatch/main.go
//
func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logrus.Errorf("%s", err)
		}

		logrus.Infof("Error serving WS: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: ws,
		send: make(chan *socketMessage, 256),
	}

	hub.register <- client
	go client.writeSocketEvents()
	go client.readSocketEvents()
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *socketMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

//
// Add stuff to original file
//

//
// broadcaster
//
func (h *Hub) eventBroadcaster() {
	for {
		select {
		case e := <-mpdEvent.Event:
			logrus.Infof("Got MPD event: %s", e)

			switch e {
			case "player":
				msg, err := createStatusMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

				msg, err = createCurrentSongMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

				msg, err = createSeekMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "playlist":
				msg, err := createPlaylistChangedMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "mixer", "options", "outputs":
				msg, err := createStatusMessage()
				if err != nil {
					break
				}
				h.broadcast <- msg

			case "update":
				msg := createUpdateDatabaseMessage()
				h.broadcast <- msg
			}

		case <-time.After(1000 * time.Millisecond):
			msg, err := createSeekMessage()
			if err != nil {
				break
			}
			if msg != nil {
				h.broadcast <- msg
			}
		}
	}
}

//
// send broadcast events to each client
//
func (c *Client) writeSocketEvents() {

	defer func() {
		logrus.Infof("Close writer")
		c.conn.Close()
	}()

	defer func() {
		r := recover()
		if r != nil {
			logrus.Infof("Recovered %s", r)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				logrus.Infof("Hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(*msg)
		}
	}
}

//
// read messages from client
//
func (c *Client) readSocketEvents() {

	defer func() {
		logrus.Infof("Close reader")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	defer func() {
		r := recover()
		if r != nil {
			time.Sleep(1000 * time.Millisecond)
			logrus.Infof("Recovered %s", r)
		}
	}()

	for {
		v := &socketMessage{}

		err := c.conn.ReadJSON(v)
		if err != nil {
			logrus.Errorf("Error reading socket %s", err)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("error: %v", err)
			}
			break
		}

		switch v.Name {
		case "seek":
			t := int64(v.Data.(float64) * 1000000000)
			mpdClient.Conn.SeekCur(time.Duration(t), false)

			// client specific playlist query
		case "playlistquery":
			d := v.Data.([]interface{})
			start := int(d[0].(float64))
			end := int(d[1].(float64))
			msg, err := createPlaylistQueryMessage(start, end)
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

			// client specific current song query
		case "currentsong":
			msg, err := createCurrentSongMessage()
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

			// global playlist items moved
			// send only and allow server to emit event
		case "playlistmove":
			d := v.Data.([]interface{})
			start := int(d[0].(float64))
			end := int(d[1].(float64))
			position := int(d[2].(float64))
			if start != position {
				mpdClient.Conn.Move(start, end, position)
			}

		case "playid":
			// -1 for play current
			d := int(v.Data.(float64))
			mpdClient.Conn.PlayID(d)

		case "stop":
			mpdClient.Conn.Stop()

		case "pause":
			mpdClient.Conn.Pause(true)

		case "playnext":
			mpdClient.Conn.Next()

		case "playprev":
			mpdClient.Conn.Previous()

		case "removeid":
			d := int(v.Data.(float64))
			mpdClient.Conn.DeleteID(d)

		case "addpath":
			d := v.Data.([]interface{})
			path := d[0].(string)
			position := int(d[1].(float64))
			mpdClient.Conn.AddID(path, position)

			// client specific database search
		case "search":
			d := v.Data.([]interface{})
			query := d[0].(string)
			start := int(d[1].(float64))
			size := int(d[2].(float64))
			msg, err := createSearchMessage(query, start, size)
			if err != nil {
				break
			}
			c.conn.WriteJSON(*msg)

		case "clear":
			mpdClient.Conn.Clear()

		case "updatedb":
			mpdClient.Conn.Update("")
		}
	}
}
