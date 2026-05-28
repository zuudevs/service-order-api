/**

 filename  : detail_task_repository_impl.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package repositories

import (
	"database/sql"

	"github.com/zuudevs/service-order-api/internal/models"
)

type detailTaskRepository struct {
	db *sql.DB
}

// ================================= Helper Repository =================================

func (r *detailTaskRepository) findDetailTasks(
	query string,
	args ...any,
) ([]models.DetailTask, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var detailTasks []models.DetailTask

	for rows.Next() {
		var detailTask models.DetailTask
		err := rows.Scan(
			&detailTask.ID,
			&detailTask.TaskID,
		)

		if err != nil {
			return nil, err
		}

		detailTasks = append(detailTasks, detailTask)
	}

	return detailTasks, nil
}

func (r *detailTaskRepository) findDetailTask(
	query string,
	args ...any,
) (*models.DetailTask, error) {
	row := r.db.QueryRow(query, args...)
	var detailTask models.DetailTask
	err := row.Scan(
		&detailTask.ID,
		&detailTask.TaskID,
	)

	if err != nil {
		return nil, err
	}

	return &detailTask, nil
}

// ================================= Basic Repository Features =================================

func NewDetailTaskRepository(
	db *sql.DB,
) DetailTaskRepository {
	return &detailTaskRepository{
		db: db,
	}
}

func (r *detailTaskRepository) Create(detailTask *models.DetailTask) error {
	_, err := r.db.Exec(
		`INSERT INTO detail_tasks(task_id) VALUES (?)`,
		detailTask.TaskID,
	)

	return err
}

func (r *detailTaskRepository) Index() (
	[]models.DetailTask,
	error,
) {
	return r.findDetailTasks(`SELECT id, task_id FROM detail_tasks`)
}

func (r *detailTaskRepository) Replace(
	id uint64,
	detailTask *models.DetailTask,
) error {
	query := `
		UPDATE detail_tasks SET
			task_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		detailTask.TaskID,
		id,
	)

	return err
}

func (r *detailTaskRepository) Delete(
	id uint64,
) error {
	_, err := r.db.Exec(`DELETE FROM detail_tasks WHERE id = ?`, id)

	return err
}

// ================================= Advance Repository Finder =================================

func (r *detailTaskRepository) GetByID(
	id uint64,
) (*models.DetailTask, error) {
	return r.findDetailTask(`SELECT id, task_id FROM detail_tasks WHERE id = ?`, id)
}

func (r *detailTaskRepository) GetByTaskID(
	taskID uint64,
) ([]models.DetailTask, error) {
	return r.findDetailTasks(`SELECT id, task_id FROM detail_tasks WHERE task_id = ?`, taskID)
}

// ================================= Setters =================================

func (r *detailTaskRepository) SetTaskID(
	id uint64,
	taskID uint64,
) error {
	_, err := r.db.Exec(`UPDATE detail_tasks SET task_id = ? WHERE id = ?`, taskID, id)
	return err
}
