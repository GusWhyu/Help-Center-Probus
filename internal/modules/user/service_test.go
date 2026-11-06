package user

import (
	"context"
	"errors"
	"testing"

	"github.com/rafialfauzan/layrd/pkg/database"
	"gorm.io/gorm"
)

type testUserRepository struct {
	findByEmailUser *User
	findByEmailErr  error
	createdUser     *User
}

func (r *testUserRepository) FindByID(context.Context, any) (*User, error) {
	return nil, database.ErrRecordNotFound
}

func (r *testUserRepository) FindAll(context.Context, int, int, map[string]any, string) ([]User, int64, error) {
	return nil, 0, nil
}

func (r *testUserRepository) Create(_ context.Context, user *User) error {
	r.createdUser = user
	return nil
}

func (r *testUserRepository) Update(context.Context, *User) error { return nil }
func (r *testUserRepository) Delete(context.Context, any) error   { return nil }
func (r *testUserRepository) WithTx(*gorm.DB) database.BaseRepository[User] {
	return nil
}
func (r *testUserRepository) DB() *gorm.DB { return nil }
func (r *testUserRepository) FindByEmail(context.Context, string) (*User, error) {
	return r.findByEmailUser, r.findByEmailErr
}
func (r *testUserRepository) FindByUsername(context.Context, string) (*User, error) {
	return nil, database.ErrRecordNotFound
}

func TestCreateUserDoesNotIgnoreEmailLookupErrors(t *testing.T) {
	lookupErr := errors.New("database unavailable")
	repo := &testUserRepository{findByEmailErr: lookupErr}

	_, err := NewService(repo).CreateUser(context.Background(), &CreateUserRequest{
		Username: "tester",
		Email:    "tester@example.com",
		Password: "long-password",
	})
	if !errors.Is(err, lookupErr) {
		t.Fatalf("CreateUser() error = %v, want %v", err, lookupErr)
	}
	if repo.createdUser != nil {
		t.Fatal("CreateUser() persisted a user after an email lookup failure")
	}
}

func TestCreateUserAssignsUnprivilegedRole(t *testing.T) {
	repo := &testUserRepository{findByEmailErr: database.ErrRecordNotFound}

	if _, err := NewService(repo).CreateUser(context.Background(), &CreateUserRequest{
		Username: "tester",
		Email:    "tester@example.com",
		Password: "long-password",
	}); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if repo.createdUser == nil || repo.createdUser.Role != "user" {
		t.Fatalf("created role = %v, want user", repo.createdUser)
	}
	if repo.createdUser.Password == "long-password" {
		t.Fatal("CreateUser() stored the plaintext password")
	}
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	repo := &testUserRepository{findByEmailUser: &User{Email: "tester@example.com"}}

	_, err := NewService(repo).CreateUser(context.Background(), &CreateUserRequest{
		Username: "tester",
		Email:    "tester@example.com",
		Password: "long-password",
	})
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("CreateUser() error = %v, want %v", err, ErrUserAlreadyExists)
	}
}
