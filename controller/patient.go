package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v4"

	"github.com/uwezo-app/chat-server/db"
	"github.com/uwezo-app/chat-server/server"
	"github.com/uwezo-app/chat-server/utils"
)

// CreatePatient an implementation of patient's creation
func CreatePatient(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user db.Patient
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("Could not read the sent data")
	}

	// Write the nickname to the database
	result := dbase.Save(&user)
	if er := result.Error; er != nil {
		log.Println(er)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Server error")
		return
	}

	if er := json.NewEncoder(w).Encode(user); er != nil {
		log.Println(er)
	}
}

func PatientLoginHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var patient db.Patient

	json.NewDecoder(r.Body).Decode(&patient)
	result := dbase.Where("NickName = ?", patient.NickName).First(&patient)
	if erro := result.Error; erro != nil {
		log.Println(erro)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("user not found")
		return
	}

	expiresAt := time.Now().Add(time.Hour * 168).Unix()
	claims := db.CustomClaims{
		UserID: patient.ID,
		Name:   patient.NickName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := t.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Token": tokenString,
		"User":  patient,
	})
}

func PatientLogoutHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var patient db.Patient
	tokenString := utils.GetTokenFromHeader(r.Header)
	claims, err := utils.ParseTokenWithClaims(tokenString)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}))
		return
	}

	result := dbase.Where(db.Patient{NickName: claims.Name}).First(&patient)
	if result.Error == nil {
		log.Println("User not found")
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}))
		return
	}

	// Invalidate user token
	tokenString, err = utils.GeneratePatientToken(&patient, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}))
		return
	}

	var conn server.ConnectedClient
	result = dbase.Find(&conn, server.ConnectedClient{UserID: claims.UserID})
	if result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("Your session has expired")
		return
	}

	conn.Client.Hub.Unregister <- &conn

	w.WriteHeader(http.StatusNotFound)
	log.Println(json.NewEncoder(w).Encode(tokenString))
}
