package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}
