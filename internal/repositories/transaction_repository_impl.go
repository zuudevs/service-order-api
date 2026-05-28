/**

 filename  : transaction_repository_impl.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"database/sql"
	"time"

	"github.com/zuudevs/service-order-api/internal/models"
)

type transactionRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *transactionRepository) findTransactions(
	query string,
	args ...any,
) ([]models.Transaction, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Timestamp,
			&transaction.Status,
			&transaction.Method,
			&transaction.Amount,
			&transaction.EvidencePath,
			&transaction.OrderID,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) findTransaction(
	query string,
	args ...any,
) (*models.Transaction, error) {
	row := r.db.QueryRow(query, args...)
	var transaction models.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.Timestamp,
		&transaction.Status,
		&transaction.Method,
		&transaction.Amount,
		&transaction.EvidencePath,
		&transaction.OrderID,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// ================================= Basic Repository Features =================================

func NewTransactionRepository(
	db *sql.DB,
) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	_, err := r.db.Exec(
		`INSERT INTO transactions(timestamp, status, method, amount, evidence_path, order_id) VALUES (?, ?, ?, ?, ?, ?)`,
		transaction.Timestamp,
		transaction.Status,
		transaction.Method,
		transaction.Amount,
		transaction.EvidencePath,
		transaction.OrderID,
	)

	return err
}

func (r *transactionRepository) Index() (
	[]models.Transaction,
	error,
) {
	return r.findTransactions(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions`)
}

func (r *transactionRepository) Replace(
	id uint64,
	transaction *models.Transaction,
) error {
	query := `
		UPDATE transactions SET
			timestamp = ?,
			status = ?,
			method = ?,
			amount = ?,
			evidence_path = ?,
			order_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		transaction.Timestamp,
		transaction.Status,
		transaction.Method,
		transaction.Amount,
		transaction.EvidencePath,
		transaction.OrderID,
		id,
	)

	return err
}

func (r *transactionRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM transactions WHERE id = ?`, id)

	return err
}

// ================================= Advance Repository Finder =================================

func (r *transactionRepository) GetByID(
	id uint64,
) (*models.Transaction, error) {
	return r.findTransaction(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions WHERE id = ?`, id)
}

func (r *transactionRepository) GetByStatus(
	status models.TransactionStatus,
) ([]models.Transaction, error) {
	return r.findTransactions(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions WHERE status = ?`, status)
}

func (r *transactionRepository) GetByMethod(
	method models.TransactionMethod,
) ([]models.Transaction, error) {
	return r.findTransactions(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions WHERE method = ?`, method)
}

func (r *transactionRepository) GetByOrderID(
	orderID uint64,
) ([]models.Transaction, error) {
	return r.findTransactions(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions WHERE order_id = ?`, orderID)
}

func (r *transactionRepository) GetByTimestamp(
	timestamp time.Time,
) ([]models.Transaction, error) {
	return r.findTransactions(`SELECT id, timestamp, status, method, amount, evidence_path, order_id FROM transactions WHERE DATE(timestamp) = DATE(?)`, timestamp)
}

// ================================= Setters =================================

func (r *transactionRepository) SetStatus(
	id uint64,
	status models.TransactionStatus,
) error {
	_, err := r.db.Exec(`UPDATE transactions SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *transactionRepository) SetAmount(
	id uint64,
	amount uint64,
) error {
	_, err := r.db.Exec(`UPDATE transactions SET amount = ? WHERE id = ?`, amount, id)
	return err
}

func (r *transactionRepository) SetOrderID(
	id uint64,
	orderID *uint64,
) error {
	_, err := r.db.Exec(`UPDATE transactions SET order_id = ? WHERE id = ?`, orderID, id)
	return err
}
