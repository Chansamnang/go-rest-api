package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-rest-api/controller"
	config "go-rest-api/internal/config"
	"go-rest-api/lang"
	"go-rest-api/middleware"
	"go-rest-api/pkg"
	"log"
	"os"
)

func main() {

	var err error

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	config.Init()

	fmt.Println("App config successful")

	fmt.Printf("config %v %v %v", config.Config.AppConfig, config.Config.MySQLConfig, config.Config.RedisConfig)

	config.InitDB()
	config.InitRedis()
	config.LoggerInit()
	lang.Init()

	buildHandler()
}

func buildHandler() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "192.168.60.211"})

	v1 := router.Group("/api")
	// Middleware
	v1.Use(middleware.Cors())
	v1.Use(middleware.LoggerHandlerMiddleware())
	v1.Use(middleware.AuthMiddleware())
	v1.Use(middleware.RateLimitMiddleware())

	controller.UserRegisterHandlers(v1)

	pkg.RouteFetcher(router)

	public := router.Group("/api")
	controller.PublicRegisterHandler(public)

	port := os.Getenv("PORT")
	err := router.Run(port)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Server Running Error %v", err.Error()))
		return
	}
}
