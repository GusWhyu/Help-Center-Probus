package user

import (
	"context"
	"errors"

	"github.com/rafialfauzan/layrd/pkg/database"

	"gorm.io/gorm"
)

var ErrUserAlreadyExists = errors.New("username or email already registered")

type Repository interface {
	database.BaseRepository[User]
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type repository struct {
	database.BaseRepository[User]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		BaseRepository: database.NewBaseRepository[User](db),
		db:             db,
	}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrUserAlreadyExists
	}
	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
