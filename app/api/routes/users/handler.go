package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cry-api/app/config"
	EmailCore "cry-api/app/core/email"
	JWT "cry-api/app/core/jwt"
	UserCore "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	EmailServices "cry-api/app/services/email"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler handles HTTP requests related to user operations
type Handler struct {
	UserService  UserServices.UserServiceInterface
	EmailService EmailServices.EmailServiceInterface
}

/*
NewHandler initializes and returns a new Handler instance with its dependencies.
It sets up the user repository, email sender, email service, and user service
using the provided GORM database connection and environment configuration.

Parameters:
  - db:   A pointer to the GORM database instance.
  - cfg:  A pointer to the environment configuration.

Returns:
  - A pointer to the initialized Handler.
*/
func NewHandler(db *gorm.DB, cfg *EnvTypes.EnvConfig) *Handler {
	userRepo := UserCore.NewGormUserRepository(db)
	emailSender := EmailCore.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)
	userService := UserServices.NewUserService(userRepo, emailService)

	return &Handler{UserService: userService, EmailService: emailService}
}

/*
Signup handles user registration requests.
It decodes the incoming JSON request into a UserDomain.User struct,
creates a new user using the userService, and sends a verification email asynchronously.
Responds with a success status if user creation is successful.
*/
func (h *Handler) Signup(c *gin.Context) {
	cfg := config.Get()

	var input UserTypes.IUserSignupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	createdUser, token, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(mapUserCreationErrorToStatusCode(err.Error()), gin.H{"error": err.Error()})
		return
	}

	go func(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) {
		verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, token)
		err := h.EmailService.SendVerifyAccountEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			verificationLink,
			user.VerificationTokens,
		)
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}(createdUser, token, cfg)

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

/*
VerifyEmailToken checks the validity of the email verification token. (This function is used to verify the email address of a user during the signup process.)
*/
func (h *Handler) VerifyEmailToken(c *gin.Context) {
	var req UserTypes.IVerificationTokenCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.UserService.VerifyUserWithTokens(req.UserToken, req.VerifyToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": user.IsVerified})
}

// SignIn method. auth + JWT
func (h *Handler) SignIn(c *gin.Context) {
	var req UserTypes.IUserSigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	jwt, err := JWT.GenerateJWT(user.UUID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
		"user": gin.H{
			"uuid":     user.UUID,
			"fullname": user.Fullname,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// VerifyAccountToken checks if the provided account verification token is valid and not expired.
// It expects a JSON body with a "token" field, retrieves the user associated with the token,
// and ensures the token matches and was created within the last 24 hours.
func (h *Handler) VerifyAccountToken(c *gin.Context) {
	var req UserTypes.IUserVerifyAccountTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Printf("VerifyAccountToken request: %+v\n", req)

	token := req.Token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := h.UserService.CheckAccountVerificationToken(token)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or expired"})
		return
	}

	timeLimit := time.Now().Add(-24 * time.Hour)
	if user.Token == nil || *user.Token != token || user.VerificationTokenCreatedAt.Before(timeLimit) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}
