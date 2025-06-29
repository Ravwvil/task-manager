package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ravwvil/task-manager/internal/models"
	"github.com/Ravwvil/task-manager/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	taskService service.TaskService
}

func NewHandler(taskService service.TaskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/tasks", h.createTask).Methods("POST")
	api.HandleFunc("/tasks", h.listTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.getTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.deleteTask).Methods("DELETE")

	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	router.Use(h.panicRecoveryMiddleware)
	router.Use(h.loggingMiddleware)
	router.Use(h.corsMiddleware)

	return router
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), req.Description)
	if err != nil {
		if err.Error() == "task description cannot be empty" {
			h.errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.errorResponse(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	h.jsonResponse(w, http.StatusCreated, &models.TaskResponse{Task: task})
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.errorResponse(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	task, err := h.taskService.GetTask(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusNotFound, "Task not found")
		return
	}

	h.jsonResponse(w, http.StatusOK, &models.TaskResponse{Task: task})
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.errorResponse(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	if err := h.taskService.DeleteTask(r.Context(), id); err != nil {
		h.errorResponse(w, http.StatusNotFound, "Task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.ListTasks(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]*models.TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = &models.TaskResponse{Task: task}
	}

	h.jsonResponse(w, http.StatusOK, response)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	h.jsonResponse(w, http.StatusOK, response)
}

func (h *Handler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) errorResponse(w http.ResponseWriter, status int, message string) {
	h.jsonResponse(w, status, &models.ErrorResponse{Error: message})
}
