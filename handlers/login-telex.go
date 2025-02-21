package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"bytes"
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

func sendWebhookNotification(payload gin.H, webhook string) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending webhook request:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Webhook notification sent successfully. Status Code:", resp.StatusCode)
}

func LoginTelex(c *gin.Context) {
	var loginReq TelexRequestBody

	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		response := gin.H{"error": "Invalid request: " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var username, password, email, webhook string

	for _, setting := range loginReq.Settings {
		switch setting.Label {
		case "username":
			username = setting.Default
		case "email":
			email = setting.Default
		case "password":
			password = setting.Default
		case "webhook":
			webhook = setting.Default
		}
	}

	log.Printf("Username: %s", username)

	message := loginReq.Message
	log.Println("Message received: ", message)

	var formattedMessage string = StripHTMLTags(message)
	log.Println("Formatted Message received: ", formattedMessage)

	if formattedMessage == "/start-mail" {
		if username == "" || email == "" || password == "" {
			response := gin.H{"message": "Login failed. Ensure username, email and password are set.", "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"}
			sendWebhookNotification(response, webhook)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		conn, err := ConnectToImapWithPassword(email, password)
		if err != nil {
			response := gin.H{"message": "Authentication failed " + err.Error(), "status": "error", "username": "Automated Email Service", "event_name": "Handling Emails"}
			c.JSON(http.StatusUnauthorized, response)
			return
		}
		defer conn.Logout()

		go EmailNoAuthHandler(email, password, username)

		log.Println("User logged in: ", email)
		log.Println("Email monitoring service started successfully.")
		response := gin.H{"status": "success", "message": "Login successful. Email monitoring started. New inbox mails would receive automated responses.", "username": "Automated Email Service", "event_name": "Handling Emails"}
		c.JSON(http.StatusOK, response)
	} else {
		log.Println("Type /start-mail to start email monitoring service.")
		response := gin.H{"status": "error", "message": "Type /start-mail to start email monitoring service.", "username": "Automated Email Service", "event_name": "Handling Emails"}
		c.JSON(http.StatusBadRequest, response)
	}
}
