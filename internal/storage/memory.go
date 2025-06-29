package storage

import (
	"fmt"
	"sync"

	"github.com/Ravwvil/task-manager/internal/models"
)

type TaskStorage interface {
	Create(task *models.Task) error
	GetByID(id string) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id string) error
	List() ([]*models.Task, error)
}

type InMemoryTaskStorage struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	return &InMemoryTaskStorage{
		tasks: make(map[string]*models.Task),
	}
}

func (s *InMemoryTaskStorage) Create(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryTaskStorage) GetByID(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

func (s *InMemoryTaskStorage) Update(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryTaskStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	delete(s.tasks, id)
	return nil
}

func (s *InMemoryTaskStorage) List() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}
