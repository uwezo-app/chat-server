package server

import (
	"time"
)

type Message struct {
	to      *Client
	message []byte
}

// Hub maintains active connections and broadcast messages
// to connections
type Hub struct {
	// Connected clients userId => Connection
	connections map[string]*ConnectionDetail

	// Incoming messages from a client
	// to all connections subscribers of a channel
	broadcast chan []byte

	// Incoming messages for a specific client
	targeted chan *Message

	// Register client's requests
	register chan *ConnectionDetail

	// Unregister requests from the connections
	unregister chan *ConnectionDetail
}

// ConnectionDetail holds the connection info
// of a specific client
type ConnectionDetail struct {
	user string

	client *Client

	lastSeen time.Time
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*ConnectionDetail),
		broadcast:   make(chan []byte),
		register:    make(chan *ConnectionDetail),
		unregister:  make(chan *ConnectionDetail),
		targeted:    make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c.user] = c

		case c := <-h.unregister:
			if _, ok := h.connections[c.user]; ok {
				close(c.client.send)
				delete(h.connections, c.user)
			}

		case msg := <-h.broadcast:
			for c := range h.connections {
				select {
				case h.connections[c].client.send <- msg:
				default:
					close(h.connections[c].client.send)
					delete(h.connections, c)
				}
			}

		case tMessage := <-h.targeted:
			select {
			// if the channel is read to receive, send the message then
			// break out of the loop
			case tMessage.to.send <- tMessage.message:
				// the default case is when the client's channel is not ready to
				// receive, which means that they are not connected
			default:
				// This is where we could saved the message into
				// the database so that when client tMessage.to is
				// back online, we send them
				close(tMessage.to.send)
				//delete(h.connections, tMessage.to)
			}
		}
	}
}
