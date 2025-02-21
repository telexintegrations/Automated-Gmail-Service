package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"hng-stage3-task-automated-email-service/handlers"
)

func TestReturnIntegrationJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/integration", handlers.ReturnIntegrationJSON)

	req, _ := http.NewRequest("GET", "/integration", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseData handlers.ResponseData
	err := json.Unmarshal(resp.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseData.Data["author"], "Author field should not be empty")
	assert.NotEmpty(t, responseData.Data["settings"], "Settings should not be empty")
	assert.Equal(t, "Automated Email Service", responseData.Data["descriptions"].(map[string]any)["app_name"])
}
