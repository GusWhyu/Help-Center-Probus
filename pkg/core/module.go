package core

import "github.com/gin-gonic/gin"

// Module represents a pluggable functional module in the framework.
// A module encapsulates its own models, repositories, services, and HTTP routes.
type Module interface {
	// Name returns the unique identifier for this module
	Name() string
	// RegisterRoutes attaches the module's HTTP endpoints to the provided Gin router group
	RegisterRoutes(router *gin.RouterGroup)
}

// Initializer is an optional interface modules can implement to perform pre-flight setup (e.g. cache warm-up, background workers)
type Initializer interface {
	Init(app *App) error
}

// Migrator is an optional interface modules can implement for explicit schema migrations.
type Migrator interface {
	AutoMigrate(app *App) error
}
