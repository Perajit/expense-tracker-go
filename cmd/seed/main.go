package main

import (
	"flag"
	"log"

	"github.com/Perajit/expense-tracker-go/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	env := *flag.String("env", "dev", "environment")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: .env not found, using system env variables")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Seeding Failed: could not connect to database: %v", err)
	}

	log.Printf("-- Start Seeding for env %s ---", env)

	database.Seed(db, env)

	log.Println("-- Seeding completed successfully ---")
}
