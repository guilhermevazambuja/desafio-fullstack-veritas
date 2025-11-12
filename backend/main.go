package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.POST("/tasks", addTask)

	err := router.Run("localhost:9090")
	if err != nil {
		return
	}
}
