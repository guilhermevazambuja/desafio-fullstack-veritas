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
	r.GET("/tasks", getTasks)
	r.POST("/tasks", addTask)
	return r
}

// Test listing all tasks
func TestGetTasks(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	router.ServeHTTP(w, req)

	// Status and content type
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Typed JSON decode
	var resp ListResp
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
	got = resp.Data[len(resp.Data)-1]
	want = tasks[len(tasks)-1]
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)
}

func TestAddTask(t *testing.T) {
	router := setupRouter()

	testTask := Task{
		ID:        "4",
		Title:     "Write Documentation",
		Completed: false,
	}

	body, err := json.Marshal(testTask)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body into map
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate returned task fields
	data := response["data"].(map[string]interface{})
	assert.Equal(t, testTask.ID, data["id"])
	assert.Equal(t, testTask.Title, data["title"])
	assert.Equal(t, testTask.Completed, data["completed"])
}
