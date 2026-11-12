package helpcenter

import (
	"context"
)

type Service interface {
	GetCategories(ctx context.Context) ([]Kategori, error)
	GetCategoryByID(ctx context.Context, id uint) (*Kategori, error)
	GetArticles(ctx context.Context, page, limit int, kategoriID, menuID uint, search string) ([]Article, int64, error)
	GetArticleByID(ctx context.Context, id uint) (*Article, error)
	GetManualBooks(ctx context.Context, menuID uint, search string) ([]ManualBook, error)
	GetManualBookByID(ctx context.Context, id uint) (*ManualBook, error)
	Search(ctx context.Context, query string, limit int) ([]SearchResultItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCategories(ctx context.Context) ([]Kategori, error) {
	return s.repo.FindCategories(ctx)
}

func (s *service) GetCategoryByID(ctx context.Context, id uint) (*Kategori, error) {
	return s.repo.FindCategoryByID(ctx, id)
}

func (s *service) GetArticles(ctx context.Context, page, limit int, kategoriID, menuID uint, search string) ([]Article, int64, error) {
	return s.repo.FindArticles(ctx, page, limit, kategoriID, menuID, search)
}

func (s *service) GetArticleByID(ctx context.Context, id uint) (*Article, error) {
	return s.repo.FindArticleByID(ctx, id)
}

func (s *service) GetManualBooks(ctx context.Context, menuID uint, search string) ([]ManualBook, error) {
	return s.repo.FindManualBooks(ctx, menuID, search)
}

func (s *service) GetManualBookByID(ctx context.Context, id uint) (*ManualBook, error) {
	return s.repo.FindManualBookByID(ctx, id)
}

func (s *service) Search(ctx context.Context, query string, limit int) ([]SearchResultItem, error) {
	return s.repo.SearchAll(ctx, query, limit)
}

