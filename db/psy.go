package db

import "gorm.io/gorm"

type Psychologist struct {
	gorm.Model

	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `gorm:"type:varchar(100);unique_index;email"`
	Password  string `json:"password"`
}

type Profile struct {
	gorm.Model

	Psychologist *Psychologist `json:"psychologist"`
	Image        string        `json:"image"`
	Description  string        `json:"description"`
}
