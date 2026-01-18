package auth_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/auth"
	"github.com/Perajit/expense-tracker-go/internal/auth/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/Perajit/expense-tracker-go/internal/user"
	userMocks "github.com/Perajit/expense-tracker-go/internal/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var accessSecret = "access-secret"
var refreshSecret = "refresh-secret"

type Credentials struct {
	Email    string
	Password string
}

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 1
		dto := auth.LoginRequest{
			Username: "test",
			Password: "pwd123",
		}
		matchedUser := GenerateUser(1, user.CreateUserRequest{Username: dto.Username, Password: dto.Password, Email: "test@example.com"})
		var refreshTokenID string

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)
		mockTokenRepo.On("WithTx", mock.Anything).Return(mockTokenRepo)
		mockTokenRepo.On("RevokeAllFromUser", userID).Return(nil).Once()
		mockTokenRepo.On("Create", mock.MatchedBy(func(t *auth.TokenEntity) bool {
			if t.UserID != userID {
				return false
			}
			refreshTokenID = t.TokenID
			return true
		})).Return(nil)

		mockUserService := new(userMocks.MockUserService)
		mockUserService.On("GetUserByUsername", dto.Username).Return(matchedUser, nil).Once()

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Login(dto, mockUserService)

		assert.NoError(t, err)
		mockUserService.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)

		// verify access token
		accessClaims := ExtractAccessClaims(tokens.AccessToken)
		assert.Equal(t, strconv.Itoa(int(matchedUser.ID)), accessClaims.UserID)
		assert.Equal(t, strconv.Itoa(int(matchedUser.ID)), accessClaims.Subject)

		timeIn15Mins := time.Now().Add(15 * time.Minute)
		assert.Less(t, accessClaims.ExpiresAt.Time, timeIn15Mins)

		// verify refresh token
		refreshClaims := ExtractRefreshClaims(tokens.RefreshToken)
		assert.Equal(t, strconv.Itoa(int(matchedUser.ID)), refreshClaims.Subject)

		timeIn7Days := time.Now().Add(7 * 24 * time.Minute)
		assert.Less(t, refreshClaims.ExpiresAt.Time, timeIn7Days)
		assert.Greater(t, refreshClaims.ExpiresAt.Time, accessClaims.ExpiresAt.Time)

		// verify saved refresh token
		assert.Equal(t, refreshClaims.ID, refreshTokenID)
	})

	t.Run("error_notfound_email", func(t *testing.T) {
		dto := auth.LoginRequest{
			Username: "test",
			Password: "pwd123",
		}

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)

		mockUserService := new(userMocks.MockUserService)
		mockUserService.On("GetUserByUsername", dto.Username).Return(nil, apperror.ErrNotFound).Once()

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Login(dto, mockUserService)

		assert.Nil(t, tokens)
		assert.Equal(t, apperror.ErrNotFound, err)
		mockUserService.AssertExpectations(t)
		mockTokenRepo.AssertNotCalled(t, "RevokeAllFromUser")
		mockTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("error_incorrect_password", func(t *testing.T) {
		dto := auth.LoginRequest{
			Username: "test",
			Password: "pwd123",
		}
		matchedUser := GenerateUser(1, user.CreateUserRequest{Username: dto.Username, Password: "pwd456", Email: "test@example.com"})

		db := testutil.SetupDB()

		mockTokenRepo := new(mocks.MockTokenRepository)

		mockUserService := new(userMocks.MockUserService)
		mockUserService.On("GetUserByUsername", dto.Username).Return(matchedUser, nil).Once()

		service := auth.NewAuthService(db, mockTokenRepo, accessSecret, refreshSecret)
		tokens, err := service.Login(dto, mockUserService)

		assert.Nil(t, tokens)
		assert.NotNil(t, err)
		mockUserService.AssertExpectations(t)
		mockTokenRepo.AssertNotCalled(t, "RevokeAllFromUser")
		mockTokenRepo.AssertNotCalled(t, "Create")
	})
}
