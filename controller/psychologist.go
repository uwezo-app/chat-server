package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

// https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8

type ErrorResponse struct {
	Code    int
	Message string
}

type Psychologist struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"Lastname"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
}

type Token struct {
	Name           string
	Email          string
	StandardClaims *jwt.StandardClaims
}

func (t Token) Valid() error {
	// Check if the token is expired
	// Check if the token has been revoked
	// by checking if the token matches the db entry
	panic("implement me")
}

// CreatePsychologist implements psychologist creation
func CreatePsychologist(w http.ResponseWriter, r *http.Request) {
	user := &Psychologist{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Could not hash your password",
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	user.Password = string(password)

	// Write to db

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err)
	}
}

// LoginHandler implements authentication for psychologists
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &Psychologist{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("An error occurred while processing your request: %s", err),
		}
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	resp, err := FindOne(user.Email, user.Password)
	if err != nil {
		log.Println(err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Username or password is incorrect",
		})
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
		return
	}
}

func FindOne(email, password string) (map[string]interface{}, error) {
	user := &Psychologist{}

	// Query to the db

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		errorResponse := ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Username or password is incorrect",
		}

		return nil, errors.New(errorResponse.Message)
	}

	tk := Token {
		Name:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var resp = map[string]interface{}{
		"Code":    http.StatusOK,
		"Message": "LoggedIn",
		"Token":   tokenString,
		"User":    user,
	}

	return resp, nil
}
