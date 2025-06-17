package tests

import (
	"testing"
	"time"

	"cry-api/app/models"
	services "cry-api/app/services/users"
	mocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser_NewUser_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	fullname := "John Doe"
	username := "johndoe"
	email := "john@example.com"
	password := "password123"

	// No existing user found
	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(nil, nil)

	// Save should be called for new user
	mockUserRepo.On("Save", mock.AnythingOfType("*models.User")).Return(nil)

	// Execute
	user, token, err := userSvc.CreateUser(fullname, username, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_ExistingUnverifiedUser_RefreshToken(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	existingUser := &models.User{
		Username:                   "johndoe",
		Email:                      "john@example.com",
		IsVerified:                 false,
		VerificationTokens:         "oldtoken",
		VerificationTokenCreatedAt: time.Now().Add(-25 * time.Hour), // expired token
	}

	mockUserRepo.On("FindByUsernameOrEmail", existingUser.Username, existingUser.Email).Return(existingUser, nil)

	// Expect Save with refreshed token called
	mockUserRepo.On("Save", mock.MatchedBy(func(u *models.User) bool {
		return u.VerificationTokens != "oldtoken"
	})).Return(nil)

	user, token, err := userSvc.CreateUser("John Doe", existingUser.Username, existingUser.Email, "password123")

	assert.NoError(t, err)
	assert.Equal(t, existingUser.Username, user.Username)
	assert.NotEqual(t, "oldtoken", token)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_AuthenticateUser_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	// Generate bcrypt hash for password123
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:   "johndoe",
		Password:   string(hashedPassword),
		IsVerified: true,
	}

	mockUserRepo.On("FindByUsername", user.Username).Return(user, nil)

	authUser, err := userSvc.AuthenticateUser(user.Username, "password123")

	assert.NoError(t, err)
	assert.NotNil(t, authUser)
	assert.Equal(t, user.Username, authUser.Username)
}

func TestUserService_AuthenticateUser_UserNotVerified(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Username:   "johndoe",
		Password:   string(hashedPassword),
		IsVerified: false,
	}

	mockUserRepo.On("FindByUsername", user.Username).Return(user, nil)

	authUser, err := userSvc.AuthenticateUser(user.Username, "password123")

	assert.Error(t, err)
	assert.Nil(t, authUser)
	assert.Equal(t, "user not verified", err.Error())
}

func TestUserService_VerifyUserWithTokens_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	userToken := "userToken123"
	verifyToken := "verifyToken456"

	user := &models.User{
		Username:                 "johndoe",
		IsVerified:               false,
		VerificationTokens:       verifyToken,
		AccountVerificationToken: &userToken,
	}

	// Setup CheckUserByBothTokens to return user
	mockUserRepo.On("FindByVerificationToken", verifyToken).Return(user, nil)
	mockUserRepo.On("Save", mock.MatchedBy(func(u *models.User) bool {
		return u.IsVerified && u.VerificationTokens == "" && u.AccountVerificationToken == nil
	})).Return(nil)

	verifiedUser, err := userSvc.VerifyUserWithTokens(userToken, verifyToken)

	assert.NoError(t, err)
	assert.True(t, verifiedUser.IsVerified)
	assert.Empty(t, verifiedUser.VerificationTokens)
	assert.Nil(t, verifiedUser.AccountVerificationToken)
}

func TestUserService_VerifyUserWithTokens_InvalidTokens(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	userSvc := services.NewUserService(mockUserRepo, mockEmailSvc)

	mockUserRepo.On("FindByVerificationToken", "invalid").Return(nil, nil)

	user, err := userSvc.VerifyUserWithTokens("someUserToken", "invalid")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid verification token")
}
