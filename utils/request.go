package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, t interface{}) {
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintln("An error occurred while processing your request"),
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return
	}
}
