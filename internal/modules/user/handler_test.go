package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/utils"
)

type testUserService struct {
	getUsersErr error
}

func (s *testUserService) GetUsers(_ context.Context, _, _ int, _ string) ([]User, int64, error) {
	return []User{}, 0, s.getUsersErr
}
func (*testUserService) GetUserByID(context.Context, uint) (*User, error) {
	return &User{ID: 1}, nil
}
func (*testUserService) CreateUser(context.Context, *CreateUserRequest) (*User, error) {
	return &User{ID: 1, Role: "user"}, nil
}
func (*testUserService) UpdateUser(context.Context, uint, *UpdateUserRequest) (*User, error) {
	return &User{ID: 1}, nil
}
func (*testUserService) DeleteUser(context.Context, uint) error { return nil }

func TestUserRoutesRequireJWTAndValidatePagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &testUserService{}
	tokenService := utils.NewTokenService("test-secret-that-is-long-enough-123", 1)
	module := &Module{
		handler:      NewHandler(service, &config.Config{}),
		tokenService: tokenService,
	}
	router := gin.New()
	module.RegisterRoutes(router.Group("/api/v1"))

	unauthenticated := httptest.NewRecorder()
	router.ServeHTTP(unauthenticated, httptest.NewRequest(http.MethodGet, "/api/v1/users", nil))
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d, want %d", unauthenticated.Code, http.StatusUnauthorized)
	}

	token, err := tokenService.GenerateToken("1", "tester", "user", nil)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	invalidPageSize := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=101", nil)
	invalidPageSize.Header.Set("Authorization", "Bearer "+token)
	invalidPageSizeResponse := httptest.NewRecorder()
	router.ServeHTTP(invalidPageSizeResponse, invalidPageSize)
	if invalidPageSizeResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid page size status = %d, want %d", invalidPageSizeResponse.Code, http.StatusBadRequest)
	}

	validRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	validRequest.Header.Set("Authorization", "Bearer "+token)
	validResponse := httptest.NewRecorder()
	router.ServeHTTP(validResponse, validRequest)
	if validResponse.Code != http.StatusOK {
		t.Fatalf("authenticated status = %d, want %d", validResponse.Code, http.StatusOK)
	}

	service.getUsersErr = errors.New("sensitive database detail")
	failingRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	failingRequest.Header.Set("Authorization", "Bearer "+token)
	failingResponse := httptest.NewRecorder()
	router.ServeHTTP(failingResponse, failingRequest)
	if failingResponse.Code != http.StatusInternalServerError {
		t.Fatalf("failing request status = %d, want %d", failingResponse.Code, http.StatusInternalServerError)
	}
	if strings.Contains(failingResponse.Body.String(), "sensitive database detail") {
		t.Fatal("internal error detail was exposed in the response")
	}
}
