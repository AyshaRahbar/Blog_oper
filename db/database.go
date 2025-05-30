package db

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	var db *gorm.DB

	dbType := os.Getenv("DB_TYPE")
	if dbType == "postgres" {
		db, err = initPostgreSQL()
	} 
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initPostgreSQL() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSLMode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

