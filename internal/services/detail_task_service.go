/**

 filename  : detail_task_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package services

import (
	"errors"

	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/repositories"
)

type DetailTaskService struct {
	detailTaskRepo repositories.DetailTaskRepository
	taskRepo       repositories.TaskRepository
}

func NewDetailTaskService(
	detailTaskRepo repositories.DetailTaskRepository,
	taskRepo repositories.TaskRepository,
) *DetailTaskService {
	return &DetailTaskService{
		detailTaskRepo: detailTaskRepo,
		taskRepo:       taskRepo,
	}
}

func (s *DetailTaskService) CreateDetailTask(taskID uint64) error {
	_, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	detailTask := models.NewDetailTask(taskID)

	return s.detailTaskRepo.Create(detailTask)
}

func (s *DetailTaskService) Index() ([]models.DetailTask, error) {
	return s.detailTaskRepo.Index()
}

func (s *DetailTaskService) GetByID(id uint64) (*models.DetailTask, error) {
	return s.detailTaskRepo.GetByID(id)
}

func (s *DetailTaskService) GetByTaskID(taskID uint64) ([]models.DetailTask, error) {
	return s.detailTaskRepo.GetByTaskID(taskID)
}

func (s *DetailTaskService) UpdateTaskID(id uint64, taskID uint64) error {
	_, err := s.detailTaskRepo.GetByID(id)
	if err != nil {
		return errors.New("detail task not found")
	}

	_, err = s.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	return s.detailTaskRepo.SetTaskID(id, taskID)
}

func (s *DetailTaskService) DeleteDetailTask(id uint64) error {
	_, err := s.detailTaskRepo.GetByID(id)
	if err != nil {
		return errors.New("detail task not found")
	}

	return s.detailTaskRepo.Delete(id)
}
