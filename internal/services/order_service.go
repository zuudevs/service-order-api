/**

 filename  : order_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package services

import (
	"errors"
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/repositories"
)

type OrderService struct {
	orderRepo  repositories.OrderRepository
	personRepo repositories.PersonRepository
	backupSvc  *BackupService
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	personRepo repositories.PersonRepository,
) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		personRepo: personRepo,
	}
}

func NewOrderServiceWithBackup(
	orderRepo repositories.OrderRepository,
	personRepo repositories.PersonRepository,
	backupSvc *BackupService,
) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		personRepo: personRepo,
		backupSvc:  backupSvc,
	}
}

func (s *OrderService) CreateOrder(
	status models.OrderStatus,
	personID *uint64,
) error {
	if personID != nil {
		_, err := s.personRepo.GetByID(*personID)
		if err != nil {
			return errors.New("person not found")
		}
	}

	order := models.NewOrder(status, time.Now().UTC())
	order.PersonID = personID

	err := s.orderRepo.Create(order)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *OrderService) Index() ([]models.Order, error) {
	return s.orderRepo.Index()
}

func (s *OrderService) GetByID(id uint64) (*models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetByStatus(status models.OrderStatus) ([]models.Order, error) {
	return s.orderRepo.GetByStatus(status)
}

func (s *OrderService) GetByPersonID(personID uint64) ([]models.Order, error) {
	return s.orderRepo.GetByPersonID(personID)
}

func (s *OrderService) UpdateStatus(id uint64, status models.OrderStatus) error {
	_, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	err = s.orderRepo.SetStatus(id, status)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *OrderService) UpdateTotalPrice(id uint64, price uint64) error {
	_, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	err = s.orderRepo.SetTotalPrice(id, price)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *OrderService) DeleteOrder(id uint64) error {
	_, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	err = s.orderRepo.Delete(id)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}
