package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TelexRequestBody struct {
	Message  string    `json:"message"`
	Settings []Setting `json:"settings"`
}

type Setting struct {
	Label    string `json:"label"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  string `json:"default"`
}

type LoginTelexRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginTelex(c *gin.Context) {
	var loginReq TelexRequestBody

	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}

	var username, password, email string

	for _, setting := range loginReq.Settings {
		switch setting.Label {
		case "username":
			username = setting.Default
		case "email":
			email = setting.Default
		case "password":
			password = setting.Default
		}
	}

	message := loginReq.Message

	log.Println("Message received: ", message)

	if username == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login failed. Ensure username, email and password are set.", "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}

	conn, err := ConnectToImapWithPassword(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
		return
	}
	defer conn.Logout()

	go EmailNoAuthHandler(email, password, username)

	log.Println("User logged in:", email)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful. Email monitoring started. New inbox mails would receive automated responses.", "username": "Automated Email Service", "status": "success", "event_name": "Handling Emails"})
}
