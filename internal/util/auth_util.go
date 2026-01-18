package util

import (
	"time"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ParseJWTWithClaims(signed string, secret []byte, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return token, err
}

func GenerateAccessToken(userIDStr string, expiresAt time.Time, secret []byte) (string, error) {
	claims := model.AccessTokenClaims{
		UserID: userIDStr,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userIDStr,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, err
}

func GenerateRefreshToken(tokenID string, userIDStr string, expiresAt time.Time, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        tokenID,
		Subject:   userIDStr,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GetAuthUserID(c *fiber.Ctx) (uint, error) {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return 0, apperror.ErrSecurityContextMissing
	}

	return userID, nil
}

func SetAuthUserID(c *fiber.Ctx, userID uint) {
	c.Locals("user_id", userID)
}
