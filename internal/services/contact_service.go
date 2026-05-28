/**

 filename  : contact_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/repositories"
)

type ContactService struct {
	contactRepo repositories.ContactRepository
	personRepo  repositories.PersonRepository
}

func NewContactService(
	contactRepo repositories.ContactRepository,
	personRepo repositories.PersonRepository,
) *ContactService {
	return &ContactService{
		contactRepo: contactRepo,
		personRepo:  personRepo,
	}
}

func (s *ContactService) CreateContact(
	value string,
	contactType models.ContactType,
	isMain bool,
	personID uint64,
) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return errors.New("value required")
	}

	switch contactType {
	case models.ContactTypePhone:
		if len(value) < 10 ||
			len(value) > 15 {
			return errors.New("phone number is not valid")
		}
	case models.ContactTypeEmail:
		matched, err := regexp.Match(
			"^[A-Za-z][A-Za-z0-9_]*@[A-Za-z0-9-]+\\.[A-Za-z]{2,}$",
			[]byte(value),
		)

		if err != nil {
			return errors.New(err.Error())
		}

		if !matched {
			return errors.New("email is not valid")
		}
	}

	_, err := s.personRepo.GetByID(personID)

	if err != nil {
		return errors.New("person not found")
	}

	contact := models.NewContact(
		value,
		contactType,
		isMain,
		personID,
	)

	return s.contactRepo.Create(contact)
}

func (s *ContactService) Index() (
	[]models.Contact,
	error,
) {
	return s.contactRepo.Index()
}

func (s *ContactService) GetByID(id uint64) (
	*models.Contact,
	error,
) {
	return s.contactRepo.GetByID(id)
}

func (s *ContactService) GetByValue(value string) (
	[]models.Contact,
	error,
) {
	return s.contactRepo.GetByValue(value)
}

func (s *ContactService) GetByContactType(contactType models.ContactType) (
	[]models.Contact,
	error,
) {
	return s.contactRepo.GetByContactType(contactType)
}

func (s *ContactService) GetByPersonID(personID uint64) (
	[]models.Contact,
	error,
) {
	return s.contactRepo.GetByPersonID(personID)
}

func (s *ContactService) UpdateContact(
	id uint64,
	value string,
	contactType models.ContactType,
	isMain bool,
	personID uint64,
) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return errors.New("value required")
	}

	switch contactType {
	case models.ContactTypePhone:
		if len(value) < 10 || len(value) > 15 {
			return errors.New("phone number is not valid")
		}
	case models.ContactTypeEmail:
		matched, err := regexp.Match(
			"^[A-Za-z][A-Za-z0-9_]*@[A-Za-z0-9-]+\\.[A-Za-z]{2,}$",
			[]byte(value),
		)
		if err != nil {
			return errors.New(err.Error())
		}
		if !matched {
			return errors.New("email is not valid")
		}
	}

	_, err := s.personRepo.GetByID(personID)
	if err != nil {
		return errors.New("person not found")
	}

	_, err = s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("contact not found")
	}

	contact := models.NewContact(value, contactType, isMain, personID)
	return s.contactRepo.Replace(id, contact)
}

func (s *ContactService) UpdateValue(id uint64, value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return errors.New("value required")
	}

	_, err := s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("contact not found")
	}

	return s.contactRepo.SetValue(id, value)
}

func (s *ContactService) UpdateContactType(id uint64, contactType models.ContactType) error {
	_, err := s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("contact not found")
	}

	return s.contactRepo.SetContactType(id, contactType)
}

func (s *ContactService) DeleteContact(id uint64) error {
	_, err := s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("contact not found")
	}

	return s.contactRepo.Delete(id)
}
