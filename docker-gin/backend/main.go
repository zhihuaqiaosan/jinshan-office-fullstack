package main

import (
	"os"
	"wordbook/handlers"
	"wordbook/middleware"
	"wordbook/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	models.InitDB()

	r := gin.Default()

	// 公开路由
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	// 需要登录的路由
	auth := r.Group("/api")
	auth.Use(middleware.Auth())
	{
		auth.GET("/word/query", handlers.QueryWord)
		auth.POST("/word/save", handlers.SaveWord)
		auth.GET("/words", handlers.GetWords)
		auth.DELETE("/word/:id", handlers.DeleteWord)
	}

	r.Run(":" + os.Getenv("BACKEND_PORT"))
}
