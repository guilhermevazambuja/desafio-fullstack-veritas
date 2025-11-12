package main

import (
	"errors"
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
