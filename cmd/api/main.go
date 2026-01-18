//go:generate mockery
package main

import (
	"fmt"
	"os"

	"log"

	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/database"
	"github.com/Perajit/expense-tracker-go/internal/middleware"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: .env not found, using system env variables")
	}

	// init db connection
	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect to database")
	}

	// init app
	app := fiber.New()

	// set up dependencies
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	validate := validator.New()

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService, validate)

	tokenRepository := auth.NewTokenReposity(db)
	authService := auth.NewAuthService(db, tokenRepository, accessSecret, refreshSecret)
	// authService := auth.NewAuthService(accessSecret, refreshSecret)
	authHandler := auth.NewAuthHandler(authService, userService, validate)

	authMiddleware := middleware.AuthMiddleware(authService)

	// routes
	userHandler.RegisterRoutes(app, authMiddleware)
	authHandler.RegisterRoutes(app)

	// start app
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server is starting on pport %v", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
