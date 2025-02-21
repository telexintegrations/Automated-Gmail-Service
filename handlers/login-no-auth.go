package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type LoginNoAuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginNoOauthHandler(c *gin.Context) {
	var loginReq LoginNoAuthRequest

	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}

	if loginReq.Username == "" || loginReq.Email == "" || loginReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login failed. Ensure username, email and password are set.", "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}

	conn, err := ConnectToImapWithPassword(loginReq.Email, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}
	defer conn.Logout()

	go EmailNoAuthHandler(loginReq.Email, loginReq.Password, loginReq.Username)

	log.Println("User logged in:", loginReq.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful. Email monitoring started. New inbox mails would receive automated responses.", "username": "Automated Email Service", "status": "success", "event_name": "Handling Emails"})
}
