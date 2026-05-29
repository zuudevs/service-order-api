/**

 filename  : transaction_repository.go
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

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	Index() ([]models.Transaction, error)
	Replace(id uint64, transaction *models.Transaction) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Transaction, error)
	GetByStatus(status models.TransactionStatus) ([]models.Transaction, error)
	GetByMethod(method models.TransactionMethod) ([]models.Transaction, error)
	GetByOrderID(orderID uint64) ([]models.Transaction, error)
	GetByTimestamp(timestamp time.Time) ([]models.Transaction, error)
	SetStatus(id uint64, status models.TransactionStatus) error
	SetAmount(id uint64, amount uint64) error
	SetOrderID(id uint64, orderID *uint64) error
}
