package main

import (
	"log"

	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/database"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: .env not found, using system env variables")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Migration failed: could not connect to databse: %v", err)
	}

	log.Println("-- Start Migration ---")

	models := []any{}
	models = append(models, user.GetModels()...)
	models = append(models, auth.GetModels()...)

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("-- Migration completed successfully ---")
}
