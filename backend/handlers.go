package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addTask(context *gin.Context) {
	var newTask Task

	if err := context.BindJSON(&newTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	context.JSON(http.StatusCreated, gin.H{"data": newTask})
}

func getTask(context *gin.Context) {
	id := context.Param("id")
	task, err := getTaskById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"data": task})
}

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, gin.H{"data": tasks})
}
