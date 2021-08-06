package server

import (
	"fmt"
	"github.com/uwezo-app/chat-server/db"
	"gorm.io/gorm"
	"log"
	"time"
)

type Message struct {
	to, from *Client
	conversationID uint
	message  []byte
}

// Hub maintains active connections and broadcast messages
// to connections
type Hub struct {
	// Incoming messages from a client
	// to all connections subscribers of a channel
	broadcast chan []byte

	// Incoming messages for a specific client
	targeted chan *Message

	// Register client's requests
	register chan *db.ConnectedClient

	// connects two peers
	pair chan *db.PairedUsers

	// Unregister requests from the connections
	unregister chan *db.ConnectedClient
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		targeted:   make(chan *Message),
		register:   make(chan *db.ConnectedClient),
		pair:       make(chan *db.PairedUsers),
		unregister: make(chan *db.ConnectedClient),
	}
}

func (h *Hub) Run(dbase *gorm.DB) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case c := <-h.register:
			conn := db.ConnectedClient{
				UserID:   c.UserID,
				Client:   c.Client,
				LastSeen: c.LastSeen,
			}
			res := dbase.Create(&conn)
			if res.Error != nil {
				c.Client.notify <- []byte("connection closed")
				log.Println(c.Client.conn.Close())
				return
			}

		case c := <-h.unregister:
			res := dbase.Where(&db.ConnectedClient{UserID: c.UserID}).Delete(&db.ConnectedClient{})
			if res.Error != nil {
				c.Client.notify <- []byte("could not close connection")
				return
			}

			close(c.Client.send)
			close(c.Client.notify)
			log.Println(c.Client.conn.Close())

		case _ = <-h.broadcast:
		//	for c := range h.connections {
		//		select {
		//		case h.connections[c].Client.send <- msg:
		//		default:
		//			close(h.connections[c].Client.send)
		//			delete(h.connections, c)
		//		}
		//	}

		case pairReq := <-h.pair:
			res := dbase.Create(pairReq)
			if res.Error != nil {
				return
			}

			// Notify the users of the connection
			var patient, psy *db.ConnectedClient
			dbase.Find(&patient, &db.ConnectedClient{UserID: pairReq.PatientID})
			dbase.Find(&psy, &db.ConnectedClient{UserID: pairReq.PsychologistID})
			patient.Client.notify <- []byte(fmt.Sprintf("Connected to %v", psy.UserID))
			psy.Client.notify <- []byte(fmt.Sprintf("Connected to %v", patient.UserID))

		case tMessage := <-h.targeted:
			select {
			// if the channel is read to receive, send the message then
			// break out of the loop
			case tMessage.to.send <- tMessage.message:
				dbase.Create(db.Conversation{
					ConversationID: tMessage.conversationID,
					From:           tMessage.from.ClientID,
					Message:        string(tMessage.message),
					SentAt:         time.Now(),
				})
				tMessage.from.send <- tMessage.message
			// the default case is when the client's channel is not ready to
			// receive, which means that they are not connected
			default:
				// This is where we could save the message into
				// the database so that when client tMessage.to is
				// back online, we send them
				close(tMessage.to.send)
				//delete(h.connections, tMessage.to)
			}
		}
	}
}
