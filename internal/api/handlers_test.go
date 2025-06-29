package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ravwvil/task-manager/internal/models"
	"github.com/Ravwvil/task-manager/internal/service"
	"github.com/Ravwvil/task-manager/internal/storage"
)

func TestCreateTask(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	reqBody := models.CreateTaskRequest{Description: "Test task"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	var response models.TaskResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Status != models.TaskStatusPending && response.Status != models.TaskStatusRunning {
		t.Errorf("Expected status %s or %s, got %s", models.TaskStatusPending, models.TaskStatusRunning, response.Status)
	}
}

func TestCreateTaskEmptyDescription(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	reqBody := models.CreateTaskRequest{Description: ""}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestGetTask(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	task := models.NewTask("Test task")
	storage.Create(task)

	req := httptest.NewRequest("GET", "/api/v1/tasks/"+task.ID, nil)
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response models.TaskResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.ID != task.ID {
		t.Errorf("Expected task ID %s, got %s", task.ID, response.ID)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/v1/tasks/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}
}

func TestListTasks(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	task1 := models.NewTask("Test task 1")
	task2 := models.NewTask("Test task 2")
	storage.Create(task1)
	storage.Create(task2)

	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response []*models.TaskResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(response))
	}
}

func TestDeleteTask(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	task := models.NewTask("Test task")
	storage.Create(task)

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+task.ID, nil)
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
	}

	_, err := storage.GetByID(task.ID)
	if err == nil {
		t.Error("Expected task to be deleted")
	}
}

func TestHealthCheck(t *testing.T) {
	storage := storage.NewInMemoryTaskStorage()
	service := service.NewTaskService(storage)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", response["status"])
	}
}
