package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/users"
	"cry-api/app/middleware"
	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompleteProfile_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService: mockUserService,
		PasswordService: mockPasswordService,
	}

	existingUserInput := UserTypes.IUserSignupRequest{
		Fullname: "randomname",
		Username: "random420",
		Email: "johnny@yahoo.com",
		Password: "Randomizepassword@12",
	}

	dummyUser := &UserModel.User{
		ID: 1,
		Email: "johnny@yahoo.com",
		Username: "random420",
	}

	isVerified := true
	isProfileCompleted := false

	// Mock: CreateUser for mock updated user
	mockUserService.On("CreateUser", existingUserInput.Fullname, existingUserInput.Username, existingUserInput.Email, existingUserInput.Password, isVerified, isProfileCompleted).Return(dummyUser, nil)	
	
	// Mock: FindUserByEmail
	input := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe420",
		Email: "johnny@yahoo.com",
		Password: "Cannotbebreak@123",
	}
	
	mockUserService.On("FindUserByEmail", input.Email).Return(dummyUser, nil)

	// Mock: HashPassword

	// Fill new information to the user model

	// Mock: UpdateUser

	// Initialize HTTP Request

	// Assert result
	

	bodyBytes, _ := json.Marshal(input)
}