package main

import (
	"fastlink/src/api"
	"fastlink/src/auth"
	"fastlink/src/config"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Starting FastLink server...")

	router := gin.Default()

	router.POST("/register", api.Register)
	router.POST("/login", api.Login)
	router.GET("/refresh", auth.ParseToken, auth.AuthRefreshToken, api.Refresh)


	
	adminGroup := router.Group("/admin")
	adminGroup.Use(auth.ParseToken, auth.AuthAdmin)
	{
		adminGroup.GET("/users", api.GetAllUser)
		adminGroup.GET("/links", api.GetAllLink)
		adminGroup.POST("/user/ban", api.BanUser)

	}

	ServerConfig := config.Server()
	router.Run(ServerConfig.PortStr())
}
