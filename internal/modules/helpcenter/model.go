package helpcenter

import (
	"time"
)

// Kategori represents system/product category (e.g. PHIS, PRESTO, PPOS)
type Kategori struct {
	ID         uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string       `gorm:"size:255;not null" json:"name"`
	CreatedAt  *time.Time   `json:"created_at,omitempty"`
	UpdatedAt  *time.Time   `json:"updated_at,omitempty"`
	Menus      []MenuSystem `gorm:"foreignKey:KategoriID" json:"menus,omitempty"`
	TotalMenus int          `gorm:"-" json:"total_menus"`
}

func (Kategori) TableName() string {
	return "kategori"
}

// MenuSystem represents a menu node in a system hierarchy
type MenuSystem struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	KategoriID  uint         `gorm:"not null" json:"kategori_id"`
	ParentID    *uint        `json:"parent_id"`
	Name        string       `gorm:"size:100;not null" json:"name"`
	Description *string      `gorm:"type:text" json:"description,omitempty"`
	Order       int          `gorm:"default:0" json:"order"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"`
	Children    []MenuSystem `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (MenuSystem) TableName() string {
	return "menu_systems"
}

// Article represents a troubleshooting / knowledge base article
type Article struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint         `gorm:"not null" json:"user_id"`
	Title       string       `gorm:"size:100;not null" json:"title"`
	Visibility  string       `gorm:"size:50;default:'publik'" json:"visibility"` // publik, internal
	Status      string       `gorm:"size:10;default:'draft'" json:"status"`      // published, draft
	Summary     *string      `gorm:"type:text" json:"summary,omitempty"`
	Content     string       `gorm:"type:longtext;not null" json:"content"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"`
	MenuSystems []MenuSystem `gorm:"many2many:menu_has_articles;joinForeignKey:ArticleID;joinReferences:MenuSystemID" json:"menu_systems,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}

// ManualBook represents a structured user manual / feature guide
type ManualBook struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint         `gorm:"not null" json:"user_id"`
	Title       string       `gorm:"size:100;not null" json:"title"`
	Order       int          `gorm:"default:1" json:"order"`
	Description *string      `gorm:"type:text" json:"description,omitempty"`
	Content     *string      `gorm:"type:longtext" json:"content,omitempty"`
	Status      string       `gorm:"size:10;default:'draft'" json:"status"` // published, draft
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"`
	MenuSystems []MenuSystem `gorm:"many2many:menu_has_manual_books;joinForeignKey:ManualBookID;joinReferences:MenuSystemID" json:"menu_systems,omitempty"`
}

func (ManualBook) TableName() string {
	return "manual_books"
}

// SearchResultItem represents unified search item
type SearchResultItem struct {
	ID       uint    `json:"id"`
	Type     string  `json:"type"` // "article" or "manual_book"
	Title    string  `json:"title"`
	Snippet  string  `json:"snippet"`
	Category string  `json:"category,omitempty"`
	Status   string  `json:"status"`
}

