/**

 filename  : person_repository_impl.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"database/sql"
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
)

type personRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *personRepository) findPersons(
	query string,
	args ...any,
) ([]models.Person, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var persons []models.Person

	for rows.Next() {
		var person models.Person
		err := rows.Scan(
			query,
			person.FirstName,
			person.MiddleName,
			person.LastName,
		)

		if err != nil {
			return nil, err
		}

		persons = append(persons, person)
	}

	return persons, nil
}

func (r *personRepository) findPerson(
	query string,
	args ...any,
) (*models.Person, error) {
	row := r.db.QueryRow(query, args...)
	var person models.Person
	err := row.Scan(
		query,
		person.FirstName,
		person.MiddleName,
		person.LastName,
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

// ================================= Basic Repository Features =================================

func NewPersonRepository(
	db *sql.DB,
) PersonRepository {
	return &personRepository{
		db: db,
	}
}

func (r *personRepository) Create(
	person *models.Person,
) error {
	_, err := r.db.Exec(
		`INSERT INTO persons(first_name, middle_name, last_name) VALUES (?, ?, ?)`,
		person.FirstName,
		person.MiddleName,
		person.LastName,
	)

	return err
}

func (r *personRepository) Index() (
	[]models.Person,
	error,
) {
	return r.findPersons(`SELECT id, first_name, middle_name, last_name, created_at FROM persons`)
}

func (r *personRepository) Replace(
	id uint64,
	person *models.Person,
) error {
	query := `
		UPDATE persons SET
			first_name = ?,
			middle_name = ?,
			last_name = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		person.FirstName,
		person.MiddleName,
		person.LastName,
		id,
	)

	return err
}

func (r *personRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM persons WHERE id = ?`, id)

	return err
}

func (r *personRepository) GetByID(
	id uint64,
) (*models.Person, error) {
	return r.findPerson(`SELECT id, first_name, middle_name, last_name, created_at FROM persons WHERE id = ?`, id)
}

func (r *personRepository) GetByFirstName(
	name string,
) ([]models.Person, error) {
	return r.findPersons(`SELECT id, first_name, middle_name, last_name, created_at FROM persons WHERE first_name = ?`, name)
}

func (r *personRepository) GetByMiddleName(
	name string,
) ([]models.Person, error) {
	return r.findPersons(`SELECT id, first_name, middle_name, last_name, created_at FROM persons WHERE middle_name = ?`, name)
}

func (r *personRepository) GetByLastName(
	name string,
) ([]models.Person, error) {
	return r.findPersons(`SELECT id, first_name, middle_name, last_name, created_at FROM persons WHERE last_name = ?`, name)
}

func (r *personRepository) GetByDateCreated(
	created_at time.Time,
) ([]models.Person, error) {
	return r.findPersons(`SELECT id, first_name, middle_name, last_name, created_at FROM persons WHERE created_at = ?`, created_at)
}

func (r *personRepository) SetFirstName(
	id uint64,
	val string,
) error {
	_, err := r.db.Exec(`UPDATE persons SET first_name = ? WHERE id = ?`, val, id)
	return err
}

func (r *personRepository) SetMiddleName(
	id uint64,
	val string,
) error {
	_, err := r.db.Exec(`UPDATE persons SET middle_name = ? WHERE id = ?`, val, id)
	return err
}

func (r *personRepository) SetLastName(
	id uint64,
	val string,
) error {
	_, err := r.db.Exec(`UPDATE persons SET last_name = ? WHERE id = ?`, val, id)
	return err
}
