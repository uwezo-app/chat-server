package controller

import (
	ws "github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	go func() {
		ct := time.Tick(5 * time.Second)

		for range ct {
			for range ch {
				for _, c := range poll {
					if err := c.WriteMessage(ws.TextMessage, <-ch); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Failure to read the message: %v", err)
		}
		ch <- msg
	}
}
