package user

import (
	"github.com/rafialfauzan/layrd/pkg/core"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	repo    Repository
}

func NewModule(app *core.App) *Module {
	db := app.DBManager.Default()
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, app.Config)

	return &Module{
		handler: handler,
		repo:    repo,
	}
}

func (m *Module) Name() string {
	return "UserModule"
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", m.handler.GetUsers)
		users.GET("/:id", m.handler.GetUserByID)
		users.POST("", m.handler.CreateUser)
		users.PUT("/:id", m.handler.UpdateUser)
		users.DELETE("/:id", m.handler.DeleteUser)
		users.POST("/:id/avatar", m.handler.UploadAvatar)
	}
}

// AutoMigrate runs schema migration on the primary database
func (m *Module) AutoMigrate(app *core.App) error {
	return app.DBManager.Default().AutoMigrate(&User{})
}
