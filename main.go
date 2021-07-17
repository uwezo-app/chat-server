package main

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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("Could upgrade the connection")
		}

		defer func(conn *ws.Conn) {
			err := conn.Close()
			if err != nil {
				log.Println("Could not close the connection")
			}
		}(conn)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
