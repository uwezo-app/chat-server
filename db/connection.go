package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres" //Gorm postgres dialect interface
	"gorm.io/gorm"
)

//ConnectDB : Make database connection
func ConnectDB() *gorm.DB {
	var db *gorm.DB
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	dns := getDNS()

	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database %v\n", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&Psychologist{},
		&Profile{},
		&Patient{},
		&PairedUsers{},
		&Conversation{},
	)

	if err != nil {
		log.Fatalf("Could not run migrations: %v\n", err)
	}

	log.Printf("Successfully connected! %v\n", dns)
	return db
}

func getDNS() string {

	if os.Getenv("APP_ENV") == "development" {
		user := os.Getenv("DATABASE_USER")
		pass := os.Getenv("DATABASE_PASSWORD")
		dbName := os.Getenv("DATABASE_NAME")
		dbHost := os.Getenv("DATABASE_HOST")
		port := os.Getenv("DATABASE_PORT")
		sslmode := os.Getenv("DATABASE_SSLMODE")

		return fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s password=%s", dbHost, user, dbName, port, sslmode, pass)
	}

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	user := os.Getenv("GOOGLE_DB_USER")
	pass := os.Getenv("GOOGLE_DB_PASSOWRD")
	dbName := os.Getenv("GOOGLE_DB_NAME")
	//dbHost := os.Getenv("DATABASE_HOST")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s/%s", user, pass, dbName, socketDir, instanceConnectionName)
}
