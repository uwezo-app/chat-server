package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/uwezo-app/chat-server/db"
	"github.com/uwezo-app/chat-server/server"
	"github.com/uwezo-app/chat-server/utils"
)

// https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8

// CreatePsychologist implements psychologist creation
func CreatePsychologist(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := &db.Psychologist{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}

	var password []byte
	password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return
	}

	user.Password = string(password)
	user.Profile = db.Profile{
		ID: 0,
	}

	rs := dbase.Create(&user)
	if rs.Error != nil {
		log.Println(rs)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Could not create your account. Please try again later",
		}))
	}

	var writer bytes.Buffer
	body := struct {
		Name string
		Link string
	}{
		Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Link: "https://google.com",
	}

	go func(dbase *gorm.DB, email string, HTMLtemp string, body interface{}, writer *bytes.Buffer) {
		err := utils.SendEmail(dbase, email, "templates/email/reset.html", body, writer)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}(dbase, user.Email, "templates/email/confirm.html", &body, &writer)

	log.Println(json.NewEncoder(w).Encode(user))
}

// LoginHandler implements authentication for psychologists
func LoginHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user = &db.Psychologist{}
	var resp = make(map[string]interface{})
	var err error

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := utils.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintln("An error occurred while processing your request"),
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return
	}

	resp, err = FindOne(dbase, user.Email, user.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}))
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(resp))
}

func FindOne(dbase *gorm.DB, email, password string) (map[string]interface{}, error) {
	var user *db.Psychologist

	dbase.Where(&db.Psychologist{Email: email}).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return nil, errors.New("username or password is incorrect")
	}

	expiresAt := time.Now().Add(time.Hour * 168).Unix() // valid for 7 days
	tokenString, err := utils.GenerateToken(user, expiresAt)
	if err != nil {
		return nil, err
	}

	var resp = map[string]interface{}{
		"Code":  http.StatusOK,
		"Token": tokenString,
		"User":  user,
	}

	return resp, nil
}

func UpdateProfileHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["Email"]
	var newProfile = &db.Psychologist{}
	var psy = &db.Psychologist{}
	var profile = &db.Profile{}

	utils.DecodeRequestBody(w, r, &newProfile)

	dbase.Find(&psy, &db.Psychologist{Email: email})
	dbase.Model(&psy).Updates(newProfile)

	dbase.Find(&psy, &db.Profile{Psychologist: psy.ID})
	dbase.Model(&profile).Updates(newProfile.Profile)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(psy))
}

func GetProfileHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["Email"]
	var profile = &db.Psychologist{}

	result := dbase.Find(&profile, &db.Psychologist{Email: email})
	if result.Error != nil {
		log.Println(result)
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Could not find your profile",
		}))
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(profile))
}

func UpdatePassword(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	password := struct {
		Password string `json:"password"`
	}{}

	json.NewDecoder(r.Body).Decode(&password)
	email := r.URL.Query().Get("email")
	var user *db.Psychologist

	dbase.Find(&user, &db.Psychologist{Email: email})

	pass, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user.Password = string(pass)

	dbase.Save(&user)
}

func LogoutHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	// remove the token from the client
	var conn *server.ConnectedClient
	result := dbase.Find(&conn, &server.ConnectedClient{UserID: claims.UserID})
	if result.Error != nil {
		log.Println(result.Error)
		json.NewEncoder(w).Encode("Your session has expired")
		return
	}

	conn.Client.Hub.Unregister <- conn
	json.NewEncoder(w).Encode("You have been logged out")
}

func ResetHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	userEmail := struct {
		Email string `json:"email"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&userEmail)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode("Could not parse your email"))
		return
	}

	var user *db.Psychologist
	result := dbase.Where(&db.Psychologist{Email: userEmail.Email}).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode("Email does not exist. Please create an account"))
		return
	}

	var writer bytes.Buffer
	body := struct {
		Name string
		Link string
	}{
		Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Link: fmt.Sprintf("https://google.com?email=%s", user.Email),
	}

	go func(dbase *gorm.DB, email string, HTMLtemp string, body interface{}, writer *bytes.Buffer) {
		err := utils.SendEmail(dbase, email, "templates/email/reset.html", body, writer)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}(dbase, userEmail.Email, "templates/email/reset.html", &body, &writer)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode("Please check your inbox for more action"))
}
