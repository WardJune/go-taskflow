package service

import (
	"context"
	"errors"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/WardJune/taskflow/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewUserService(userRepo domain.UserRepository, jwtSecret string) domain.UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	t, err := token.Generate(user.ID, user.Email, s.jwtSecret)

	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: t, User: *user}, nil
}

func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {

	user, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	t, err := token.Generate(user.ID, user.Email, s.jwtSecret)

	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: t, User: *user}, nil
}

func (s *userService) FindAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}
