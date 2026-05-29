/**

 filename  : order_repository_impl.go
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

type orderRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *orderRepository) findOrders(
	query string,
	args ...any,
) ([]models.Order, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.OrderDate,
			&order.LastModified,
			&order.TotalPrice,
			&order.PersonID,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) findOrder(
	query string,
	args ...any,
) (*models.Order, error) {
	row := r.db.QueryRow(query, args...)
	var order models.Order
	err := row.Scan(
		&order.ID,
		&order.Status,
		&order.OrderDate,
		&order.LastModified,
		&order.TotalPrice,
		&order.PersonID,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// ================================= Basic Repository Features =================================

func NewOrderRepository(
	db *sql.DB,
) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(order *models.Order) error {
	_, err := r.db.Exec(
		`INSERT INTO orders(status, order_date, last_modified, total_price, person_id) VALUES (?, ?, ?, ?, ?)`,
		order.Status,
		order.OrderDate,
		order.LastModified,
		order.TotalPrice,
		order.PersonID,
	)

	return err
}

func (r *orderRepository) Index() (
	[]models.Order,
	error,
) {
	return r.findOrders(`SELECT id, status, order_date, last_modified, total_price, person_id FROM orders`)
}

func (r *orderRepository) Replace(
	id uint64,
	order *models.Order,
) error {
	query := `
		UPDATE orders SET
			status = ?,
			order_date = ?,
			last_modified = ?,
			total_price = ?,
			person_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		order.Status,
		order.OrderDate,
		order.LastModified,
		order.TotalPrice,
		order.PersonID,
		id,
	)

	return err
}

func (r *orderRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM orders WHERE id = ?`, id)

	return err
}

// ================================= Advance Repository Finder =================================

func (r *orderRepository) GetByID(
	id uint64,
) (*models.Order, error) {
	return r.findOrder(`SELECT id, status, order_date, last_modified, total_price, person_id FROM orders WHERE id = ?`, id)
}

func (r *orderRepository) GetByStatus(
	status models.OrderStatus,
) ([]models.Order, error) {
	return r.findOrders(`SELECT id, status, order_date, last_modified, total_price, person_id FROM orders WHERE status = ?`, status)
}

func (r *orderRepository) GetByOrderDate(
	date time.Time,
) ([]models.Order, error) {
	return r.findOrders(`SELECT id, status, order_date, last_modified, total_price, person_id FROM orders WHERE DATE(order_date) = DATE(?)`, date)
}

func (r *orderRepository) GetByPersonID(
	id uint64,
) ([]models.Order, error) {
	return r.findOrders(`SELECT id, status, order_date, last_modified, total_price, person_id FROM orders WHERE person_id = ?`, id)
}

// ================================= Setters =================================

func (r *orderRepository) SetStatus(
	id uint64,
	status models.OrderStatus,
) error {
	_, err := r.db.Exec(`UPDATE orders SET status = ?, last_modified = ? WHERE id = ?`, status, time.Now().UTC(), id)
	return err
}

func (r *orderRepository) SetTotalPrice(
	id uint64,
	price uint64,
) error {
	_, err := r.db.Exec(`UPDATE orders SET total_price = ?, last_modified = ? WHERE id = ?`, price, time.Now().UTC(), id)
	return err
}

func (r *orderRepository) SetPersonID(
	id uint64,
	personID *uint64,
) error {
	_, err := r.db.Exec(`UPDATE orders SET person_id = ?, last_modified = ? WHERE id = ?`, personID, time.Now().UTC(), id)
	return err
}
