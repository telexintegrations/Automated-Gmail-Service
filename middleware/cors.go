package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"https://telex.im", "http://localhost:8080",
			"https://staging.telex.im",
			"http://telextest.im",
			"http://staging.telextest.im"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
