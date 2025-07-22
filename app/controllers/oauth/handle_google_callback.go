package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	cfg "cry-api/app/config"
	"cry-api/app/factories"
	jwtService "cry-api/app/services/jwt"
)

func (h *OAuthController) HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	cryAppUrl := cfg.Get().CryAppURL
	appEnv := cfg.Get().AppEnv

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code invalid"})
		log.Println("No code found")
		return
	}

	token, err := h.OAuthService.ExchangeToken(context.Background(), code)

	if err != nil {
		log.Println("error exchanging: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed oauth exchange"})
		return
	}

	googleUserInfo, err := h.OAuthService.FetchUserInfo(token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info from Google"})
		return
	}

	existingUser, err := h.UserService.FindUserByEmail(googleUserInfo.Email)

	if existingUser == nil || err != nil {
		randomPassword, genErr := factories.Generate32ByteToken()

		if genErr != nil {
			log.Println("Error generating random password", genErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user's data"})
			return
		}

		isVerified := true
		isProfileComplted := false

		// TODO: Create transaction for two insert update
		createdUser, userErr := h.UserService.CreateUser(googleUserInfo.GivenName, googleUserInfo.Email, googleUserInfo.Email, randomPassword, isVerified, isProfileComplted)

		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
			return
		}

		oauthErr := h.OAuthService.CreateGoogleAccount(createdUser, googleUserInfo, token)

		if oauthErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save authorized account"})
			return
		}

		signUpUrl := fmt.Sprintf("%s/auth/signup?email=%s&fullname=%s", cryAppUrl, createdUser.Email, createdUser.Fullname)
		c.Redirect(302, signUpUrl)
		return
	}

	// Check if oauth account already exists and if user has not completed their profile

	oauthErr := h.OAuthService.CreateGoogleAccount(existingUser, googleUserInfo, token)

	if oauthErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save the authorized account"})
		return
	}

	jwt, err := jwtService.GenerateJWT(existingUser.UUID, existingUser.Email, existingUser.TwoFAEnabled, existingUser.TwoFAEnabled)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate jwt"})
	}

	isSecured := appEnv == "production"

	c.SetCookie("jwt", jwt, 3000, "/", cryAppUrl, isSecured, true)
	c.Redirect(302, cryAppUrl)
}
