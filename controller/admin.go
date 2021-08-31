package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/uwezo-app/chat-server/db"
	"github.com/uwezo-app/chat-server/utils"
)

const (
	ErrorResponse = "An error occured while processinf your request"
)

func AdminRegistrationHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := &db.Admin{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
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
	}(dbase, user.Email, "Welcome", "templates/email/admin.html", body)

	w.WriteHeader(http.StatusCreated)
	log.Println(json.NewEncoder(w).Encode(user))
}

func AdminLoginHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var admin *db.Admin
	type authInfo struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}
	var adminAuth *authInfo

	err := json.NewDecoder(r.Body).Decode(&adminAuth)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse)
		return
	}

	result := dbase.Where(db.Admin{Email: adminAuth.Email}).Find(&admin)
	if e := result.Error; e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("User not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(adminAuth.Password))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse)
		return
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	token, err := utils.GenerateAdminToken(admin, expiresAt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"Token": token,
		"User":  admin,
	})
}

func AdminLogoutHandler(dbase *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	var user *db.Admin
	result := dbase.Where(db.Admin{Email: claims.Email}).First(&user)
	if e := result.Error; e == nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}))
		return
	}

	// Invalidate user token
	tokenString, err = utils.GenerateAdminToken(user, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode(utils.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	log.Println(json.NewEncoder(w).Encode(tokenString))
}
