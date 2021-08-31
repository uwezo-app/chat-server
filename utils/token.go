package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/uwezo-app/chat-server/db"
)

var (
	ErrTokenInvalid   = errors.New("token invalid")
	ErrTokenIsExpired = errors.New("token is expired")
)

func GenerateToken(user *db.Psychologist, expiresAt int64) (token string, err error) {
	claims := db.CustomClaims{
		UserID: user.ID,
		Name:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	var tokenString string
	tokenString, err = t.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func GeneratePatientToken(user *db.Patient, expiresAt int64) (token string, err error) {
	claims := db.CustomClaims{
		UserID: user.ID,
		Name:   user.NickName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := t.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func ParseTokenWithClaims(tokenString string) (*db.CustomClaims, error) {
	tk, err := jwt.ParseWithClaims(tokenString, &db.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := tk.Claims.(*db.CustomClaims)
	if !ok || !tk.Valid {
		log.Println(err)
		return nil, ErrTokenIsExpired
	}

	return claims, nil
}

func GetTokenFromHeader(h http.Header) string {

	token := h.Get("Authorization")
	return strings.Split(token, " ")[1]
}
