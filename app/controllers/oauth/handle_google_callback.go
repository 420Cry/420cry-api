package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	cfg "cry-api/app/config"
	jwtService "cry-api/app/services/jwt"
)


func (h *OAuthController) HandleGoogleCallback(c *gin.Context) {
	log.Println("Go into the callback")
	code := c.Query("code")
	cryAppUrl := cfg.Get().CryAppURL
	appEnv := cfg.Get().AppEnv
	isSecured := appEnv == "production"

	if code == "" {
		c.Redirect(302, fmt.Sprintf("%s/auth/signup", cryAppUrl))
		return
	}

	token, err := h.OAuthService.ExchangeToken(context.Background(), code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed oauth exchange"})
		return
	}

	googleUserInfo, err := h.OAuthService.FetchUserInfo(token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info from Google"})
		return
	}

	existingUser, err := h.UserService.FindUserByEmail(googleUserInfo.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	if existingUser == nil {
		createdUser, transactionErr := h.UserService.CreateUserByGoogleAuth(googleUserInfo, token)

		if transactionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
			return
		}

		signUpUrl := fmt.Sprintf("%s/auth/signup?email=%s&fullname=%s", cryAppUrl, createdUser.Email, createdUser.Fullname)
		c.Redirect(302, signUpUrl)
		return
	}

	oauthAccount, err := h.OAuthService.FindAccountByProviderAndId("Google", googleUserInfo.Sub)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find user's google account"})
		return
	}

	if !existingUser.IsProfileCompleted && oauthAccount != nil {
		signUpUrl := fmt.Sprintf("%s/auth/signup?email=%s&fullname=%s", cryAppUrl, existingUser.Email, existingUser.Fullname)
		c.Redirect(302, signUpUrl)
		return
	}

	if oauthAccount == nil {
		if err := h.OAuthService.CreateGoogleAccount(existingUser, googleUserInfo, token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save the authorized account"})
			return
		}

		jwt, err := jwtService.GenerateJWT(existingUser.UUID, existingUser.Email, existingUser.TwoFAEnabled, existingUser.TwoFAEnabled)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate jwt"})
		}

		// Wont work in local
		c.SetCookie("jwt", jwt, 3000, "/", "localhost", isSecured, true)
		c.Redirect(302, cryAppUrl)
		return
	}

	jwt, err := jwtService.GenerateJWT(existingUser.UUID, existingUser.Email, existingUser.TwoFAEnabled, existingUser.TwoFAEnabled)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate jwt"})
	}

	// Wont work in local
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		MaxAge:   3600,
		Secure:   isSecured,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	c.Redirect(302, cryAppUrl)
}
