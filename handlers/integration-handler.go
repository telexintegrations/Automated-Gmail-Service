package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Data map[string]any `json:"data"`
}

var responseData = ResponseData{
	Data: map[string]any{
		"date": map[string]string{
			"created_at": "2025-02-20",
			"updated_at": "2025-02-21",
		},
		"descriptions": map[string]any{
			"app_name":         "Automated Email Service",
			"app_description":  "This is an automated email service that when integrated to a user's email, sends an automated mail to every new mail sender. It can be activated by typing /start-mail",
			"app_logo":         "https://www.shutterstock.com/image-vector/single-black-email-refresh-line-600nw-2455287007.jpg",
			"app_url":          "https://automated-gmail-service.onrender.com",
			"background_color": "#fff",
		},
		"is_active":            true,
		"integration_type":     "output",
		"integration_category": "Email & Messaging",
		"key_features":         []string{"Sends an automated mail to every new mail sender", "Service is started by typing /start-mail", "Responses are sent to the webhook you provide while configuring the integration"},
		"author":               "Tonyrealzy",
		"website":              "https://automated-gmail-service.onrender.com",
		"settings": []map[string]any{
			{
				"label":    "username",
				"type":     "text",
				"required": true,
				"default":  "",
			},
			{
				"label":    "email",
				"type":     "text",
				"required": true,
				"default":  "",
			},
			{
				"label":    "password",
				"type":     "text",
				"required": true,
				"default":  "",
			},
			{
				"label":    "webhook",
				"type":     "text",
				"required": true,
				"default":  "",
			},
		},
		"target_url": "https://automated-gmail-service.onrender.com/target_url",
	},
}

func ReturnIntegrationJSON(c *gin.Context) {
	c.JSON(http.StatusOK, responseData)
}
