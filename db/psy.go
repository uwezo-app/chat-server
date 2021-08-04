package db

import "gorm.io/gorm"

type Psychologist struct {
	gorm.Model

	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `gorm:"primaryKey;autoIncrement:false"`
	Password  string `json:"password"`

	Profile Profile `gorm:"foreignKey:Psychologist"`
}

type Profile struct {
	gorm.Model

	Psychologist string `gorm:"primaryKey"`
	Image        string `json:"Image"`
	Description  string `json:"Description"`
}
