package handlers

import (
	"context"
	"fmt"
	"hng-stage3-task-automated-email-service/config"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var credentials = config.LoadConfig()

var clientId = credentials.GoogleClientID
var clientSecret = credentials.GoogleClientSecret
var redirectUrl = credentials.GoogleRedirectURI

var googleOAuthConfig = &oauth2.Config{
	ClientID:     clientId,
	ClientSecret: clientSecret,
	RedirectURL:  redirectUrl,
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

// var tokenStore = make(map[string]*oauth2.Token)

func LoginHandler(c *gin.Context) {
	fmt.Println(clientId, clientSecret, redirectUrl)
	if googleOAuthConfig.ClientID == "" || googleOAuthConfig.ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth configuration missing. Ensure client_id and client_secret are set."})
		return
	}
	authUrl := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	if authUrl == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication URL"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, authUrl)
}

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code missing"})
		return
	}
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	userEmail, err := GetEmailFromToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user email"})
		return
	}

	c.SetCookie("auth-token", token.AccessToken, 3600, "/", "", true, true)
	c.SetCookie("refresh-token", token.RefreshToken, 0, "/", "", true, true)

	// tokenStore["user"] = token

	// token, exists := tokenStore["user"]
	// if !exists || token.AccessToken == "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
	// 	return
	// }

	// if token.Expiry.Before(time.Now()) {
	// 	newToken, err := RefreshAccessToken(token.RefreshToken)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
	// 		return
	// 	}
	// 	tokenStore["user"] = newToken
	// 	token = newToken
	// }

	go EmailHandler(userEmail, token.AccessToken)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful. Email listener started", "email": userEmail})
}

func RefreshAccessToken(refreshToken string) (*oauth2.Token, error) {
	tokenSource := googleOAuthConfig.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}
