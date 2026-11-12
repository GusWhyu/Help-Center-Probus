package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/internal/buildinfo"
	"github.com/rafialfauzan/layrd/pkg/core"
	"github.com/rafialfauzan/layrd/pkg/response"
)

func registerRoutes(app *core.App) {
	// Welcome / Root info endpoint
	apiInfo := func(c *gin.Context) {
		response.OK(c, "Help Center Probus API is running", gin.H{
			"app_name":            "Help Center Probus Backend API",
			"version":             "v1",
			"status":              "online",
			"frontend_portal_url": "http://localhost:5173",
			"endpoints": gin.H{
				"categories":   "/api/v1/categories",
				"articles":     "/api/v1/articles",
				"manual_books": "/api/v1/manual-books",
				"search":       "/api/v1/search?q={query}",
				"health":       "/health",
				"version":      "/version",
			},
		})
	}

	app.Router.GET("/", apiInfo)
	app.APIGroup.GET("", apiInfo)
	app.APIGroup.GET("/", apiInfo)

	app.Router.GET("/version", func(c *gin.Context) {
		response.OK(c, "Application version info", gin.H{
			"version":    buildinfo.Version,
			"build_time": buildinfo.BuildTime,
			"git_commit": buildinfo.GitCommit,
		})
	})

	app.SetupWebSocketRoute("/ws")
}
