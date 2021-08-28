package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/uwezo-app/chat-server/db"
	"github.com/uwezo-app/chat-server/server"
	"gorm.io/gorm"
)

func GetConversation(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	conv := db.PairedUsers{}
	psy, _ := strconv.Atoi(query.Get("PsychologistID"))
	pat, _ := strconv.Atoi(query.Get("PatientID"))

	result := dbase.Where(db.PairedUsers{PsychologistID: uint(psy), PatientID: uint(pat)}).Find(&conv)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(result.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conv)
}


func PairUsers(hub *server.Hub, dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var pair db.PairedUsers
	err := json.NewDecoder(r.Body).Decode(&pair)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	_, ok := hub.Connections[pair.PsychologistID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Psychologist not found")
		return
	}

	result := dbase.Create(&pair)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result.Error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pair)
}
