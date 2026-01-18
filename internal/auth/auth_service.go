package auth

import (
	"strconv"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/model"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var accessExpiresIn = 15 * time.Minute      // 15 minutes
var refershExpiresIn = 7 * 24 * time.Minute // 7 days

type UserProvider interface {
	GetUserByUsername(email string) (*user.UserEntity, error)
}

type AuthService interface {
	Login(dto LoginRequest, userProvider UserProvider) (*TokenResponse, error)
	Verify(access string) (uint, error)
	Refresh(refresh string) (*TokenResponse, error)
}

type authService struct {
	db            *gorm.DB
	tokenRepo     TokenRepository
	accessSecret  []byte
	refreshSecret []byte
}

func NewAuthService(db *gorm.DB, tokenRepo TokenRepository, accessSecret, refreshSecret string) AuthService {
	return &authService{
		db:            db,
		tokenRepo:     tokenRepo,
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *authService) Login(dto LoginRequest, userProvider UserProvider) (*TokenResponse, error) {
	u, err := userProvider.GetUserByUsername(dto.Username)
	if err != nil {
		return nil, err
	}

	if err := util.VerifyPassword(u.Password, dto.Password); err != nil {
		return nil, err
	}

	var result *TokenResponse

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.revokeAllfreshTokensFromUser(tx, uint(u.ID)); err != nil {
			return err
		}

		userIDStr := strconv.FormatUint(uint64(u.ID), 10)
		tokens, err := s.issueTokens(tx, userIDStr)
		if err != nil {
			return err
		}

		result = tokens

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *authService) Verify(access string) (uint, error) {
	var accessClaims model.AccessTokenClaims
	accessToken, err := util.ParseJWTWithClaims(access, s.accessSecret, &accessClaims)
	if err != nil || !accessToken.Valid {
		return 0, err
	}

	userIDInt, err := strconv.Atoi(accessClaims.UserID)
	if err != nil {
		return 0, err
	}

	return uint(userIDInt), nil
}

func (s *authService) Refresh(refresh string) (*TokenResponse, error) {
	var refreshClaims jwt.RegisteredClaims
	refreshToken, err := util.ParseJWTWithClaims(refresh, s.refreshSecret, &refreshClaims)
	if err != nil || !refreshToken.Valid {
		return nil, apperror.ErrInvalidToken
	}

	tokenID := refreshClaims.ID

	t, err := s.tokenRepo.GetByTokenID(tokenID)
	if err != nil {
		return nil, err
	}

	var result *TokenResponse

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if t.IsRevoked {
			if err := s.revokeAllfreshTokensFromUser(tx, t.UserID); err != nil {
				return err
			}
		} else {
			if err := s.revokeRefreshToken(tx, t); err != nil {
				return err
			}
		}

		userIDStr := refreshClaims.Subject
		tokens, err := s.issueTokens(tx, userIDStr)
		if err != nil {
			return err
		}

		result = tokens

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *authService) revokeRefreshToken(tx *gorm.DB, t *TokenEntity) error {
	tokenRepository := s.tokenRepo.WithTx(tx)

	return tokenRepository.Revoke(t)
}

func (s *authService) revokeAllfreshTokensFromUser(tx *gorm.DB, userID uint) error {
	tokenRepository := s.tokenRepo.WithTx(tx)

	return tokenRepository.RevokeAllFromUser(userID)
}

func (s *authService) issueTokens(tx *gorm.DB, userIDStr string) (*TokenResponse, error) {
	tokenRepository := s.tokenRepo.WithTx(tx)

	accessExpiresAt := time.Now().Add(accessExpiresIn)
	access, err := util.GenerateAccessToken(userIDStr, accessExpiresAt, s.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenID := uuid.New().String()
	refreshExpiresAt := time.Now().Add(refershExpiresIn)
	refresh, err := util.GenerateRefreshToken(refreshTokenID, userIDStr, refreshExpiresAt, s.refreshSecret)
	if err != nil {
		return nil, err
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	t := &TokenEntity{TokenID: refreshTokenID, UserID: uint(userIDInt), ExpiresAt: refreshExpiresAt}
	if err := tokenRepository.Create(t); err != nil {
		return nil, err
	}

	return &TokenResponse{AccessToken: access, RefreshToken: refresh}, nil
}
