package db

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID uint
	Name   string
	Email  string
	Role   string

	jwt.StandardClaims
}

//
//func (t CustomClaims) Valid() error {
//	// Check if the token is expired
//	// Check if the token has been revoked
//	// by checking if the token matches the db entry
//	panic("implement me")
//}
