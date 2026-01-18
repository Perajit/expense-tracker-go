package auth_test

import (
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/auth/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		access := GenerateAccessToken(1, time.Now().Add(time.Hour))

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		userID, err := service.Verify(access)

		assert.Equal(t, uint(1), userID)
		assert.NoError(t, err)
	})

	t.Run("error_invalid_token", func(t *testing.T) {
		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		userID, err := service.Verify("invalid")

		assert.Equal(t, uint(0), userID)
		assert.Error(t, apperror.ErrInvalidToken, err)
	})

	t.Run("error_expired_token", func(t *testing.T) {
		access := GenerateAccessToken(1, time.Now().Add(-time.Hour))

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		userID, err := service.Verify(access)

		assert.Equal(t, uint(0), userID)
		assert.Error(t, apperror.ErrInvalidToken, err)
	})
}
