package db

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type CustomClaims struct {
	UserID uint
	Name   string
	Email  string

	jwt.StandardClaims
}

type Token struct {
	gorm.Model

	UserID string `gorm:"primaryKey"`
	Token  string `json:"CustomClaims"`
}

//
//func (t CustomClaims) Valid() error {
//	// Check if the token is expired
//	// Check if the token has been revoked
//	// by checking if the token matches the db entry
//	panic("implement me")
//}
