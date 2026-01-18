package auth_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/auth/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefresh(t *testing.T) {
	t.Run("success_unrevoked_token", func(t *testing.T) {
		refreshToken := &auth.TokenEntity{
			TokenID:   "123",
			UserID:    1,
			IsRevoked: false,
		}
		refresh := GenerateRefreshToken(refreshToken.TokenID, refreshToken.UserID, time.Now().Add(time.Hour))

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)
		mockTokenRepo.On("WithTx", mock.Anything).Return(mockTokenRepo)
		mockTokenRepo.On("GetByTokenID", refreshToken.TokenID).Return(refreshToken, nil).Once()
		mockTokenRepo.On("Revoke", mock.MatchedBy(func(t *auth.TokenEntity) bool {
			if t != refreshToken {
				return false
			}
			refreshToken.IsRevoked = false
			return true
		})).Return(nil).Once()
		mockTokenRepo.On("Create", mock.Anything).Return(nil)

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Refresh(refresh)

		assert.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)

		// verify access token
		accessClaims := ExtractAccessClaims(tokens.AccessToken)
		assert.Equal(t, strconv.Itoa(1), accessClaims.UserID)
		assert.Equal(t, strconv.Itoa(1), accessClaims.Subject)

		timeIn15Mins := time.Now().Add(15 * time.Minute)
		assert.Less(t, accessClaims.ExpiresAt.Time, timeIn15Mins)

		// verify refresh token
		refreshClaims := ExtractRefreshClaims(tokens.RefreshToken)
		assert.Equal(t, strconv.Itoa(1), refreshClaims.Subject)

		timeIn7Days := time.Now().Add(7 * 24 * time.Minute)
		assert.Less(t, refreshClaims.ExpiresAt.Time, timeIn7Days)
		assert.Greater(t, refreshClaims.ExpiresAt.Time, accessClaims.ExpiresAt.Time)
	})

	t.Run("success_revoked_token", func(t *testing.T) {
		refreshToken := &auth.TokenEntity{
			TokenID:   "123",
			UserID:    1,
			IsRevoked: true,
		}
		refresh := GenerateRefreshToken(refreshToken.TokenID, refreshToken.UserID, time.Now().Add(time.Hour))

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)
		mockTokenRepo.On("WithTx", mock.Anything).Return(mockTokenRepo)
		mockTokenRepo.On("GetByTokenID", refreshToken.TokenID).Return(refreshToken, nil).Once()
		mockTokenRepo.On("RevokeAllFromUser", mock.MatchedBy(func(userID uint) bool {
			if userID != refreshToken.UserID {
				return false
			}
			refreshToken.IsRevoked = false
			return true
		})).Return(nil).Once()
		mockTokenRepo.On("Create", mock.Anything).Return(nil)

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Refresh(refresh)

		assert.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)

		// verify access token
		accessClaims := ExtractAccessClaims(tokens.AccessToken)
		assert.Equal(t, strconv.Itoa(1), accessClaims.UserID)
		assert.Equal(t, strconv.Itoa(1), accessClaims.Subject)

		timeIn15Mins := time.Now().Add(15 * time.Minute)
		assert.Less(t, accessClaims.ExpiresAt.Time, timeIn15Mins)

		// verify refresh token
		refreshClaims := ExtractRefreshClaims(tokens.RefreshToken)
		assert.Equal(t, strconv.Itoa(1), refreshClaims.Subject)

		timeIn7Days := time.Now().Add(7 * 24 * time.Minute)
		assert.Less(t, refreshClaims.ExpiresAt.Time, timeIn7Days)
		assert.Greater(t, refreshClaims.ExpiresAt.Time, accessClaims.ExpiresAt.Time)
	})

	t.Run("error_invalid_token", func(t *testing.T) {
		mockToken := new(mocks.MockTokenRepository)

		db := testutil.SetupDB()

		service := auth.NewAuthService(db, mockToken, accessSecret, refreshSecret)
		tokens, err := service.Refresh("invalid")

		assert.Nil(t, tokens)
		assert.Equal(t, apperror.ErrInvalidToken, err)
		mockToken.AssertNotCalled(t, "Create")
	})

	t.Run("error_expired_token", func(t *testing.T) {
		refresh := GenerateRefreshToken("123", 1, time.Now().Add(-time.Hour))

		mockToken := new(mocks.MockTokenRepository)

		db := testutil.SetupDB()

		service := auth.NewAuthService(db, mockToken, accessSecret, refreshSecret)
		tokens, err := service.Refresh(refresh)

		assert.Nil(t, tokens)
		assert.Equal(t, apperror.ErrInvalidToken, err)
		mockToken.AssertNotCalled(t, "Create")
	})

	t.Run("error_revoke_token", func(t *testing.T) {
		refreshToken := &auth.TokenEntity{
			TokenID:   "123",
			UserID:    1,
			IsRevoked: false,
		}
		refresh := GenerateRefreshToken(refreshToken.TokenID, refreshToken.UserID, time.Now().Add(time.Hour))

		mockTokenRepo := new(mocks.MockTokenRepository)
		mockTokenRepo.On("WithTx", mock.Anything).Return(mockTokenRepo)
		mockTokenRepo.On("GetByTokenID", refreshToken.TokenID).Return(refreshToken, nil).Once()
		mockTokenRepo.On("Revoke", refreshToken).Return(apperror.ErrDefault).Once()

		db := testutil.SetupDB()

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Refresh(refresh)

		assert.Nil(t, tokens)
		assert.Equal(t, apperror.ErrDefault, err)
		mockTokenRepo.AssertExpectations(t)
		mockTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("error_save_token", func(t *testing.T) {
		refreshToken := &auth.TokenEntity{
			TokenID:   "123",
			UserID:    1,
			IsRevoked: false,
		}
		refresh := GenerateRefreshToken(refreshToken.TokenID, refreshToken.UserID, time.Now().Add(time.Hour))

		mockTokenRepo := new(mocks.MockTokenRepository)
		mockTokenRepo.On("WithTx", mock.Anything).Return(mockTokenRepo)
		mockTokenRepo.On("GetByTokenID", refreshToken.TokenID).Return(refreshToken, nil).Once()
		mockTokenRepo.On("Revoke", mock.MatchedBy(func(t *auth.TokenEntity) bool {
			if t != refreshToken {
				return false
			}
			refreshToken.IsRevoked = false
			return true
		})).Return(nil).Once()
		mockTokenRepo.On("Create", mock.Anything).Return(apperror.ErrDefault)

		db := testutil.SetupDB()

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Refresh(refresh)

		assert.Nil(t, tokens)
		assert.Equal(t, apperror.ErrDefault, err)
		mockTokenRepo.AssertExpectations(t)
	})
}
