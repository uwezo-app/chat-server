package server

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/uwezo-app/chat-server/db"
)

type Message struct {
	to, from       *Client
	conversationID uint
	message        []byte
}

// Hub maintains active connections and broadcast messages
// to connections
type Hub struct {
	Connections map[uint]*ConnectedClient

	// Incoming messages from a client
	// to all connections subscribers of a channel
	Broadcast chan []byte

	// Incoming messages for a specific client
	Targeted chan *Message

	// returns connected users
	GetUsers chan *Client

	// Register client's requests
	Register chan *ConnectedClient

	// connects two peers
	Pair chan *db.PairedUsers

	// Unregister requests from the connections
	Unregister chan *ConnectedClient
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:   make(chan []byte),
		Targeted:    make(chan *Message),
		Register:    make(chan *ConnectedClient),
		Pair:        make(chan *db.PairedUsers),
		Unregister:  make(chan *ConnectedClient),
		Connections: make(map[uint]*ConnectedClient),
		GetUsers:    make(chan *Client),
	}
}

func (h *Hub) Run(dbase *gorm.DB) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case c := <-h.Register:
			conn := ConnectedClient{
				UserID:   c.UserID,
				Client:   c.Client,
				LastSeen: c.LastSeen,
			}
			h.Connections[c.UserID] = &conn

		case c := <-h.Unregister:
			if _, ok := h.Connections[c.UserID]; ok {
				close(c.Client.Send)
				close(c.Client.Notify)
				_ = c.Client.Conn.Close()
			}

		case <-h.Broadcast:
		//	for c := range h.connections {
		//		select {
		//		case h.connections[c].Client.send <- msg:
		//		default:
		//			close(h.connections[c].Client.send)
		//			delete(h.connections, c)
		//		}
		//	}

		case c := <-h.GetUsers:
			c.SendJSON <- &struct {
				Users []interface{} `json:"users"`
			}{
				Users: func() []interface{} {
					type user struct {
						ID   uint   `json:"id"`
						Name string `json:"name"`
					}
					users := make([]interface{}, 0)
					for _, c := range h.Connections {
						user_ := &user{
							ID:   c.UserID,
							Name: c.Client.Name,
						}
						users = append(users, user_)
					}
					return users
				}(),
			}

		case pairReq := <-h.Pair:
			res := dbase.Create(&pairReq)
			if res.Error != nil {
				return
			}

			// Notify the users of the connection
			patient := h.Connections[pairReq.PatientID]
			psy := h.Connections[pairReq.PsychologistID]
			patient.Client.Send <- []byte(fmt.Sprintf("ConversationID %v", pairReq.ID))
			psy.Client.Send <- []byte(fmt.Sprintf("ConversationID %v", pairReq.ID))

		case tMessage := <-h.Targeted:
			select {
			// if the channel is read to receive, send the message then
			// break out of the loop
			case tMessage.to.Send <- tMessage.message:
				dbase.Create(&db.Conversation{
					ConversationID: tMessage.conversationID,
					From:           tMessage.from.ClientID,
					Message:        string(tMessage.message),
					SentAt:         time.Now(),
				})
				tMessage.from.Send <- tMessage.message
			// the default case is when the client's channel is not ready to
			// receive, which means that they are not connected
			default:
				// This is where we could save the message into
				// the database so that when client tMessage.to is
				// back online, we send them
				close(tMessage.to.Send)
				//delete(h.connections, tMessage.to)
			}
		}
	}
}
