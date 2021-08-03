package controller

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go

var poll []*ws.Conn

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could upgrade the connection")
	}

	defer func(conn *ws.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close the connection: %v", err)
		}
	}(conn)

	poll = append(poll, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Failure to read the message: %v", err)
		}

		for _, c := range poll {
			if c != conn {
				if err := c.WriteMessage(ws.TextMessage, msg); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
