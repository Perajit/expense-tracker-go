package database

import (
	"fmt"
	"log"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"gorm.io/gorm"
)

type CategorySeed struct {
	Name   string `json:"name"`
	UserID uint   `json:"userId"`
}

func seedCategories(db *gorm.DB, env string) {
	// seed from file
	fileName := fmt.Sprintf("%s_categories.json", env)
	var seeds []CategorySeed
	if err := loadSeedFile(fileName, &seeds); err != nil {
		log.Printf("Skip category: could not load: %s", fileName)
		return
	}

	for _, seed := range seeds {
		err := insertCategory(db, seed)
		if err != nil {
			log.Println(err)
		}
	}
}

func insertCategory(db *gorm.DB, data CategorySeed) error {
	var u expense.CategoryEntity

	// // check duplication
	err := db.Where("user_id = ?", data.UserID).Where("name = ?", data.Name).First(&u).Error
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("Skip category: category [%s] already exists for user [%d]", data.Name, data.UserID)
	}

	// save to database
	u = expense.CategoryEntity{UserID: data.UserID, Name: data.Name, IsDefault: data.UserID == 0}
	if err := db.Create(&u).Error; err != nil {
		return fmt.Errorf("Skip category: could not save category: %v", err)
	}

	log.Printf("Category [%s] for user [%d] created successfully", data.Name, data.UserID)
	return nil
}
