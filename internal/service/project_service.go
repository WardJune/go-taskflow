package service

import (
	"context"
	"errors"

	"github.com/WardJune/taskflow/internal/domain"
)

type projectService struct {
	projectRepo domain.ProjectRepository
	taskRepo    domain.TaskRepository
	userRepo    domain.UserRepository
}

func NewProjectService(
	projectRepo domain.ProjectRepository,
	taskRepo domain.TaskRepository,
	userRepo domain.UserRepository,
) domain.ProjectService {
	return &projectService{
		projectRepo,
		taskRepo,
		userRepo,
	}
}

func (s *projectService) Create(ctx context.Context, ownerID int64, req *domain.CreateProjectRequest) (*domain.Project, error) {
	project := &domain.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	member := &domain.ProjectMember{
		ProjectID: project.ID,
		UserID:    ownerID,
		Role:      domain.MemberRoleOwner,
	}

	if err := s.projectRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) GetByID(ctx context.Context, projectID, requesterID int64) (*domain.ProjectDetail, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, errors.New("project not found")
	}

	//check if member
	isMember, err := s.projectRepo.IsMember(ctx, projectID, requesterID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("access denied")
	}

	//get members
	members, err := s.projectRepo.GetMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	//get tasks
	tasks, err := s.taskRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &domain.ProjectDetail{
		Project: *project,
		Members: members,
		Tasks:   tasks,
	}, nil

}

func (s *projectService) GetUserProjects(ctx context.Context, userID int64) ([]domain.Project, error) {
	projects, err := s.projectRepo.FindByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *projectService) GetAvailableUserProjects(ctx context.Context, projectID int64) ([]domain.User, error) {
	users, err := s.projectRepo.GetAvailableUser(ctx, projectID)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *projectService) AddMember(ctx context.Context, projectID, ownerID int64, req *domain.AddMemberRequest) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)

	if err != nil {
		return err
	}

	if project == nil {
		return errors.New("project not found")
	}

	// only owner add member
	if project.OwnerID != ownerID {
		return errors.New("only project onwner can add members")
	}

	user, err := s.userRepo.FindById(ctx, req.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	member := &domain.ProjectMember{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      domain.MemberRoleMember,
	}

	return s.projectRepo.AddMember(ctx, member)
}
