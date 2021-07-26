package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/uwezo-app/chat-server/router"

	"github.com/joho/godotenv"
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
	r := router.Handlers()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
