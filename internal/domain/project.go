package domain

import (
	"context"
	"time"
)

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleMember MemberRole = "member"
)

func (r MemberRole) IsValid() bool {
	switch r {
	case MemberRoleOwner, MemberRoleMember:
		return true
	}

	return false
}

type Project struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	OwnerID     int64     `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type ProjectMember struct {
	ProjectID int64      `db:"project_id" json:"project_id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	Role      MemberRole `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type ProjectDetail struct {
	Project
	Members []User `json:"members"`
	Tasks   []Task `json:"tasks"`
}

// Request Struct
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description"`
}

type AddMemberRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// Interface
type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id int64) (*Project, error)
	FindByUserID(ctx context.Context, userID int64) ([]Project, error)
	AddMember(ctx context.Context, member *ProjectMember) error
	IsMember(ctx context.Context, projectID, userID int64) (bool, error)
	GetMembers(ctx context.Context, projectID int64) ([]User, error)
}

type ProjectService interface {
	Create(ctx context.Context, ownerID int64, req *CreateProjectRequest) (*Project, error)
	GetByID(ctx context.Context, projectID, requesterID int64) (*ProjectDetail, error)
	GetUserProjects(ctx context.Context, userID int64) ([]Project, error)
	AddMember(ctx context.Context, projectID, ownerID int64, req *AddMemberRequest) error
}
