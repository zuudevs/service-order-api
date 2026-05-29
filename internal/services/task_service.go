/**

 filename  : task_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package services

import (
	"errors"
	"strings"
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/repositories"
)

type TaskService struct {
	taskRepo repositories.TaskRepository
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(
	subject string,
	description string,
	price uint64,
	due time.Time,
) error {
	subject = strings.TrimSpace(subject)

	if subject == "" {
		return errors.New("subject is required")
	}

	if due.Before(time.Now().UTC()) {
		return errors.New("due date must be in the future")
	}

	task := models.NewTask(subject, description, price, due)

	return s.taskRepo.Create(task)
}

func (s *TaskService) Index() ([]models.Task, error) {
	return s.taskRepo.Index()
}

func (s *TaskService) GetByID(id uint64) (*models.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *TaskService) GetByStatus(status models.TaskStatus) ([]models.Task, error) {
	return s.taskRepo.GetByStatus(status)
}

func (s *TaskService) GetBySubject(subject string) ([]models.Task, error) {
	return s.taskRepo.GetBySubject(subject)
}

func (s *TaskService) UpdateStatus(id uint64, status models.TaskStatus) error {
	_, err := s.taskRepo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.taskRepo.SetStatus(id, status)
}

func (s *TaskService) UpdatePrice(id uint64, price uint64) error {
	_, err := s.taskRepo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.taskRepo.SetPrice(id, price)
}

func (s *TaskService) UpdateSubject(id uint64, subject string) error {
	subject = strings.TrimSpace(subject)

	if subject == "" {
		return errors.New("subject is required")
	}

	_, err := s.taskRepo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.taskRepo.SetSubject(id, subject)
}

func (s *TaskService) DeleteTask(id uint64) error {
	_, err := s.taskRepo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.taskRepo.Delete(id)
}
