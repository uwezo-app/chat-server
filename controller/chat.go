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
	ch := make(chan []byte)
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

	broadcaster := func(ch chan []byte) {

		for range ch {
			for _, c := range poll {
				if err := c.WriteMessage(ws.TextMessage, <-ch); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Failure to read the message: %v", err)
		}
		ch <- msg
		go broadcaster(ch)
	}
}
