package main

import (
	"demo-server/handlers"
	"demo-server/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"demo-server/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userCenterURL := os.Getenv("USER_CENTER_URL")
	if userCenterURL == "" {
		userCenterURL = "http://localhost:8082"
	}

	sysName := os.Getenv("SYS_NAME")
	if sysName == "" {
		sysName = "demo"
	}

	jwtHashKey := os.Getenv("JWT_HASH_KEY")
	if jwtHashKey == "" {
		jwtHashKey = "1040f1b0fa1ef69d804de8d5bc996830e7f049e"
	}
	utils.SetJwtHashKey(jwtHashKey)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.GET("/health", handlers.HealthCheck)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
	}

	user := r.Group("/api/user")
	{
		user.Use(middleware.AuthMiddleware())
		user.GET("/profile", handlers.GetProfile)
		user.PUT("/profile", handlers.UpdateProfile)
	}

	userCenter := handlers.NewUserCenterHandler(userCenterURL, sysName)
	uc := r.Group("/api/uc")
	{
		uc.GET("/login/methods", userCenter.GetLoginMethods)
		uc.GET("/login/qrcode", userCenter.GetWxQrcode)
		uc.GET("/login/check-scan", userCenter.CheckWxScanStatus)
		uc.GET("/login/check-token", userCenter.CheckToken)
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
