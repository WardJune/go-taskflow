package domain

import (
	"context"
	"time"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          int64      `db:"id" json:"id"`
	ProjectID   int64      `db:"project_id" json:"project_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status" json:"status"`
	AssigneeID  *int64     `db:"assignee_id" json:"assignee_id"`
	DueDate     *time.Time `db:"due_date" json:"due_date"`
	CreatedBy   int64      `db:"created_by" json:"created_by"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=2"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssigneeID  *int64     `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
	AssigneeID  *int64      `json:"assignee_id"`
	DueDate     *time.Time  `json:"due_date"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id int64) (*Task, error)
	FindByProjectID(ctx context.Context, projectID int64) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}

type TaskService interface {
	Create(ctx context.Context, projectID, creatorID int64, req *CreateTaskRequest) (*Task, error)
	GetByProject(ctx context.Context, projectID, requesterID int64) ([]Task, error)
	Update(ctx context.Context, taskID, requesterID int64, req *UpdateTaskRequest) (*Task, error)
	Delete(ctx context.Context, taskID, requesterID int64) error
}
