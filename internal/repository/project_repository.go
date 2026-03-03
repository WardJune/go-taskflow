package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type projectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	query := `INSERT INTO projects (name, description, owner_id)
	VALUES ($1, $2, $3) RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query, project.Name, project.Description, project.OwnerID).Scan(&project.ID, &project.CreatedAt)
}

func (r *projectRepository) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project

	query := `
		SELECT id, name,description, owner_id, created_at FROM projects
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &project, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindByUserID(ctx context.Context, userID int64) ([]domain.Project, error) {
	projects := make([]domain.Project, 0)

	query := `
		SELECT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		INNER JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1
		ORDER BY p.created_at DESC
	`

	if err := r.db.SelectContext(ctx, &projects, query, userID); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepository) AddMember(ctx context.Context, member *domain.ProjectMember) error {

	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, member.ProjectID, member.UserID, member.Role)
	return err
}

func (r *projectRepository) IsMember(ctx context.Context, projectID, userID int64) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM project_members WHERE project_id = $1 AND user_id = $2`

	err := r.db.GetContext(ctx, &count, query, projectID, userID)

	return count > 0, err
}

func (r *projectRepository) GetMembers(ctx context.Context, projectID int64) ([]domain.User, error) {
	members := make([]domain.User, 0)

	query := `
		SELECT u.id, u.name, u.email, u.created_at
		FROM users u
		INNER JOIN project_members pm ON u.id = pm.user_id
		WHERE pm.project_id = $1
	`

	if err := r.db.SelectContext(ctx, &members, query, projectID); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *projectRepository) GetAvailableUser(ctx context.Context, projectID int64) ([]domain.User, error) {
	members := make([]domain.User, 0)

	query := `SELECT * FROM users
		WHERE id NOT IN
			(
				SELECT user_id FROM project_members
				WHERE project_id = $1
			);
	;`

	if err := r.db.SelectContext(ctx, &members, query, projectID); err != nil {
		return nil, err
	}

	return members, nil
}
