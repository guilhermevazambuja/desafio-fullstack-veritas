package main

import (
	"bytes"
	"encoding/json"
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
	return r
}

// Test adding a new task
func TestAddTask(t *testing.T) {
	router := setupRouter()

	testTask := Task{ID: "4", Title: "Write Documentation", Completed: false}
	body, err := json.Marshal(testTask)
	require.NoError(t, err)

	w := performRequest(router, http.MethodPost, "/tasks", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, testTask, resp.Data)
}

// Test adding a task with an invalid request body
func TestAddTaskInvalidPayload(t *testing.T) {
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

func performRequest(r http.Handler, method string, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func assertJSON(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}
