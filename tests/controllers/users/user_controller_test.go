// Package tests provides test functionality for user controllers
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	UserController "cry-api/app/controllers/users"
	"cry-api/app/middleware"
	UserTypes "cry-api/app/types/users"
	"cry-api/tests/suites"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserControllerTestSuite provides tests for the user controller
type UserControllerTestSuite struct {
	suites.UserTestSuite
	router *gin.Engine
}

// SetupSuite initializes the test suite
func (suite *UserControllerTestSuite) SetupSuite() {
	suite.UserTestSuite.SetupSuite()

	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Add error handler middleware
	suite.router.Use(middleware.ErrorHandler())

	// Initialize controller with test container
	userController := UserController.NewUserController(suite.GetContainer())

	// Setup routes
	api := suite.router.Group("/api/v1")
	users := api.Group("/users")
	{
		users.POST("/signup", userController.Signup)
		users.POST("/signin", userController.SignIn)
		users.POST("/reset-password", userController.HandleResetPasswordRequest)
	}
}

// TestUserSignup tests the user signup functionality
func (suite *UserControllerTestSuite) TestUserSignup() {
	// Test data
	signupData := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(signupData)
	suite.Require().NoError(err)

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.True(suite.T(), response["success"].(bool))
}

// TestUserSignupInvalidData tests user signup with invalid data
func (suite *UserControllerTestSuite) TestUserSignupInvalidData() {
	// Test data with invalid email
	signupData := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "invalid-email",
		Password: "SecurePass123!",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(signupData)
	suite.Require().NoError(err)

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), response, "error")
}

// TestUserSignupDuplicateUser tests user signup with duplicate user
func (suite *UserControllerTestSuite) TestUserSignupDuplicateUser() {
	// Test data
	signupData := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
	}

	// First signup
	jsonData, err := json.Marshal(signupData)
	suite.Require().NoError(err)

	req1, err := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	suite.router.ServeHTTP(w1, req1)
	assert.Equal(suite.T(), http.StatusCreated, w1.Code)

	// Second signup with same data (should fail)
	req2, err := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)
	assert.Equal(suite.T(), http.StatusConflict, w2.Code)
}

// TestUserSignupWeakPassword tests user signup with weak password
func (suite *UserControllerTestSuite) TestUserSignupWeakPassword() {
	// Test data with weak password
	signupData := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "123", // Weak password
	}

	// Convert to JSON
	jsonData, err := json.Marshal(signupData)
	suite.Require().NoError(err)

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/users/signup", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), response, "error")
}

// RunUserControllerTestSuite runs the user controller test suite
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
