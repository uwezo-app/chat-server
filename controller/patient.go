package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type Patient struct {
	NickName string `json:"NickName"`
}

// CreatePatient an implementation of patient's creation
func CreatePatient(w http.ResponseWriter, r *http.Request) {
	user := &Patient{}
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
	user := &Patient{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var errorResponse = ErrorResponse {
			Code: http.StatusInternalServerError,
			Message: "Could not decode your request",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
	}
}
