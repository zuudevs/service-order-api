/**

 filename  : transaction.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

import "time"

type TransactionStatus uint8
type TransactionMethod uint8

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusCanceled
	TransactionStatusSuccess
)

const (
	TransactionMethodCash TransactionMethod = iota
	TransactionMethodEWallet
	TransactionMethodBank
)

type Transaction struct {
	ID            uint64            `json:"id"            db:"id"`
	Timestamp     time.Time         `json:"timestamp"     db:"timestamp"`
	Status        TransactionStatus `json:"status"        db:"status"`
	Method        TransactionMethod `json:"method"        db:"method"`
	Amount        uint64            `json:"amount"        db:"amount"`
	EvidencePath  string            `json:"evidance_path" db:"evidence_path"`
	OrderID       *uint64           `json:"order_id"      db:"order_id"`
}

func NewTransaction(
	status TransactionStatus,
	method TransactionMethod,
	amount uint64,
	evidence_path string,
) *Transaction {
	return &Transaction{
		Timestamp:     time.Now().UTC(),
		Status:        status,
		Method:        method,
		Amount:        amount,
		EvidencePath:  evidence_path,
		OrderID:       nil,
	}
}
