package repository

import (
	"database/sql"
	"log"
	"pet-project/pkg/model"
)

type PostgresTaskRepository struct {
	DB *sql.DB
}

type TaskRepository interface {
	CreateTask(task *model.Task) error                      ///
	UpdateTask(task *model.Task) error                      ///
	GetByIDTask(id int) (*model.Task, error)                ///
	ListByProjectTask(projectID int) ([]*model.Task, error) ///
	DeleteTask(id int) error                                ///
}

func (rt *PostgresTaskRepository) CreateTask(task *model.Task) error {
	query := `INSERT INTO tasks (title, description, status, priority, assigned_to, project_id, created_at, updated_at, due_date)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := rt.DB.Exec(query, task.Title, task.Description,
		task.Status, task.Priority, task.AssignedTo, task.ProjectID, task.CreatedAt, task.UpdatedAt, task.DueDate)
	if err != nil {
		log.Println("Failed to create task:", err)
		return err
	}
	return nil
}

func (rt *PostgresTaskRepository) UpdateTask(task *model.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, assigned_to = $5, updated_at = $6, due_date = $7 WHERE id = $8`
	_, err := rt.DB.Exec(query, task.Title, task.Description,
		task.Status, task.Priority, task.AssignedTo, task.UpdatedAt, task.DueDate, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (rt *PostgresTaskRepository) GetByIDTask(id int) (*model.Task, error) {
	task := &model.Task{}
	query := `SELECT id, title, description, status, priority, assigned_to, created_at, updated_at, due_date FROM tasks WHERE id = $1`
	row := rt.DB.QueryRow(query, id)
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.AssignedTo, &task.CreatedAt, &task.UpdatedAt, &task.DueDate)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (rt *PostgresTaskRepository) ListByProjectTask(projectID int) ([]*model.Task, error) {
	query := `SELECT id, title, description, status, priority, assigned_to, created_at, updated_at, due_date FROM tasks WHERE project_id = $1`
	rows, err := rt.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []*model.Task{}

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.AssignedTo, &task.CreatedAt, &task.UpdatedAt, &task.DueDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil

}

func (rt *PostgresTaskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := rt.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
