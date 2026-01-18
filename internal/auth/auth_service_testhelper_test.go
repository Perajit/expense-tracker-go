package auth_test

import (
	"strconv"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/model"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var AccessSecret = "access-secret"
var RefreshSecret = "refresh-secret"

func ExtractAccessClaims(access string) *model.AccessTokenClaims {
	accessClaims := &model.AccessTokenClaims{}
	util.ParseJWTWithClaims(access, []byte(AccessSecret), accessClaims)

	return accessClaims
}

func ExtractRefreshClaims(refresh string) *jwt.RegisteredClaims {
	refreshClaims := &jwt.RegisteredClaims{}
	util.ParseJWTWithClaims(refresh, []byte(RefreshSecret), refreshClaims)

	return refreshClaims
}

func GenerateAccessToken(userID uint, expiresAt time.Time) string {
	signed, _ := util.GenerateAccessToken(strconv.Itoa(int(userID)), expiresAt, []byte(AccessSecret))

	return signed
}

func GenerateRefreshToken(tokenID string, userID uint, expiresAt time.Time) string {
	signed, _ := util.GenerateRefreshToken(tokenID, strconv.Itoa(int(userID)), expiresAt, []byte(RefreshSecret))

	return signed
}

func GenerateUser(id uint, dto user.CreateUserRequest) *user.UserEntity {
	hashedPassword, _ := util.HashPassword(dto.Password)

	return &user.UserEntity{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
		},
		Username: dto.Username,
		Password: string(hashedPassword),
		Email:    dto.Email,
	}
}
