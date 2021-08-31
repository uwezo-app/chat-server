package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/uwezo-app/chat-server/db"
	"gorm.io/gorm"
)

type Item struct {
	ID             uint      `json:"ID"`
	Name           string    `json:"Name"`
	PairID         uint      `json:"PairID"`
	PairEncryption string    `json:"PairEncryption"`
	PairedAt       time.Time `json:"PairedAt"`
}

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

func PairUsers(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var pair db.PairedUsers
	err := json.NewDecoder(r.Body).Decode(&pair)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
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

/**
1. Choose to shake the world for God
2. know how to co-operate with the Holy Spirit
3. Believe in miracle
4. Depend on the power of the Holy Spirit
5. Priase Jesus => They key to the glory store is Praise to Jesus
6. If you are not called in to ministry by God, don't do it
7. You don't have to understand just obey
8. Give your heart
9. God is absolute perfection
10.Have a conviction about God's call
*/
func GetConnections(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var queries = r.URL.Query()
	var rows *sql.Rows
	var err error

	if psychologistID := queries.Get("psychologistID"); psychologistID != "" {
		p, _ := strconv.Atoi(psychologistID)
		rows, err = dbase.
			Table("paired_users").
			Where("paired_users.psychologist_id = ?", p).
			Joins("join patients on patients.id = paired_users.patient_id").
			Select("patients.id, patients.nick_name, paired_users.id, paired_users.encryption_key, paired_users.paired_at").
			Rows()
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error")
		return
	}

	defer rows.Close()

	var response = []interface{}{}

	for rows.Next() {
		item := Item{}

		err = rows.Scan(&item.ID, &item.Name, &item.PairID, &item.PairEncryption, &item.PairedAt)
		if err != nil {
			log.Println(err)
		}

		response = append(response, item)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPatientConnections(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var queries = r.URL.Query()
	var rows *sql.Rows
	var err error

	if patientID := queries.Get("patientID"); patientID != "" {
		p, _ := strconv.Atoi(patientID)
		rows, _ = dbase.
			Table("paired_users").
			Where("paired_users.patient_id = ?", p).
			Joins("join psychologists on psychologists.id = paired_users.psychologist_id").
			Select("psychologists.id, psychologists.first_name, psychologists.last_name, paired_users.id, paired_users.encryption_key, paired_users.paired_at").
			Rows()
	}

	defer rows.Close()

	var response = []interface{}{}
	for rows.Next() {
		item := Item{}
		var fn string
		var ln string

		err = rows.Scan(&item.ID, &fn, &ln, &item.PairID, &item.PairEncryption, &item.PairedAt)
		if err != nil {
			log.Println(err)
		}

		item.Name = fmt.Sprintf("%s %s", fn, ln)
		response = append(response, item)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func NewChat(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var chat db.Conversation
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	result := dbase.Create(&chat)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result.Error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat.Message)
}

func GetChats(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var queries = r.URL.Query()
	var chat = []db.Conversation{}
	cnvId, _ := strconv.Atoi(queries.Get("ConversationID"))

	res := dbase.Where(db.Conversation{ConversationID: uint(cnvId)}).Find(&chat)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}
