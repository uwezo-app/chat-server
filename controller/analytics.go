package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/uwezo-app/chat-server/db"
)

func GetNumberofPsychologists(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var psychologists []db.Psychologist
	dbase.Find(&psychologists)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(len(psychologists)))
}

func GetNumberofPatients(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var patients []db.Patient
	dbase.Find(&patients)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(len(patients)))
}

func Get5LatestPsychologists(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var psychologists []db.Psychologist
	dbase.Order("created_at desc").Find(&psychologists)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(psychologists))
}

func GetMonthlyActiveUsers(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {

}
