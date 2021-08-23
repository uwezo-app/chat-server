package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred",
		})
		return
	}

	user.Password = utils.HashPassword(user.Password, w)
	if user.Password == "" {
		return
	}

	rs := dbase.Create(&user)
	if rs.Error != nil {
		log.Println(rs)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Could not create your account. Please try again later",
		}))
		return
	}

	body := struct {
		Name string
		Link string
	}{
		Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Link: "https://google.com",
	}

	go func(dbase *gorm.DB, email string, subject string, HTMLTemp string, body interface{}) {
		err := utils.SendEmail(dbase, email, subject, HTMLTemp, body)
		if err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	}(dbase, user.Email, "Welcome", "templates/email/confirm.html", body)

	w.WriteHeader(http.StatusCreated)
	log.Println(json.NewEncoder(w).Encode(user))
}

// LoginHandler implements authentication for psychologists
func LoginHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user *db.Psychologist
	var resp map[string]interface{}
	var err error

	err = json.NewDecoder(r.Body).Decode(&user)
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

	result := dbase.Where(db.Psychologist{Email: email}).First(&user)
	if result.Error != nil {
		log.Println(result)
		return nil, errors.New("user not found")
	}

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
	// email := mux.Vars(r)["Email"]
	var newProfile = &db.Psychologist{}

	utils.DecodeRequestBody(w, r, &newProfile)

	result := dbase.Session(&gorm.Session{FullSaveAssociations: true}).Save(&newProfile)
	if result.Error != nil {
		log.Println(result)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Could not update your profile. Please try again later",
		}))
		return
	}

	body := struct {
		Name  string
		Email string
	}{
		Name:  fmt.Sprintf("%s %s", newProfile.FirstName, newProfile.LastName),
		Email: "security@uwezo.app",
	}

	go func(dbase *gorm.DB, email string, subject string, HTMLTemp string, body interface{}) {
		err := utils.SendEmail(dbase, email, subject, HTMLTemp, body)
		if err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	}(dbase, newProfile.Email, "Profile update", "templates/email/profile.html", body)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(newProfile))
}

func GetProfileHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["Email"]
	var profile = &db.Psychologist{}

	result := dbase.First(&profile, db.Psychologist{Email: email})
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

func GetPsychologists(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	var psychologists []db.Psychologist
	if limit := queries.Get("limit"); limit != "" {
		l, _ := strconv.Atoi(limit)
		dbase.Limit(l).Find(&psychologists)
	} else if orderBy := queries.Get("order_by"); orderBy != "" {
		var limit string
		if limit = queries.Get("limit"); limit == "" {
			limit = "5"
		}
		l, _ := strconv.Atoi(limit)
		dbase.Order(orderBy + " desc").Limit(l).Find(&psychologists)
	} else {
		dbase.Find(&psychologists)
	}

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(psychologists))
}

func UpdatePassword(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	password := struct {
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Cod not read the input")
		return
	}
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

	var user *db.Psychologist
	dbase.Where(db.Psychologist{Email: claims.Email}).First(&user)
	if user == nil {
		log.Println("User not found")
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}))
		return
	}

	// Invalidate user token
	tokenString, err = utils.GenerateToken(user, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}))
		return
	}

	// remove the token from the client
	var conn *server.ConnectedClient
	result := dbase.Find(&conn, &server.ConnectedClient{UserID: claims.UserID})
	if result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("Your session has expired")
		return
	}

	conn.Client.Hub.Unregister <- conn

	w.WriteHeader(http.StatusNotFound)
	log.Println(json.NewEncoder(w).Encode(tokenString))
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

	body := struct {
		Name string
		Link string
	}{
		Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Link: fmt.Sprintf("https://google.com?email=%s", user.Email),
	}

	go func(dbase *gorm.DB, email string, subject string, HTMLTemp string, body interface{}) {
		err := utils.SendEmail(dbase, email, subject, HTMLTemp, body)
		if err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	}(dbase, userEmail.Email, "Reset Password", "templates/email/reset.html", body)

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode("Please check your inbox for more action"))
}
