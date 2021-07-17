package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8

type ErrorResponse struct {
	Code int
	Message string
}

type Psychologist struct {
	FirstName string `json:"FirstName"`
	LastName string `json:"Lastname"`
	Email      string `json:"Email"`
	Password string `json:"Password"`
}

type Token struct {
	UserID int
	Name string
	Email  string
	StandardClaims *jwt.StandardClaims
}

func (t Token) Valid() error {
	panic("implement me")
}

// CreatePsychologist implements psychologist creation
func CreatePsychologist(w http.ResponseWriter, r *http.Request) {
	psy := &Psychologist{}
	err := json.NewDecoder(r.Body).Decode(psy)
	if err != nil {
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(psy.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse{
			Code: http.StatusBadRequest,
			Message: "Could not decode your password",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	psy.Password = string(password)

	// Write to db

	_ = json.NewEncoder(w).Encode(psy)
}

// LoginHandler implements authentication for psychologists
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	psy := &Psychologist{}
	err := json.NewDecoder(r.Body).Decode(psy)
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("An error occurred while processing your request: %V", err),
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	resp, err := FindOne(psy.Email, psy.Password)
	if err != nil {
		log.Println(err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Code: http.StatusNotFound,
			Message: "Username or password is incorrect",
		})
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func FindOne(email, password string) (map[string]interface{}, error) {
	user := &Psychologist {}

	// Query to the db

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse {
			Code: http.StatusNotFound,
			Message: "Username or password is incorrect",
		}

		return nil, errors.New(errorResponse.Message)
	}

	tk := Token {
		UserID: 1, // user.ID
		Name:   user.FirstName + " " + user.LastName,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte("SECRET"))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var resp = map[string]interface{}{
		"Code": http.StatusOK,
		"Message": "LoggedIn",
		"Token": tokenString,
		"User": user,
	}

	return resp, nil
}
