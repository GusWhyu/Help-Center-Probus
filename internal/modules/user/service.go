package user

import (
	"context"
	"errors"

	"github.com/rafialfauzan/layrd/pkg/database"
	"github.com/rafialfauzan/layrd/pkg/utils"
)

type Service interface {
	GetUsers(ctx context.Context, page, limit int, role string) ([]User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUsers(ctx context.Context, page, limit int, role string) ([]User, int64, error) {
	conditions := make(map[string]any)
	if role != "" {
		conditions["role = ?"] = role
	}
	return s.repo.FindAll(ctx, page, limit, conditions, "id desc")
}

func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	switch {
	case err == nil:
		return nil, ErrUserAlreadyExists
	case !errors.Is(err, database.ErrRecordNotFound):
		return nil, err
	}

	_, err = s.repo.FindByUsername(ctx, req.Username)
	switch {
	case err == nil:
		return nil, ErrUserAlreadyExists
	case !errors.Is(err, database.ErrRecordNotFound):
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
