package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	hub := server.NewHub()
	go hub.Run()
	r := router.Handlers(hub)

	log.Printf("%v Starting server\n", time.Now())
	if err = http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Could not start the server: %v\n", err)
	}
}
