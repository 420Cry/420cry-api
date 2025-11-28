package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/users"
	"cry-api/app/middleware"
	UserModel "cry-api/app/models"
	services "cry-api/app/services/jwt"
	app_errors "cry-api/app/types/errors"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestRouter creates a test router with error handling middleware
func setupTestRouter(userController *controller.UserController, claims *services.Claims) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	// Add user claims middleware if provided
	if claims != nil {
		router.Use(func(c *gin.Context) {
			c.Set("user", claims)
			c.Next()
		})
	}

	router.PUT("/update-account-name", userController.UpdateAccountName)
	return router
}

// setupTestRouterWithString creates a test router with invalid claims (string)
func setupTestRouterWithString(userController *controller.UserController, invalidClaims string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("user", invalidClaims)
		c.Next()
	})
	router.PUT("/update-account-name", userController.UpdateAccountName)
	return router
}

func TestUpdateAccountName_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	existingUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: "oldusername",
	}

	updatedUser := &UserModel.User{
		UUID:     existingUser.UUID,
		Fullname: existingUser.Fullname,
		Email:    existingUser.Email,
		Username: input.AccountName,
	}

	// Mock JWT claims in context
	claims := &services.Claims{
		UUID:  existingUser.UUID,
		Email: existingUser.Email,
	}

	// Mock: Get user by UUID
	mockUserService.
		On("GetUserByUUID", existingUser.UUID).
		Return(existingUser, nil)

	// Mock: Check if username already exists (should return nil = not found)
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return((*UserModel.User)(nil), nil)

	// Mock: Update user
	mockUserService.
		On("UpdateUser", mock.MatchedBy(func(u *UserModel.User) bool {
			return u.UUID == existingUser.UUID && u.Username == input.AccountName
		})).
		Return(nil).
		Run(func(args mock.Arguments) {
			user := args.Get(0).(*UserModel.User)
			user.Username = updatedUser.Username
		})

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	c.Set("user", claims)

	userController.UpdateAccountName(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"].(bool))
	assert.Equal(t, "Username updated successfully", respBody["message"])

	mockUserService.AssertExpectations(t)
}

func TestUpdateAccountName_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`)
	router := setupTestRouter(userController, nil)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid input")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
	mockUserService.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateAccountName_MissingUserClaims(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	router := setupTestRouter(userController, nil) // No claims

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "User not authenticated")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
}

func TestUpdateAccountName_InvalidUserClaims(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	router := setupTestRouterWithString(userController, "invalid-claims-string")

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	// When claims are set but have the wrong type, the type assertion fails
	// and the controller returns "Invalid user claims"
	assert.Contains(t, respBody["error"], "Invalid user claims")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
}

func TestUpdateAccountName_UserNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	claims := &services.Claims{
		UUID:  "non-existent-uuid",
		Email: "test@example.com",
	}

	// Mock: User not found
	mockUserService.
		On("GetUserByUUID", claims.UUID).
		Return((*UserModel.User)(nil), app_errors.NewNotFoundError("user", "User not found"))

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "User not found")

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateAccountName_UsernameAlreadyExists(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "existingusername",
	}
	bodyBytes, _ := json.Marshal(input)

	currentUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: "oldusername",
	}

	existingUser := &UserModel.User{
		UUID:     "different-uuid",
		Fullname: "Jane Smith",
		Email:    "jane@example.com",
		Username: input.AccountName,
	}

	claims := &services.Claims{
		UUID:  currentUser.UUID,
		Email: currentUser.Email,
	}

	// Mock: Get current user by UUID
	mockUserService.
		On("GetUserByUUID", currentUser.UUID).
		Return(currentUser, nil)

	// Mock: Username already exists (different user)
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return(existingUser, nil)

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusConflict, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Username is already in use")

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateAccountName_SameUsername(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	currentUsername := "johndoe"
	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: currentUsername, // Same as current username
	}
	bodyBytes, _ := json.Marshal(input)

	currentUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: currentUsername,
	}

	claims := &services.Claims{
		UUID:  currentUser.UUID,
		Email: currentUser.Email,
	}

	// Mock: Get current user by UUID
	mockUserService.
		On("GetUserByUUID", currentUser.UUID).
		Return(currentUser, nil)

	// Mock: Username exists but it's the same user
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return(currentUser, nil)

	// Mock: Update user (should still be called)
	mockUserService.
		On("UpdateUser", mock.MatchedBy(func(u *UserModel.User) bool {
			return u.UUID == currentUser.UUID && u.Username == input.AccountName
		})).
		Return(nil)

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"].(bool))

	mockUserService.AssertExpectations(t)
}

func TestUpdateAccountName_CheckUsernameAvailabilityError(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	currentUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: "oldusername",
	}

	claims := &services.Claims{
		UUID:  currentUser.UUID,
		Email: currentUser.Email,
	}

	// Mock: Get current user by UUID
	mockUserService.
		On("GetUserByUUID", currentUser.UUID).
		Return(currentUser, nil)

	// Mock: Database error when checking username availability
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return((*UserModel.User)(nil), app_errors.NewInternalServerError("Failed to check username availability"))

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Failed to check username availability")

	mockUserService.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateAccountName_UpdateUserError(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "newusername",
	}
	bodyBytes, _ := json.Marshal(input)

	currentUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: "oldusername",
	}

	claims := &services.Claims{
		UUID:  currentUser.UUID,
		Email: currentUser.Email,
	}

	// Mock: Get current user by UUID
	mockUserService.
		On("GetUserByUUID", currentUser.UUID).
		Return(currentUser, nil)

	// Mock: Username is available
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return((*UserModel.User)(nil), nil)

	// Mock: Update user fails
	mockUserService.
		On("UpdateUser", mock.AnythingOfType("*models.User")).
		Return(app_errors.NewInternalServerError("Failed to update username"))

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Failed to update username")

	mockUserService.AssertExpectations(t)
}

func TestUpdateAccountName_EmptyUsername(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "", // Empty username
	}
	bodyBytes, _ := json.Marshal(input)

	router := setupTestRouter(userController, nil)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid input")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
}

func TestUpdateAccountName_UsernameTooShort(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "ab", // Only 2 characters (minimum is 3)
	}
	bodyBytes, _ := json.Marshal(input)

	claims := &services.Claims{
		UUID:  "test-uuid-1234",
		Email: "test@example.com",
	}

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Username must be at least 3 characters long")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
}

func TestUpdateAccountName_UsernameTooLong(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	// Create a username with 51 characters (maximum is 50)
	longUsername := "a123456789012345678901234567890123456789012345678901"
	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: longUsername,
	}
	bodyBytes, _ := json.Marshal(input)

	claims := &services.Claims{
		UUID:  "test-uuid-1234",
		Email: "test@example.com",
	}

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Username must not exceed 50 characters")

	mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
}

func TestUpdateAccountName_UsernameInvalidCharacters(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	// Test with various invalid characters
	testCases := []struct {
		name     string
		username string
	}{
		{"with spaces", "user name"},
		{"with hyphens", "user-name"},
		{"with dots", "user.name"},
		{"with special chars", "user@name"},
		{"with symbols", "user$name"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := UserTypes.IUserUpdateAccountNameRequest{
				AccountName: tc.username,
			}
			bodyBytes, _ := json.Marshal(input)

			claims := &services.Claims{
				UUID:  "test-uuid-1234",
				Email: "test@example.com",
			}

			router := setupTestRouter(userController, claims)

			req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			var respBody map[string]string
			err := json.NewDecoder(res.Body).Decode(&respBody)
			assert.NoError(t, err)
			assert.Contains(t, respBody["error"], "Username can only contain letters, numbers, and underscores")

			mockUserService.AssertNotCalled(t, "GetUserByUUID", mock.Anything)
		})
	}
}

func TestUpdateAccountName_ValidUsernameWithUnderscores(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserUpdateAccountNameRequest{
		AccountName: "user_name_123", // Valid username with underscores and numbers
	}
	bodyBytes, _ := json.Marshal(input)

	existingUser := &UserModel.User{
		UUID:     "test-uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: "oldusername",
	}

	claims := &services.Claims{
		UUID:  existingUser.UUID,
		Email: existingUser.Email,
	}

	// Mock: Get user by UUID
	mockUserService.
		On("GetUserByUUID", existingUser.UUID).
		Return(existingUser, nil)

	// Mock: Check if username already exists (should return nil = not found)
	mockUserService.
		On("FindUserByUsername", input.AccountName).
		Return((*UserModel.User)(nil), nil)

	// Mock: Update user
	mockUserService.
		On("UpdateUser", mock.MatchedBy(func(u *UserModel.User) bool {
			return u.UUID == existingUser.UUID && u.Username == input.AccountName
		})).
		Return(nil)

	router := setupTestRouter(userController, claims)

	req := httptest.NewRequest(http.MethodPut, "/update-account-name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"].(bool))
	assert.Equal(t, "Username updated successfully", respBody["message"])

	mockUserService.AssertExpectations(t)
}
