package tests

import (
	"bytes"
	"encoding/json"
	"hng-stage3-task-automated-email-service/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func mockConnectToImapWithPassword(email, password string) (string, error) {
// 	if email == "valid@gmail.com" && password == "password123" {
// 		return "mock_connection", nil
// 	}
// 	return "", assert.AnError
// }

func TestLoginTelex_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/target_url", handlers.LoginTelex)

	requestBody := map[string]any{
		"message": "/start-mail",
		"settings": []map[string]any{
			{"label": "username", "type": "text", "required": true, "default": "testuser"},
			{"label": "email", "type": "text", "required": true, "default": "valid@gmail.com"},
			{"label": "password", "type": "text", "required": true, "default": "password123"},
		},
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/target_url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
}

func TestLoginTelex_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/target_url", handlers.LoginTelex)

	requestBody := map[string]any{
		"message": "/start-mail",
		"settings": []map[string]any{
			{"label": "username", "type": "text", "required": true, "default": "testuser"},
			// Left out email on purpose...
			{"label": "password", "type": "text", "required": true, "default": "password123"},
		},
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/target_url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Login failed. Ensure username, email and password are set.")
}

func TestLoginTelex_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/target_url", handlers.LoginTelex)

	// Invalid JSON request...
	req, _ := http.NewRequest("POST", "/target_url", bytes.NewBuffer([]byte(`invalid_json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}

func TestLoginTelex_AuthFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/target_url", handlers.LoginTelex)

	// Invalid login credentials...
	requestBody := map[string]any{
		"message": "/start-mail",
		"settings": []map[string]any{
			{"label": "username", "type": "text", "required": true, "default": "testuser"},
			{"label": "email", "type": "text", "required": true, "default": "invalid@gmail.com"},
			{"label": "password", "type": "text", "required": true, "default": "wrongpassword"},
		},
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/target_url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication failed")
}

func TestLoginTelex_MissingMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/target_url", handlers.LoginTelex)

	// Request missing message field...
	requestBody := map[string]any{
		"settings": []map[string]any{
			{"label": "username", "type": "text", "required": true, "default": "testuser"},
			{"label": "email", "type": "text", "required": true, "default": "valid@gmail.com"},
			{"label": "password", "type": "text", "required": true, "default": "password123"},
		},
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/target_url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Type /start-mail to start email monitoring service.")
}
