package db

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	NickName string `json:"NickName"`

	PairedUsers []PairedUsers `gorm:"foreignKey:PatientID"`
}
