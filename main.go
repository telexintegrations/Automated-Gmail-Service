package main

import (
	"fmt"
	"hng-stage3-task-automated-email-service/handlers"
	"hng-stage3-task-automated-email-service/middleware"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Use(middleware.SetUpCORS())

	r.POST("/auth/login", handlers.LoginNoOauthHandler)

	fmt.Println("Server running on port 8080...")
	log.Fatal(r.Run(":8080"))
}
