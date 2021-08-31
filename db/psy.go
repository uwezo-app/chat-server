package db

import "gorm.io/gorm"

type Psychologist struct {
	gorm.Model

	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email" gorm:"primaryKey;autoIncrement:false"`
	Password  string `json:"Password"`
	IsVerfied bool   `json:"IsVerfied" gorm:"default:false"`
	IsDeleted bool   `json:"IsDeleted" gorm:"default:false"`

	PairedUsers []PairedUsers `gorm:"foreignKey:PsychologistID"`
	Profile     Profile       `gorm:"foreignKey:Psychologist"`
}

type Profile struct {
	gorm.Model

	ID           uint
	Psychologist uint   `gorm:"primaryKey;autoIncrement:false"`
	PhoneNumber  string `json:"PhoneNumber"`
	Image        string `json:"Image"`
	Country      string `json:"Country"`
	FocusedArea  string `json:"FocusedArea"`
	Address      string `json:"Address"`
	Description  string `json:"Description"`
}
