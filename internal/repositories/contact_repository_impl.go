/**

 filename  : contact_repository_impl.go
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

type contactRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *contactRepository) findContacts(
	query string,
	args ...any,
) ([]models.Contact, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
			&contact.Value,
			&contact.ContactType,
			&contact.IsMain,
			&contact.CreatedAt,
			&contact.PersonID,
		)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *contactRepository) findContact(
	query string,
	args ...any,
) (*models.Contact, error) {
	row := r.db.QueryRow(query, args...)
	var contact models.Contact
	err := row.Scan(
		&contact.ID,
		&contact.Value,
		&contact.ContactType,
		&contact.IsMain,
		&contact.CreatedAt,
		&contact.PersonID,
	)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

// ================================= Basic Repository Features =================================

func NewContactRepository(
	db *sql.DB,
) ContactRepository {
	return &contactRepository{
		db: db,
	}
}

func (r *contactRepository) Create(contact *models.Contact) error {
	_, err := r.db.Exec(
		`INSERT INTO contacts(value, contact_type, is_main, person_id) VALUES (?, ?, ?, ?)`,
		contact.Value,
		contact.ContactType,
		contact.IsMain,
		contact.PersonID,
	)

	return err
}

func (r *contactRepository) Index() (
	[]models.Contact,
	error,
) {
	return r.findContacts(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts`)
}

func (r *contactRepository) Replace(
	id uint64,
	contact *models.Contact,
) error {
	query := `
		UPDATE contacts SET
			value = ?, 
			contact_type = ?, 
			is_main = ?,
			person_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		contact.Value,
		contact.ContactType,
		contact.IsMain,
		contact.PersonID,
		id,
	)

	return err
}

func (r *contactRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM contacts WHERE id = ?`, id)

	return err
}

// ================================= Advance Repository Finder =================================

func (r *contactRepository) GetByID(
	id uint64,
) (
	*models.Contact,
	error,
) {
	return r.findContact(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts WHERE id = ?`, id)
}

func (r *contactRepository) GetByValue(
	value string,
) (
	[]models.Contact,
	error,
) {
	return r.findContacts(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts WHERE value = ?`, value)
}

func (r *contactRepository) GetByContactType(
	contactType models.ContactType,
) (
	[]models.Contact,
	error,
) {
	return r.findContacts(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts WHERE contact_type = ?`, contactType)
}

func (r *contactRepository) GetByDateCreated(
	created_at time.Time,
) (
	[]models.Contact,
	error,
) {
	return r.findContacts(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts WHERE created_at = ?`, created_at)
}

func (r *contactRepository) GetByPersonID(
	id uint64,
) (
	[]models.Contact,
	error,
) {
	return r.findContacts(`SELECT id, value, contact_type, is_main, created_at, person_id FROM contacts WHERE person_id = ?`, id)
}

// ================================= Advance Repository Finder =================================

func (r *contactRepository) SetValue(
	id uint64,
	value string,
) error {
	_, err := r.db.Exec(`UPDATE contacts SET value = ? WHERE id = ?`, value, id)
	return err
}

func (r *contactRepository) SetContactType(
	id uint64,
	val models.ContactType,
) error {
	_, err := r.db.Exec(`UPDATE contacts SET contact_type = ? WHERE id = ?`, val, id)
	return err
}
