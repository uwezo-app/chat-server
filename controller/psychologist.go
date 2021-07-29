package controller

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/uwezo-app/chat-server/db"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

// https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8

type ErrorResponse struct {
	Code    int
	Message string
}

var dbase = db.ConnectDB()

// CreatePsychologist implements psychologist creation
func CreatePsychologist(w http.ResponseWriter, r *http.Request) {
	user := &db.Psychologist{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return
	}

	var password []byte
	password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := ErrorResponse{
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
		log.Println(json.NewEncoder(w).Encode(ErrorResponse{
			Code:    	http.StatusInternalServerError,
			Message:	"Could not create your account. Please try again later",
		}))
	}
	log.Printf("Row affected %v\n", rs.RowsAffected)

	log.Println(json.NewEncoder(w).Encode(user))
}

// LoginHandler implements authentication for psychologists
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user = &db.Psychologist{}
	var resp = make(map[string]interface{})
	var err error

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("An error occurred while processing your request: %s", err),
		}
		log.Println(json.NewEncoder(w).Encode(errorResponse))
		return
	}

	resp, err = FindOne(user.Email, user.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}))
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode(resp))
}

func FindOne(email, password string) (map[string]interface{}, error) {
	user := &db.Psychologist{}

	dbase.Where(&db.Psychologist{Email: email}).First(user)

	expiresAt := time.Now().Add(time.Minute * 10080).Unix() // valid for 7 days

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return nil, errors.New("username or password is incorrect")
	}

	tk := db.Token {
		UserID: 			user.ID,
		Name:   			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:  			user.Email,
		StandardClaims: 	&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	var tokenString string
	tokenString, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dbase.Create(&db.TokenString{Token: tokenString, ID: user.ID})

	var resp = map[string]interface{} {
		"Code":    http.StatusOK,
		"Token":   tokenString,
		"User":    user,
	}

	return resp, nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request){}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println(json.NewEncoder(w).Encode("Could not parse your email").Error())
		return
	}

	var user *db.Psychologist
	result := dbase.Where(&db.Psychologist{Email: userEmail.Email}).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(json.NewEncoder(w).Encode("Email does not exist. Please create an account"))
		return
	}

	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")

	to := []string {
		user.Email,
	}

	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From": {m.FormatAddress(from, "Uwezo Team")},
		"To": to,
		"Subject": {"Reset Password"},
	})

	t, _ := template.ParseFiles("templates/email/reset.html")
	var body bytes.Buffer

	err = t.Execute(&body, struct {
		Name    string
		Link string
	}{
		Name:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Link: "https://google.com",
	})
	if err != nil {
		log.Println(err)
		return
	}

	m.SetBody("text/html", string(body.Bytes()))

	p, _ := strconv.Atoi(port)
	d := gomail.NewDialer(host, p, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(json.NewEncoder(w).Encode("Could not send you a confirmation email"))
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println(json.NewEncoder(w).Encode("Please check your inbox for more action"))
}
