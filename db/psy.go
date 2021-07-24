package db

import "gorm.io/gorm"

type Psychologist struct {
	gorm.Model

	FirstName string `json:"FirstName"`
	LastName  string `json:"Lastname"`
	Email     string `gorm:"type:varchar(100);unique_index;Email"`
	Password  string `json:"Password"`
}
