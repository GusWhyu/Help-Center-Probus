package database

import (
	"context"
	"errors"
	"testing"
)

func TestFindAllRejectsInvalidPagination(t *testing.T) {
	repo := &baseRepository[struct{}]{}
	maxInt := int(^uint(0) >> 1)

	for _, tt := range []struct {
		name  string
		page  int
		limit int
	}{
		{name: "zero page", page: 0, limit: 10},
		{name: "negative limit", page: 1, limit: -1},
		{name: "offset overflow", page: maxInt, limit: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := repo.FindAll(context.Background(), tt.page, tt.limit, nil, "")
			if !errors.Is(err, ErrInvalidPage) {
				t.Fatalf("FindAll() error = %v, want %v", err, ErrInvalidPage)
			}
		})
	}
}
