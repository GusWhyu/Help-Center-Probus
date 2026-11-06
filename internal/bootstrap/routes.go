package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/internal/buildinfo"
	"github.com/rafialfauzan/layrd/pkg/core"
	"github.com/rafialfauzan/layrd/pkg/response"
)

func registerRoutes(app *core.App) {
	app.Router.GET("/version", func(c *gin.Context) {
		response.OK(c, "Application version info", gin.H{
			"version":    buildinfo.Version,
			"build_time": buildinfo.BuildTime,
			"git_commit": buildinfo.GitCommit,
		})
	})

	app.SetupWebSocketRoute("/ws")
}
