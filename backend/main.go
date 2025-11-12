package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoints
	router.POST("/tasks", addTask)
	router.GET("/tasks/:id", getTask)
	router.GET("/tasks", getTasks)
	router.PUT("/tasks/:id", replaceTask)
	router.PATCH("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)

	if err := router.Run("localhost:9090"); err != nil {
		panic(err)
	}
}
