package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/uwezo-app/chat-server/db"
	"github.com/uwezo-app/chat-server/utils"
)

// CreatePatient an implementation of patient's creation
func CreatePatient(w http.ResponseWriter, r *http.Request) {
	user := &db.Patient{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Fatal("Could not read the sent data")
	}

	// Write the nickname to the database

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatal(err)
	}
}

func AddProfileInformation(w http.ResponseWriter, r *http.Request) {}

func PatientLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &db.Patient{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var errorResponse = utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Could not decode your request",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
	}
}
