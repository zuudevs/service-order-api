/**

 filename  : person_repository.go
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

type PersonRepository interface {
	Create(person *models.Person) error
	Index() ([]models.Person, error)
	Replace(id uint64, person *models.Person) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Person, error)
	GetByFirstName(name string) ([]models.Person, error)
	GetByMiddleName(name string) ([]models.Person, error)
	GetByLastName(name string) ([]models.Person, error)
	GetByDateCreated(created_at time.Time) ([]models.Person, error)
	SetFirstName(id uint64, val string) error
	SetMiddleName(id uint64, val string) error
	SetLastName(id uint64, val string) error
}
