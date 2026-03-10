package service

import (
	"context"
	"errors"

	"github.com/WardJune/taskflow/internal/domain"
	ws "github.com/WardJune/taskflow/internal/websocket"
)

type taskService struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
	hub         *ws.Hub
}

func NewTaskService(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository, hub *ws.Hub) domain.TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		hub:         hub,
	}
}

func (s *taskService) Create(ctx context.Context, projectID, creatorID int64, req *domain.CreateTaskRequest) (*domain.Task, error) {
	// is member
	isMember, err := s.projectRepo.IsMember(ctx, projectID, creatorID)

	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("access denied")
	}

	// is asignee not nil
	if req.AssigneeID != nil {
		isAssigneeMember, err := s.projectRepo.IsMember(ctx, projectID, *req.AssigneeID)
		if err != nil {
			return nil, err
		}

		if !isAssigneeMember {
			return nil, errors.New("assignee is not a member of the project")
		}
	}

	task := &domain.Task{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
		CreatedBy:   creatorID,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	s.hub.BroadcastToProject(projectID, ws.Notification{
		Type:      "task_created",
		ProjectID: projectID,
		Message:   "New task created: " + task.Title,
		Data:      task,
	})

	return task, nil
}

func (s *taskService) GetByProject(ctx context.Context, projectID, requesterID int64) ([]domain.Task, error) {
	isMember, err := s.projectRepo.IsMember(ctx, projectID, requesterID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("access denied")
	}

	return s.taskRepo.FindByProjectID(ctx, projectID)
}

func (s *taskService) Update(ctx context.Context, taskID, requesterID int64, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	isMember, err := s.projectRepo.IsMember(ctx, task.ProjectID, requesterID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("access denied")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	s.hub.BroadcastToProject(task.ProjectID, ws.Notification{
		Type:      "task_updated",
		ProjectID: task.ProjectID,
		Message:   "Task updated: " + task.Title,
		Data:      task,
	})

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, taskID, requesterID int64) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)

	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	if task.CreatedBy != requesterID {
		return errors.New("only task creator can delete this task")
	}

	s.hub.BroadcastToProject(task.ProjectID, ws.Notification{
		Type:      "task_deleted",
		ProjectID: task.ProjectID,
		Message:   "Task deleted: " + task.Title,
		Data:      map[string]int64{"task_id": taskID},
	})

	return s.taskRepo.Delete(ctx, taskID)
}
