package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Perajit/expense-tracker-go/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeed struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func seedUsers(db *gorm.DB, env string) {
	// seed from env
	envJSON := os.Getenv("SEED_USER_JSON")
	if envJSON != "" {
		seed := &UserSeed{}
		if err := json.Unmarshal([]byte(envJSON), seed); err != nil {
			log.Printf("Skip user: could not read seed user from environment: %v", err)
			return
		}
		if err := insertUser(db, *seed); err != nil {
			log.Printf("Skip user: cound not insert seed user from environment: %v", err)
			return
		}
		return
	}

	// seed from file
	fileName := fmt.Sprintf("%s_users.json", env)
	var seeds []UserSeed
	if err := loadSeedFile(fileName, &seeds); err != nil {
		log.Printf("Skip user: could not load: %s", fileName)
		return
	}

	for _, seed := range seeds {
		err := insertUser(db, seed)
		if err != nil {
			log.Println(err)
		}
	}
}

func insertUser(db *gorm.DB, data UserSeed) error {
	var u user.UserEntity

	// check duplication
	err := db.Where("username = ?", data.Username).First(&u).Error
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("Skip user: user [%s] already exists", data.Username)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Skip user: could not hash password: %v", err)
	}

	// save to database
	u = user.UserEntity{Email: data.Email, Password: string(hashedPassword), Username: data.Username}
	if err := db.Create(&u).Error; err != nil {
		return fmt.Errorf("Skip user: could not save user: %v", err)
	}

	log.Printf("User [%s] created successfully", data.Email)
	return nil
}
