package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addTask(context *gin.Context) {
	var newTask Task

	if err := context.BindJSON(&newTask); err != nil {
		context.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidPayload.Error()})
		return
	}

	tasks = append(tasks, newTask)
	context.JSON(http.StatusCreated, SuccessResponse[Task]{Data: newTask})
}

func getTask(context *gin.Context) {
	id := context.Param("id")
	task, err := getTaskById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: ErrTaskNotFound.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, SuccessResponse[Task]{Data: *task})
}

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, SuccessResponse[[]Task]{Data: tasks})
}
