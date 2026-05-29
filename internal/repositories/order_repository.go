/**

 filename  : order_repository.go
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

type OrderRepository interface {
	Create(order *models.Order) error
	Index() ([]models.Order, error)
	Replace(id uint64, order *models.Order) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Order, error)
	GetByStatus(status models.OrderStatus) ([]models.Order, error)
	GetByOrderDate(date time.Time) ([]models.Order, error)
	GetByPersonID(id uint64) ([]models.Order, error)
	SetStatus(id uint64, status models.OrderStatus) error
	SetTotalPrice(id uint64, price uint64) error
	SetPersonID(id uint64, personID *uint64) error
}
