package main

import (
	"fmt"
	"github.com/uwezo-app/chat-server/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	r := router.Handlers()
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
