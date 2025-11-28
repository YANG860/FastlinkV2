package main

import (
	"fastlink/src/api"
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
	router.GET("/refresh", api.Refresh)

	router.GET("/users", api.GetAllUser)
	router.GET("/links", api.GetAllLink)
	router.POST("/user/ban", api.BanUser)
	

	
	ServerConfig := config.Server()
	router.Run(ServerConfig.PortStr())
}
