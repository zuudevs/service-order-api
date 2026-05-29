/**

 filename  : task_repository_impl.go
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

type taskRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *taskRepository) findTasks(
	query string,
	args ...any,
) ([]models.Task, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Subject,
			&task.Description,
			&task.Status,
			&task.Price,
			&task.Due,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) findTask(
	query string,
	args ...any,
) (*models.Task, error) {
	row := r.db.QueryRow(query, args...)
	var task models.Task
	err := row.Scan(
		&task.ID,
		&task.Subject,
		&task.Description,
		&task.Status,
		&task.Price,
		&task.Due,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// ================================= Basic Repository Features =================================

func NewTaskRepository(
	db *sql.DB,
) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) Create(task *models.Task) error {
	_, err := r.db.Exec(
		`INSERT INTO tasks(subject, description, status, price, due) VALUES (?, ?, ?, ?, ?)`,
		task.Subject,
		task.Description,
		task.Status,
		task.Price,
		task.Due,
	)

	return err
}

func (r *taskRepository) Index() (
	[]models.Task,
	error,
) {
	return r.findTasks(`SELECT id, subject, description, status, price, due FROM tasks`)
}

func (r *taskRepository) Replace(
	id uint64,
	task *models.Task,
) error {
	query := `
		UPDATE tasks SET
			subject = ?,
			description = ?,
			status = ?,
			price = ?,
			due = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		task.Subject,
		task.Description,
		task.Status,
		task.Price,
		task.Due,
		id,
	)

	return err
}

func (r *taskRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)

	return err
}

// ================================= Advance Repository Finder =================================

func (r *taskRepository) GetByID(
	id uint64,
) (*models.Task, error) {
	return r.findTask(`SELECT id, subject, description, status, price, due FROM tasks WHERE id = ?`, id)
}

func (r *taskRepository) GetByStatus(
	status models.TaskStatus,
) ([]models.Task, error) {
	return r.findTasks(`SELECT id, subject, description, status, price, due FROM tasks WHERE status = ?`, status)
}

func (r *taskRepository) GetBySubject(
	subject string,
) ([]models.Task, error) {
	return r.findTasks(`SELECT id, subject, description, status, price, due FROM tasks WHERE subject = ?`, subject)
}

func (r *taskRepository) GetByDueDate(
	due time.Time,
) ([]models.Task, error) {
	return r.findTasks(`SELECT id, subject, description, status, price, due FROM tasks WHERE DATE(due) = DATE(?)`, due)
}

// ================================= Setters =================================

func (r *taskRepository) SetStatus(
	id uint64,
	status models.TaskStatus,
) error {
	_, err := r.db.Exec(`UPDATE tasks SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *taskRepository) SetPrice(
	id uint64,
	price uint64,
) error {
	_, err := r.db.Exec(`UPDATE tasks SET price = ? WHERE id = ?`, price, id)
	return err
}

func (r *taskRepository) SetSubject(
	id uint64,
	subject string,
) error {
	_, err := r.db.Exec(`UPDATE tasks SET subject = ? WHERE id = ?`, subject, id)
	return err
}
