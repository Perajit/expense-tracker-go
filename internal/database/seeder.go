package database

import (
	"encoding/json"

	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB, env string) {
	log.Println("Starting database seeding")

	seedUsers(db, env)
	seedCategories(db, env)

	log.Println("Completed database seeding")
}

func loadSeedFile(fileName string, target interface{}) error {
	path := filepath.Join("internal", "database", "seeds", fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}
