/**

 filename  : order.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

import "time"

type OrderStatus uint8

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusAccepted
	OrderStatusRejected
	OrderStatusRevised
	OrderStatusCompleted
)

type Order struct {
	ID           uint64      `json:"id"            db:"id"`
	Status       OrderStatus `json:"status"        db:"status"`
	OrderDate    time.Time   `json:"order_date"    db:"order_date"`
	LastModified time.Time   `json:"last_modified" db:"last_modified"`
	TotalPrice   uint64      `json:"total_price"   db:"total_price"`
	PersonID     *uint64     `json:"person_id"     db:"person_id"`
}

func NewOrder(
	status OrderStatus,
	dt time.Time,
) *Order {
	return &Order{
		Status:       OrderStatusPending,
		OrderDate:    dt,
		LastModified: dt,
		TotalPrice:   0,
		PersonID:     nil,
	}
}
