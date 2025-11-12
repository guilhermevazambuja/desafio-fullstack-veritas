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

	testTask := Task{
		ID:        "4",
		Title:     "Write Documentation",
		Completed: false,
	}

	body, err := json.Marshal(testTask)
	require.NoError(t, err)

	w := performRequest(router, http.MethodPost, "/tasks", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assertJSON(t, w)

	var resp SuccessResponse[Task]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Deterministic expectations
	want := testTask
	got := resp.Data
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)
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

	// Deterministic expectations
	want := tasks[0]
	got := resp.Data
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)
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

	// Deterministic expectations
	require.NotEmpty(t, resp.Data)
	assert.Equal(t, len(tasks), len(resp.Data))

	// Validate first item fields
	got := resp.Data[0]
	want := tasks[0]
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)

	// Validate last item fields
	want = tasks[len(tasks)-1]
	got = resp.Data[len(resp.Data)-1]
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)
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
