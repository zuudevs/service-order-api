/**

 filename  : task.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

import "time"

type TaskStatus uint8

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusOnProgress
	TaskStatusComplete
)

type Task struct {
	ID          uint64     `json:"id"          db:"id"`
	Subject     string     `json:"subject"     db:"subject"`
	Description string     `json:"description" db:"description"`
	Status      TaskStatus `json:"status"      db:"status"`
	Price       uint64     `json:"price"       db:"price"`
	Due         time.Time  `json:"due"         db:"due"`
}

func NewTask(
	subject string,
	description string,
	price uint64,
	due time.Time,
) *Task {
	return &Task{
		Subject:     subject,
		Description: description,
		Status:      TaskStatusPending,
		Price:       price,
		Due:         due,
	}
}
