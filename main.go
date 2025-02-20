package main

import (
	"fmt"
	"hng-stage3-task-automated-email-service/handlers"
	"hng-stage3-task-automated-email-service/middleware"
	"hng-stage3-task-automated-email-service/config"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

	r.Use(middleware.SetUpCORS())

	r.POST("/auth/login", handlers.LoginNoOauthHandler)

	// r.GET("/auth/oauth", handlers.LoginHandler)

	// r.GET("/callback", handlers.CallbackHandler)

	// r.POST("/check-emails", handlers.CheckTokenHandler)

	fmt.Println("Server running on port 8080...")
	log.Fatal(r.Run(":8080"))
}
