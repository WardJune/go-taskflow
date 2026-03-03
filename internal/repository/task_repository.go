package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (project_id, title, description, assignee_id, due_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, status, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query, task.ProjectID, task.Title, task.Description, task.AssigneeID, task.DueDate, task.CreatedBy).Scan(&task.ID, &task.Status, &task.CreatedAt, &task.UpdatedAt)
}

func (r *taskRepository) FindByID(ctx context.Context, id int64) (*domain.Task, error) {
	var task domain.Task

	query := `SELECT * FROM tasks WHERE id = $1`

	err := r.db.GetContext(ctx, &task, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) FindByProjectID(ctx context.Context, projectID int64) ([]domain.Task, error) {
	tasks := make([]domain.Task, 0)

	query := `
		SELECT * FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(ctx, &tasks, query, projectID); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks SET
			title = $1,
			description = $2,
			status = $3,
			assignee_id = $4,
			due_date = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query, task.Title, task.Description, task.Status, task.AssigneeID, task.DueDate, task.ID).Scan(&task.UpdatedAt)
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
