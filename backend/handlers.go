package main

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func addTask(context *gin.Context) {
	var newTask Task
	err := context.ShouldBindJSON(&newTask)

	statusCode := http.StatusInternalServerError
	if err != nil {
		statusCode = http.StatusBadRequest
		context.IndentedJSON(statusCode, ErrorResponse{Error: ErrInvalidPayload.Error()})
		return
	}
	statusCode = http.StatusCreated

	// TODO Add validation steps
	tasks = append(tasks, newTask)
	context.IndentedJSON(statusCode, SuccessResponse[Task]{Data: newTask})
}

func getTask(context *gin.Context) {
	id := context.Param("id")
	task, err := getTaskById(id)

	statusCode := http.StatusInternalServerError
	if err != nil {
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

func replaceTask(context *gin.Context) {
	id := context.Param("id")
	var updatedTask Task

	errBindJSON := context.ShouldBindJSON(&updatedTask)
	if errBindJSON != nil {
		context.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidPayload.Error()})
		return
	}

	if updatedTask.ID != nil && *updatedTask.ID != id {
		context.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrIDMismatch.Error()})
		return
	}

	existingTask, errGetTaskById := getTaskById(id)
	if errGetTaskById != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(errGetTaskById, ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		}
		context.IndentedJSON(statusCode, ErrorResponse{Error: errGetTaskById.Error()})
		return
	}

	vTask := reflect.ValueOf(updatedTask)
	for i := 0; i < vTask.NumField(); i++ {
		if vTask.Field(i).IsNil() {
			context.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrIncompletePayload.Error()})
			return
		}
	}

	*existingTask = updatedTask
	context.IndentedJSON(http.StatusOK, SuccessResponse[Task]{Data: *existingTask})
}

func updateTask(context *gin.Context) {
	id := context.Param("id")
	var updatedProps Task

	errBindJSON := context.ShouldBindJSON(&updatedProps)
	if errBindJSON != nil {
		context.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidPayload.Error()})
		return
	}

	if updatedProps.ID != nil && *updatedProps.ID != id {
		context.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrIDMismatch.Error()})
		return
	}

	existingTask, errGetTaskById := getTaskById(id)
	if errGetTaskById != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(errGetTaskById, ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		}
		context.IndentedJSON(statusCode, ErrorResponse{Error: errGetTaskById.Error()})
		return
	}

	vTask := reflect.ValueOf(existingTask).Elem()
	vUpdates := reflect.ValueOf(updatedProps)
	for i := 0; i < vTask.NumField(); i++ {
		fieldTask := vTask.Field(i)
		fieldUpdate := vUpdates.Field(i)

		if !fieldUpdate.IsNil() {
			fieldTask.Set(fieldUpdate)
		}
	}
	context.IndentedJSON(http.StatusOK, SuccessResponse[Task]{Data: *existingTask})
}
