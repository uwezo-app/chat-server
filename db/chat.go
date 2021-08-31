package db

import (
	"time"

	"gorm.io/gorm"
)

type PairedUsers struct {
	gorm.Model

	PsychologistID uint

	PatientID uint

	Conversation []Conversation `gorm:"foreignKey:ConversationID"`

	EncryptionKey string

	PairedAt time.Time
}

type Conversation struct {
	gorm.Model

	ConversationID uint

	From uint

	Message string

	Url string

	SentAt time.Time
}
