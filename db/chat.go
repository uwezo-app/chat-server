package db

import (
	"time"

	"github.com/uwezo-app/chat-server/server"
	"gorm.io/gorm"
)

// ConnectedClient holds the connection info of a specific client
type ConnectedClient struct {
	gorm.Model

	UserID uint `gorm:"primaryKey"`

	Client *server.Client `gorm:"embedded"`

	LastSeen time.Time
}

type PairedUsers struct {
	gorm.Model

	PsychologistID uint

	PatientID uint

	Conversation []Conversation `gorm:"foreignKey:ConversationID"`

	EncryptionID string

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
