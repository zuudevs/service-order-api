/**

 filename  : contact.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

import "time"

type ContactType uint8;

const (
	ContactTypeEmail = iota
	ContactTypePhone
)

type Contact struct {
	ID          uint64      `json:"id"           db:"id"`
	Value       string      `json:"value"        db:"value"`
	ContactType ContactType `json:"contact_type" db:"contact_type"`
	IsMain      bool        `json:"is_main"      db:"is_main"`
	PersonID    uint64      `json:"person_id"    db:"person_id"`
	CreatedAt   time.Time   `json:"created_at"   db:"created_at"`
}

func NewContact(
	value string,
	contactType ContactType,
	isMain bool,
	personID uint64,
) *Contact {
	return &Contact{
		Value:       value,
		ContactType: contactType,
		CreatedAt:   time.Now().UTC(),
		IsMain:      isMain,
		PersonID:    personID,
	}
}