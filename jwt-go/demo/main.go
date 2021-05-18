package main

import (
	"github.com/gin-gonic/gin"
	"hello/WeeklyLearning/jwt-go/demo/controller"
	"hello/WeeklyLearning/jwt-go/demo/middleware"
)

func initRouter()*gin.Engine {
	router := gin.Default()
	router.POST("/login",controller.Login)
	Auth := router.Group("/auth")
	Auth.Use(middleware.ValidateWelcome)
	{
		Auth.GET("/welcome",controller.Welcome)
	}

	return router
}

func main() {
	router := initRouter()
	router.Run(":8080")
}
