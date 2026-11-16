package helpcenter

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/rafialfauzan/layrd/pkg/database"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Repository interface {
	FindCategories(ctx context.Context) ([]Kategori, error)
	FindCategoryByID(ctx context.Context, id uint) (*Kategori, error)
	FindArticles(ctx context.Context, page, limit int, kategoriID, menuID uint, search string) ([]Article, int64, error)
	FindArticleByID(ctx context.Context, id uint) (*Article, error)
	FindManualBooks(ctx context.Context, kategoriID, menuID uint, search string) ([]ManualBook, error)
	FindManualBookByID(ctx context.Context, id uint) (*ManualBook, error)
	SearchAll(ctx context.Context, query string, limit int) ([]SearchResultItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindCategories(ctx context.Context) ([]Kategori, error) {
	var categories []Kategori
	err := r.db.WithContext(ctx).
		Preload("Menus", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC, id ASC")
		}).
		Order("id ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	for i := range categories {
		categories[i].TotalMenus = len(categories[i].Menus)
	}
	return categories, nil
}

func (r *repository) FindCategoryByID(ctx context.Context, id uint) (*Kategori, error) {
	var cat Kategori
	err := r.db.WithContext(ctx).
		Preload("Menus.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC, id ASC")
		}).
		Preload("Menus", "parent_id IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC, id ASC")
		}).
		Preload("Menus.ManualBooks", "status = ?", "published", func(db *gorm.DB) *gorm.DB {
			return db.Order("manual_books.`order` ASC, manual_books.id ASC")
		}).
		First(&cat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Urutkan manual books pada setiap menu berdasarkan kolom order tabel pivot menu_has_manual_books
	if len(cat.Menus) > 0 {
		var menuIDs []uint
		for _, m := range cat.Menus {
			menuIDs = append(menuIDs, m.ID)
		}

		type PivotOrder struct {
			MenuSystemID uint `gorm:"column:menu_system_id"`
			ManualBookID uint `gorm:"column:manual_book_id"`
			Order        int  `gorm:"column:order"`
		}
		var pivots []PivotOrder
		if err := r.db.WithContext(ctx).Table("menu_has_manual_books").
			Where("menu_system_id IN (?)", menuIDs).
			Find(&pivots).Error; err == nil {
			pivotMap := make(map[uint]map[uint]int)
			for _, p := range pivots {
				if _, ok := pivotMap[p.MenuSystemID]; !ok {
					pivotMap[p.MenuSystemID] = make(map[uint]int)
				}
				pivotMap[p.MenuSystemID][p.ManualBookID] = p.Order
			}

			for i := range cat.Menus {
				mID := cat.Menus[i].ID
				if orders, ok := pivotMap[mID]; ok {
					for j := range cat.Menus[i].ManualBooks {
						if ord, exists := orders[cat.Menus[i].ManualBooks[j].ID]; exists {
							cat.Menus[i].ManualBooks[j].Order = ord
						}
					}
					if len(cat.Menus[i].ManualBooks) > 1 {
						sort.SliceStable(cat.Menus[i].ManualBooks, func(a, b int) bool {
							ordA, okA := orders[cat.Menus[i].ManualBooks[a].ID]
							ordB, okB := orders[cat.Menus[i].ManualBooks[b].ID]
							if !okA {
								ordA = 999
							}
							if !okB {
								ordB = 999
							}
							if ordA != ordB {
								return ordA < ordB
							}
							return cat.Menus[i].ManualBooks[a].ID < cat.Menus[i].ManualBooks[b].ID
						})
					}
				}
			}
		}
	}

	return &cat, nil
}

func (r *repository) FindArticles(ctx context.Context, page, limit int, kategoriID, menuID uint, search string) ([]Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	query := r.db.WithContext(ctx).Model(&Article{}).
		Where("articles.status = ?", "published").
		Where("articles.visibility = ?", "publik")

	if menuID > 0 {
		query = query.Joins("JOIN menu_has_articles ON menu_has_articles.article_id = articles.id").
			Where("menu_has_articles.menu_system_id = ?", menuID)
	} else if kategoriID > 0 {
		query = query.Joins("JOIN menu_has_articles ON menu_has_articles.article_id = articles.id").
			Joins("JOIN menu_systems ON menu_systems.id = menu_has_articles.menu_system_id").
			Where("menu_systems.kategori_id = ?", kategoriID)
	}

	if search != "" {
		likeTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("(LOWER(articles.title) LIKE ? OR LOWER(articles.summary) LIKE ?)", likeTerm, likeTerm)
	}

	var total int64
	if err := query.Distinct("articles.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var articles []Article
	err := query.Distinct("articles.*").
		Preload("MenuSystems").
		Order("articles.id ASC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	// Deduplikasi in-memory
	seen := make(map[uint]bool)
	uniqueArticles := make([]Article, 0, len(articles))
	for _, a := range articles {
		if !seen[a.ID] {
			seen[a.ID] = true
			uniqueArticles = append(uniqueArticles, a)
		}
	}

	return uniqueArticles, total, nil
}

func (r *repository) FindArticleByID(ctx context.Context, id uint) (*Article, error) {
	var article Article
	err := r.db.WithContext(ctx).
		Where("status = ?", "published").
		Where("visibility = ?", "publik").
		Preload("MenuSystems").
		First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *repository) FindManualBooks(ctx context.Context, kategoriID, menuID uint, search string) ([]ManualBook, error) {
	query := r.db.WithContext(ctx).Model(&ManualBook{}).
		Where("manual_books.status = ?", "published")

	if menuID > 0 {
		query = query.Joins("JOIN menu_has_manual_books ON menu_has_manual_books.manual_book_id = manual_books.id").
			Where("menu_has_manual_books.menu_system_id = ?", menuID)
	} else if kategoriID > 0 {
		query = query.Joins("JOIN menu_has_manual_books ON menu_has_manual_books.manual_book_id = manual_books.id").
			Joins("JOIN menu_systems ON menu_systems.id = menu_has_manual_books.menu_system_id").
			Where("menu_systems.kategori_id = ?", kategoriID)
	}

	if search != "" {
		likeTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("(LOWER(manual_books.title) LIKE ? OR LOWER(manual_books.description) LIKE ?)", likeTerm, likeTerm)
	}

	var manuals []ManualBook
	if menuID > 0 {
		query = query.Select("manual_books.*, menu_has_manual_books.`order` AS pivot_order").
			Order("menu_has_manual_books.`order` ASC, manual_books.id ASC")
	} else if kategoriID > 0 {
		query = query.Select("manual_books.*, menu_has_manual_books.`order` AS pivot_order").
			Order("menu_systems.`order` ASC, menu_has_manual_books.`order` ASC, manual_books.id ASC")
	} else {
		query = query.Select("manual_books.*").
			Order("manual_books.`order` ASC, manual_books.id ASC")
	}
	err := query.Preload("MenuSystems").Find(&manuals).Error
	if err != nil {
		return nil, err
	}

	// Deduplikasi in-memory
	seen := make(map[uint]bool)
	uniqueManuals := make([]ManualBook, 0, len(manuals))
	for _, m := range manuals {
		if !seen[m.ID] {
			seen[m.ID] = true
			uniqueManuals = append(uniqueManuals, m)
		}
	}

	return uniqueManuals, nil
}

func (r *repository) FindManualBookByID(ctx context.Context, id uint) (*ManualBook, error) {
	var manual ManualBook
	err := r.db.WithContext(ctx).
		Preload("MenuSystems").
		First(&manual, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	return &manual, nil
}

func (r *repository) SearchAll(ctx context.Context, query string, limit int) ([]SearchResultItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	likeTerm := "%" + strings.ToLower(query) + "%"
	var results []SearchResultItem

	// 1. Search Articles
	var articles []Article
	err := r.db.WithContext(ctx).Model(&Article{}).
		Where("status = ?", "published").
		Where("visibility = ?", "publik").
		Where("(LOWER(title) LIKE ? OR LOWER(summary) LIKE ?)", likeTerm, likeTerm).
		Order("id ASC").
		Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	for _, a := range articles {
		snippet := ""
		if a.Summary != nil {
			snippet = *a.Summary
		}
		results = append(results, SearchResultItem{
			ID:      a.ID,
			Type:    "article",
			Title:   a.Title,
			Snippet: snippet,
			Status:  a.Status,
		})
	}

	// 2. Search Manual Books
	remaining := limit - len(results)
	if remaining > 0 {
		var manuals []ManualBook
		err := r.db.WithContext(ctx).Model(&ManualBook{}).
			Where("status = ?", "published").
			Where("(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)", likeTerm, likeTerm).
			Limit(remaining).
			Find(&manuals).Error
		if err != nil {
			return nil, err
		}

		for _, m := range manuals {
			snippet := ""
			if m.Description != nil {
				snippet = *m.Description
			}
			results = append(results, SearchResultItem{
				ID:      m.ID,
				Type:    "manual_book",
				Title:   m.Title,
				Snippet: snippet,
				Status:  m.Status,
			})
		}
	}

	return results, nil
}

