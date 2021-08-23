package db

import "gorm.io/gorm"

type Admin struct {
	gorm.Model

	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
}
