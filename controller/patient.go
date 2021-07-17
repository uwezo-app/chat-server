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
	patient := &Patient{}
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		log.Fatal("Could not read the sent data")
	}


}
