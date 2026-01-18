package middleware

import (
	"strings"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get access token from header
		authHeader := c.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// verify token and extract user id
		userID, err := authService.Verify(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrInvalidToken.Error()})
		}

		util.SetAuthUserID(c, userID)

		return c.Next()
	}
}
