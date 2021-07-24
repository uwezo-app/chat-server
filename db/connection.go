package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"log"
	"os"

	_ "gorm.io/driver/postgres" //Gorm postgres dialect interface
	"gorm.io/gorm"
)

//ConnectDB : Make database connection
func ConnectDB() *gorm.DB {
	var db *gorm.DB
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbHost := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")

	dns := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, user, dbName, port, pass)

	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database %v\n", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Psychologist{}, &TokenString{})
	if err != nil {
		log.Fatalf("Could not run migrations: %v\n", err)
	}

	log.Println("Successfully connected!", db)
	return db
}
