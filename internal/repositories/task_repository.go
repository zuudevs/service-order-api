/**

 filename  : task_repository.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
)

type TaskRepository interface {
	Create(task *models.Task) error
	Index() ([]models.Task, error)
	Replace(id uint64, task *models.Task) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Task, error)
	GetByStatus(status models.TaskStatus) ([]models.Task, error)
	GetBySubject(subject string) ([]models.Task, error)
	GetByDueDate(due time.Time) ([]models.Task, error)
	SetStatus(id uint64, status models.TaskStatus) error
	SetPrice(id uint64, price uint64) error
	SetSubject(id uint64, subject string) error
}
