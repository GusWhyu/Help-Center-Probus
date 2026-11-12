package helpcenter

import (
	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/core"
)

type Module struct {
	handler *Handler
	repo    Repository
}

func NewModule(app *core.App) *Module {
	db := app.DBManager.Default()
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
		repo:    repo,
	}
}

func (m *Module) Name() string {
	return "HelpCenterModule"
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	// Public routes for Help Center
	categories := rg.Group("/categories")
	{
		categories.GET("", m.handler.GetCategories)
		categories.GET("/:id", m.handler.GetCategoryByID)
	}

	articles := rg.Group("/articles")
	{
		articles.GET("", m.handler.GetArticles)
		articles.GET("/:id", m.handler.GetArticleByID)
	}

	manuals := rg.Group("/manual-books")
	{
		manuals.GET("", m.handler.GetManualBooks)
		manuals.GET("/:id", m.handler.GetManualBookByID)
	}

	rg.GET("/search", m.handler.Search)
}

