package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, gin.H{"data": tasks})
}

func addTask(context *gin.Context) {
	var newTask Task

	if err := context.BindJSON(&newTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	context.JSON(http.StatusCreated, gin.H{"data": newTask})
}
