module github.com/uwezo-app/chat-server

// +heroku goVersion go1.16
go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/joho/godotenv v1.3.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.12
)
