package helpcenter

import (
	"context"
	"testing"
)

type mockHelpCenterRepo struct {
	categories []Kategori
	articles   []Article
	manuals    []ManualBook
	searchRes  []SearchResultItem
}

func (m *mockHelpCenterRepo) FindCategories(ctx context.Context) ([]Kategori, error) {
	return m.categories, nil
}

func (m *mockHelpCenterRepo) FindCategoryByID(ctx context.Context, id uint) (*Kategori, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockHelpCenterRepo) FindArticles(ctx context.Context, page, limit int, kategoriID, menuID uint, search string) ([]Article, int64, error) {
	return m.articles, int64(len(m.articles)), nil
}

func (m *mockHelpCenterRepo) FindArticleByID(ctx context.Context, id uint) (*Article, error) {
	for _, a := range m.articles {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockHelpCenterRepo) FindManualBooks(ctx context.Context, menuID uint, search string) ([]ManualBook, error) {
	return m.manuals, nil
}

func (m *mockHelpCenterRepo) FindManualBookByID(ctx context.Context, id uint) (*ManualBook, error) {
	for _, mb := range m.manuals {
		if mb.ID == id {
			return &mb, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockHelpCenterRepo) SearchAll(ctx context.Context, query string, limit int) ([]SearchResultItem, error) {
	return m.searchRes, nil
}

func TestHelpCenterService(t *testing.T) {
	repo := &mockHelpCenterRepo{
		categories: []Kategori{
			{ID: 1, Name: "PRESTO"},
			{ID: 2, Name: "PHIS"},
		},
		articles: []Article{
			{ID: 1, Title: "Printer Thermal Tidak Mencetak", Status: "published", Visibility: "publik"},
		},
		manuals: []ManualBook{
			{ID: 1, Title: "01. Login Sales", Order: 1, Status: "published"},
		},
		searchRes: []SearchResultItem{
			{ID: 1, Type: "article", Title: "Printer Thermal Tidak Mencetak"},
		},
	}

	service := NewService(repo)
	ctx := context.Background()

	// 1. Categories
	cats, err := service.GetCategories(ctx)
	if err != nil || len(cats) != 2 {
		t.Fatalf("expected 2 categories, got %v, err=%v", len(cats), err)
	}

	// 2. Articles
	arts, total, err := service.GetArticles(ctx, 1, 10, 0, 0, "")
	if err != nil || total != 1 || len(arts) != 1 {
		t.Fatalf("expected 1 article, got %v, total=%v, err=%v", len(arts), total, err)
	}

	// 3. Manuals
	mans, err := service.GetManualBooks(ctx, 0, "")
	if err != nil || len(mans) != 1 {
		t.Fatalf("expected 1 manual, got %v, err=%v", len(mans), err)
	}

	// 4. Search
	results, err := service.Search(ctx, "printer", 10)
	if err != nil || len(results) != 1 {
		t.Fatalf("expected 1 search result, got %v, err=%v", len(results), err)
	}
}

