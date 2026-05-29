/**

 filename  : contact_repository.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
)

type ContactRepository interface {
	Create(contact *models.Contact) error
	Index() ([]models.Contact, error)
	Replace(id uint64, contact *models.Contact) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Contact, error)
	GetByValue(value string) ([]models.Contact, error)
	GetByContactType(contactType models.ContactType) ([]models.Contact, error)
	GetByDateCreated(created_at time.Time) ([]models.Contact, error)
	GetByPersonID(id uint64) ([]models.Contact, error)
	SetValue(id uint64, val string) error
	SetContactType(id uint64, val models.ContactType) error
}
