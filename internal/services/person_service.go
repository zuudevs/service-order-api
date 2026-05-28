/**

 filename  : person_service.go
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

type PersonService struct {
	personRepo repositories.PersonRepository
}

func NewPersonService(
	personRepo repositories.PersonRepository,
) *PersonService {
	return &PersonService{
		personRepo: personRepo,
	}
}

func (s *PersonService) Create(
	firstname string,
	middlename *string,
	lastname *string,
) error {
	firstname = strings.TrimSpace(firstname)

	if firstname == "" {
		return errors.New("firstname required")
	}

	person := models.NewPerson(
		firstname,
		middlename,
		lastname,
	)

	return s.personRepo.Create(person)
}

func (s *PersonService) Index() (
	[]models.Person,
	error,
) {
	return s.personRepo.Index()
}

func (s *PersonService) GetByID(id uint64) (
	*models.Person,
	error,
) {
	return s.personRepo.GetByID(id)
}

func (s *PersonService) GetByFirstName(name string) (
	[]models.Person,
	error,
) {
	return s.personRepo.GetByFirstName(name)
}

func (s *PersonService) GetByLastName(name string) (
	[]models.Person,
	error,
) {
	return s.personRepo.GetByLastName(name)
}

func (s *PersonService) UpdatePerson(
	id uint64,
	firstname string,
	middlename *string,
	lastname *string,
) error {
	firstname = strings.TrimSpace(firstname)

	if firstname == "" {
		return errors.New("firstname required")
	}

	_, err := s.personRepo.GetByID(id)
	if err != nil {
		return errors.New("person not found")
	}

	person := models.NewPerson(firstname, middlename, lastname)
	return s.personRepo.Replace(id, person)
}

func (s *PersonService) UpdateFirstName(id uint64, firstname string) error {
	firstname = strings.TrimSpace(firstname)

	if firstname == "" {
		return errors.New("firstname required")
	}

	_, err := s.personRepo.GetByID(id)
	if err != nil {
		return errors.New("person not found")
	}

	return s.personRepo.SetFirstName(id, firstname)
}

func (s *PersonService) UpdateMiddleName(id uint64, middlename string) error {
	_, err := s.personRepo.GetByID(id)
	if err != nil {
		return errors.New("person not found")
	}

	return s.personRepo.SetMiddleName(id, middlename)
}

func (s *PersonService) UpdateLastName(id uint64, lastname string) error {
	_, err := s.personRepo.GetByID(id)
	if err != nil {
		return errors.New("person not found")
	}

	return s.personRepo.SetLastName(id, lastname)
}

func (s *PersonService) DeletePerson(id uint64) error {
	_, err := s.personRepo.GetByID(id)
	if err != nil {
		return errors.New("person not found")
	}

	return s.personRepo.Delete(id)
}
