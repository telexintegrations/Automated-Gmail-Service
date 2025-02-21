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

func LoginTelex(c *gin.Context) {
	var loginReq TelexRequestBody

	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
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

	log.Printf("Username: %s", username)

	message := loginReq.Message
	log.Println("Message received: ", message)

	var formattedMessage string = StripHTMLTags(message)
	log.Println("Formatted Message received: ", formattedMessage)

	if formattedMessage == "/start-mail" {
		if username == "" || email == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed. Ensure username, email and password are set.", "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
			return
		}

		conn, err := ConnectToImapWithPassword(email, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"})
			return
		}
		defer conn.Logout()

		go EmailNoAuthHandler(email, password, username)

		log.Println("User logged in: ", email)
		log.Println("Email monitoring service started successfully.")
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful. Email monitoring started. New inbox mails would receive automated responses.", "username": "Automated Email Service", "event_name": "Handling Emails"})
	} else {
		log.Println("Type /start-mail to start email monitoring service.")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Type /start-mail to start email monitoring service.", "username": "Automated Email Service", "event_name": "Handling Emails"})
	}
}
