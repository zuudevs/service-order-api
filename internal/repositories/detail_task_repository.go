/**

 filename  : detail_task_repository.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"github.com/zuudevs/service-order-api/internal/models"
)

type DetailTaskRepository interface {
	Create(detailTask *models.DetailTask) error
	Index() ([]models.DetailTask, error)
	Replace(id uint64, detailTask *models.DetailTask) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.DetailTask, error)
	GetByTaskID(taskID uint64) ([]models.DetailTask, error)
	SetTaskID(id uint64, taskID uint64) error
}
