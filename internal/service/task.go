package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Ravwvil/task-manager/internal/models"
	"github.com/Ravwvil/task-manager/internal/storage"
)

type TaskService interface {
	CreateTask(ctx context.Context, description string) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context) ([]*models.Task, error)
}

type taskService struct {
	storage           storage.TaskStorage
	runningTasks      map[string]chan struct{}
	runningTasksMutex sync.RWMutex
}

func NewTaskService(storage storage.TaskStorage) TaskService {
	return &taskService{
		storage:      storage,
		runningTasks: make(map[string]chan struct{}),
	}
}

func (s *taskService) CreateTask(ctx context.Context, description string) (*models.Task, error) {
	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("task description cannot be empty")
	}

	task := models.NewTask(description)

	if err := s.storage.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	stopChan := make(chan struct{})
	s.runningTasksMutex.Lock()
	s.runningTasks[task.ID] = stopChan
	s.runningTasksMutex.Unlock()

	go s.processTask(task.ID, stopChan)

	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.storage.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	_, err := s.storage.GetByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	s.runningTasksMutex.Lock()
	if stopChan, exists := s.runningTasks[id]; exists {
		close(stopChan)
		delete(s.runningTasks, id)
	}
	s.runningTasksMutex.Unlock()

	if err := s.storage.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (s *taskService) ListTasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := s.storage.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) processTask(taskID string, stopChan chan struct{}) {
	defer func() {
		s.runningTasksMutex.Lock()
		delete(s.runningTasks, taskID)
		s.runningTasksMutex.Unlock()
	}()

	task, err := s.storage.GetByID(taskID)
	if err != nil {
		return
	}

	task.Start()
	if err := s.storage.Update(task); err != nil {
		return
	}

	minDuration := 3 * time.Minute
	maxDuration := 5 * time.Minute
	duration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration)))

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-stopChan:
			task.Cancel("Task was cancelled")
			s.storage.Update(task)
			return

		case <-ticker.C:
			if time.Since(startTime) >= duration {
				result := fmt.Sprintf("Task completed successfully after %v", time.Since(startTime).Round(time.Second))
				task.Complete(result)
				s.storage.Update(task)
				return
			}

			if task.StartedAt != nil {
				task.Duration = int64(time.Since(*task.StartedAt).Seconds())
			}

		case <-time.After(duration):
			result := fmt.Sprintf("Task completed successfully after %v", time.Since(startTime).Round(time.Second))
			task.Complete(result)
			s.storage.Update(task)
			return
		}
	}
}
