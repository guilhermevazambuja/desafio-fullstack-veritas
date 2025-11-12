package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/tasks", addTask)
	router.GET("/tasks/:id", getTask)
	router.GET("/tasks", getTasks)

	err := router.Run("localhost:9090")
	if err != nil {
		return
	}
}
