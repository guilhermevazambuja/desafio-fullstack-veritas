package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addTask(context *gin.Context) {
	var newTask Task
	err := context.ShouldBindJSON(&newTask)

	statusCode := http.StatusInternalServerError
	if err != nil {
		statusCode = http.StatusBadRequest
		context.JSON(statusCode, ErrorResponse{Error: ErrInvalidPayload.Error()})
		return
	}
	statusCode = http.StatusCreated

	// TODO Add validation steps
	tasks = append(tasks, newTask)
	context.JSON(statusCode, SuccessResponse[Task]{Data: newTask})
}

func getTask(context *gin.Context) {
	id := context.Param("id")
	task, err := getTaskById(id)

	statusCode := http.StatusInternalServerError
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: ErrTaskNotFound.Error()})
		if errors.Is(err, ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		}
		context.IndentedJSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}
	statusCode = http.StatusOK

	context.IndentedJSON(statusCode, SuccessResponse[Task]{Data: *task})
}

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, SuccessResponse[[]Task]{Data: tasks})
}
