package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/tasks", addTask)
	r.GET("/tasks/:id", getTask)
	r.GET("/tasks", getTasks)
	r.PUT("/tasks/:id", replaceTask)
	r.PATCH("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	return r
}

// Test adding a new task
func TestAddTask(t *testing.T) {
	resetTasks()
	router := setupRouter()

	testTask := Task{
		Title:  strPtr("Write Documentation"),
		Status: strPtr("to_do"),
	}
	body, err := json.Marshal(testTask)
	require.NoError(t, err)

	w := performRequest(router, http.MethodPost, "/tasks", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.NotNil(t, resp.Data.ID)
	assertTaskMatch(t, testTask, resp.Data)
}

// Test adding a task with an invalid request body
func TestAddTaskInvalidPayload(t *testing.T) {
	resetTasks()
	router := setupRouter()

	invalidJSON := []byte(`{"ids": "5", "title": 123}`)

	w := performRequest(router, http.MethodPost, "/tasks", invalidJSON)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertJSON(t, w)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, ErrInvalidPayload.Error(), resp.Error)
}

// Test getting a specific task
func TestGetTask(t *testing.T) {
	resetTasks()
	router := setupRouter()

	w := performRequest(router, http.MethodGet, "/tasks/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, tasks[0], resp.Data)
}

// Test getting a task that doesn't exist
func TestGetTaskNotFound(t *testing.T) {
	resetTasks()
	router := setupRouter()

	invalidId := "nonexistent-id"

	w := performRequest(router, http.MethodGet, "/tasks/"+invalidId, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assertJSON(t, w)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, ErrTaskNotFound.Error(), resp.Error)
}

// Test listing all tasks
func TestGetTasks(t *testing.T) {
	resetTasks()
	router := setupRouter()

	w := performRequest(router, http.MethodGet, "/tasks", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[[]Task]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, len(tasks), len(resp.Data))

	// Validate first item fields
	assert.Equal(t, tasks[0], resp.Data[0])

	// Validate last item fields
	assert.Equal(t, tasks[len(tasks)-1], resp.Data[len(resp.Data)-1])

}

// Test fully replacing a specific task data
func TestReplaceTask(t *testing.T) {
	resetTasks()
	router := setupRouter()

	testTask := Task{
		ID:     strPtr("1"),
		Title:  strPtr("Study Go"),
		Status: strPtr("in_progress"),
	}
	body, err := json.Marshal(testTask)
	require.NoError(t, err)

	originalTask, _ := getTaskById("1")
	assert.NotEqual(t, "Study Go", *originalTask.Title)

	w := performRequest(router, http.MethodPut, "/tasks/1", body)

	assert.Equal(t, http.StatusOK, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, testTask, resp.Data)

	updatedTask, _ := getTaskById("1")
	assert.Equal(t, testTask, *updatedTask)
}

// Test replacing task data with an incomplete request
func TestReplaceTaskIncompletePayload(t *testing.T) {
	resetTasks()
	router := setupRouter()

	badTask := Task{
		Title: strPtr("Incomplete request"),
	}
	body, err := json.Marshal(badTask)
	require.NoError(t, err)

	originalTask, _ := getTaskById("1")

	w := performRequest(router, http.MethodPut, "/tasks/1", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertJSON(t, w)

	var resp ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, ErrIncompletePayload.Error(), resp.Error)

	unchangedTask, _ := getTaskById("1")
	assert.Equal(t, *originalTask, *unchangedTask)
}

// Test partially modifying an existing task
func TestUpdateTaskSuccess(t *testing.T) {
	resetTasks()
	router := setupRouter()

	updatePayload := Task{
		Title: strPtr("Clean the kitchen"),
	}
	body, err := json.Marshal(updatePayload)
	require.NoError(t, err)

	originalTask, _ := getTaskById("1")
	assert.NotEqual(t, "Clean the kitchen", *originalTask.Title)

	w := performRequest(router, http.MethodPatch, "/tasks/1", body)

	assert.Equal(t, http.StatusOK, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "Clean the kitchen", *resp.Data.Title)
	assert.Equal(t, *originalTask.Status, *resp.Data.Status)

	updatedTask, _ := getTaskById("1")
	assert.Equal(t, "Clean the kitchen", *updatedTask.Title)
}

// Test updating a task with mismatched ID in payload
func TestUpdateTaskIDMismatch(t *testing.T) {
	resetTasks()
	router := setupRouter()

	updatePayload := Task{
		ID:    strPtr("999"),
		Title: strPtr("Should fail"),
	}
	body, err := json.Marshal(updatePayload)
	require.NoError(t, err)

	w := performRequest(router, http.MethodPatch, "/tasks/1", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertJSON(t, w)

	var resp ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, ErrIDMismatch.Error(), resp.Error)

	unchangedTask, _ := getTaskById("1")
	assert.NotEqual(t, "Should fail", *unchangedTask.Title)
}

// Test deleting an existing task successfully
func TestDeleteTaskSuccess(t *testing.T) {
	resetTasks()
	router := setupRouter()

	taskToDeleteId := "3"

	w := performRequest(router, http.MethodDelete, "/tasks/"+taskToDeleteId, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Verifica que a resposta contém a task deletada
	assert.Equal(t, "Record Video", *resp.Data.Title)
	assert.Equal(t, taskToDeleteId, *resp.Data.ID)

	// Verifica que a task foi realmente removida do slice
	_, errGet := getTaskById(taskToDeleteId)
	assert.Error(t, errGet)
	assert.True(t, errors.Is(errGet, ErrTaskNotFound))
}

// Test deleting a non-existent task
func TestDeleteTaskNotFound(t *testing.T) {
	resetTasks()
	router := setupRouter()

	invalidId := "nonexistent-id"

	w := performRequest(router, http.MethodDelete, "/tasks/"+invalidId, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assertJSON(t, w)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrTaskNotFound.Error(), resp.Error)
}

func performRequest(r http.Handler, method string, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func assertTaskMatch(t *testing.T, expected, actual Task) {
	assert.Equal(t, *expected.Title, *actual.Title)
	assert.Equal(t, *expected.Status, *actual.Status)
}

func assertJSON(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func resetTasks() {
	tasks = []Task{
		{ID: strPtr("1"), Title: strPtr("Clean Room"), Status: strPtr("to_do")},
		{ID: strPtr("2"), Title: strPtr("Read Book"), Status: strPtr("in_progress")},
		{ID: strPtr("3"), Title: strPtr("Record Video"), Status: strPtr("done")},
	}
}
