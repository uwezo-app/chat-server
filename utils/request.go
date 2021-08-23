package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, t interface{}) {
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintln("An error occurred while processing your request"),
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return
	}
}

func HashPassword(pass string, w http.ResponseWriter) string {
	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return ""
	}

	return string(password)
}
