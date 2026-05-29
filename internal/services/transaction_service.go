/**

 filename  : transaction_service.go
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

	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/repositories"
)

type TransactionService struct {
	transactionRepo repositories.TransactionRepository
	orderRepo       repositories.OrderRepository
	backupSvc       *BackupService
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	orderRepo repositories.OrderRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		orderRepo:       orderRepo,
	}
}

func NewTransactionServiceWithBackup(
	transactionRepo repositories.TransactionRepository,
	orderRepo repositories.OrderRepository,
	backupSvc *BackupService,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		orderRepo:       orderRepo,
		backupSvc:       backupSvc,
	}
}

func (s *TransactionService) CreateTransaction(
	status models.TransactionStatus,
	method models.TransactionMethod,
	amount uint64,
	evidencePath string,
	orderID *uint64,
) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}

	evidencePath = strings.TrimSpace(evidencePath)

	if orderID != nil {
		_, err := s.orderRepo.GetByID(*orderID)
		if err != nil {
			return errors.New("order not found")
		}
	}

	transaction := models.NewTransaction(status, method, amount, evidencePath)
	transaction.OrderID = orderID

	err := s.transactionRepo.Create(transaction)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *TransactionService) Index() ([]models.Transaction, error) {
	return s.transactionRepo.Index()
}

func (s *TransactionService) GetByID(id uint64) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

func (s *TransactionService) GetByStatus(status models.TransactionStatus) ([]models.Transaction, error) {
	return s.transactionRepo.GetByStatus(status)
}

func (s *TransactionService) GetByMethod(method models.TransactionMethod) ([]models.Transaction, error) {
	return s.transactionRepo.GetByMethod(method)
}

func (s *TransactionService) GetByOrderID(orderID uint64) ([]models.Transaction, error) {
	return s.transactionRepo.GetByOrderID(orderID)
}

func (s *TransactionService) UpdateStatus(id uint64, status models.TransactionStatus) error {
	_, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	err = s.transactionRepo.SetStatus(id, status)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *TransactionService) UpdateAmount(id uint64, amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}

	_, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	err = s.transactionRepo.SetAmount(id, amount)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}

func (s *TransactionService) DeleteTransaction(id uint64) error {
	_, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	err = s.transactionRepo.Delete(id)
	if err == nil && s.backupSvc != nil {
		s.backupSvc.Increment()
	}
	return err
}
