package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/uwezo-app/chat-server/db"

	"github.com/joho/godotenv"

	"github.com/uwezo-app/chat-server/router"
	"github.com/uwezo-app/chat-server/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(8000)
	}

	var dbase = db.ConnectDB()

	hub := server.NewHub()
	go hub.Run(dbase)
	r := router.Handlers(hub, dbase)

	log.Printf("Starting server on port: %v\n", port)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Could not start the server: %v\n", err)
	}
}
