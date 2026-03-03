package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES (:name, :email, :password)
		RETURNING id, created_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&user.ID, &user.CreatedAt)
	}
	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)

	query := `SELECT * FROM users ORDER BY name ASC`

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	fmt.Printf("users: %v\n", users)

	return users, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, name, email, password, created_at
		FROM users
 		WHERE id = $1
  `

	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
