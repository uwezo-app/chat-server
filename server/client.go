package server

import (
	"bytes"
	"log"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
)

// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client stands between the hub and the ws connections
type Client struct {
	hub *Hub

	// Websockets connection
	conn *ws.Conn

	// Buffered channel for the outgoing messages
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(){
	defer func (){
		c.hub.unregister <- c
		log.Println(c.conn.Close())
	}()

	c.conn.SetReadLimit(maxMessageSize)
	log.Println(c.conn.SetReadDeadline(time.Now().Add(pongWait)))
	c.conn.SetPongHandler(func(string) error {
		log.Println(c.conn.SetReadDeadline(time.Now().Add(pongWait)))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("error %v\n", err)
			}
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.Println(c.conn.Close())
	}()

	for {
		select {
		case message, ok := <- c.send:
			log.Println(c.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if !ok {
				// Connection closed by the hub
				log.Println(c.conn.WriteMessage(ws.CloseMessage, []byte{}))
				return
			}

			w, err := c.conn.NextWriter(ws.TextMessage)
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}

			if _, err = w.Write(message); err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				log.Println(w.Write(newline))
				log.Println(w.Write(<-c.send))
			}

			if err := w.Close(); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
		case <-ticker.C:
			log.Println(c.conn.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


func ChatHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could upgrade the connection")
	}

	send := make(chan []byte)
	client := &Client{hub, conn, send}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
