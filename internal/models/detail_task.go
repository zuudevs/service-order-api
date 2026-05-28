/**

 filename  : detail_task.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package models

type DetailTask struct {
	ID     uint64  `json:"id"      db:"id"`
	TaskID uint64  `json:"task_id" db:"task_id"`
}

func NewDetailTask(task_id uint64) *DetailTask {
	return &DetailTask{
		TaskID: task_id,
	}
}