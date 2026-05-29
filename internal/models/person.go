/**

 filename  : person.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

import "time"

type Person struct {
	ID         uint64    `json:"id"         db:"id"`
	FirstName  string    `json:"firstname"  db:"first_name"`
	MiddleName *string   `json:"middlename" db:"middle_name"`
	LastName   *string   `json:"lastname"   db:"last_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func NewPerson(
	firstname string,
	middlename *string,
	lastname *string,
) *Person {
	return &Person{
		FirstName:  firstname,
		MiddleName: middlename,
		LastName:   lastname,
		CreatedAt:  time.Now().UTC(),
	}
}
