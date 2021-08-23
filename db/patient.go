package db

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	NickName string `json:"NickName"`
	Password string `json:"Password"`

	PairedUsers []PairedUsers `gorm:"foreignKey:PatientID"`
}

func (p *Patient) TableName() string {
	return "patients"
}

func (p *Patient) GetPairedUsers() []PairedUsers {
	return p.PairedUsers
}
