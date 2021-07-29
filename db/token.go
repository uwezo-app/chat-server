package db

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Token struct {
	*jwt.StandardClaims

	UserID		uint
	Name        string
	Email       string
}

type TokenString struct {
	gorm.Model

	ID			uint 	`gorm:"primarykey"`
	Token       string 	`json:"Token"`
}

func (t Token) Valid() error {
	// Check if the token is expired
	// Check if the token has been revoked
	// by checking if the token matches the db entry
	panic("implement me")
}
