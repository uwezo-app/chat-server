package db

import "gorm.io/gorm"

type Psychologist struct {
	gorm.Model

	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `gorm:"type:varchar(100);unique_index;Email"`
	Password  string `json:"password"`
}

type Profile struct {
	gorm.Model

	Psychologist *Psychologist `json:"Psychologist"`
	Image        string        `json:"Image"`
	Description  string        `json:"Description"`
}
