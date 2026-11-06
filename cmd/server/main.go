package main

import (
	"github.com/rafialfauzan/layrd/internal/modules/user"
	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/core"
	"github.com/rafialfauzan/layrd/pkg/response"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// 1. Load application configuration from .env or environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Instantiate Framework Application Engine
	app := core.NewApp(cfg)

	// Version metadata endpoint
	app.Router.GET("/version", func(c *gin.Context) {
		response.OK(c, "Application version info", gin.H{
			"version":    Version,
			"build_time": BuildTime,
			"git_commit": GitCommit,
		})
	})

	// 3. Connect to configured multi-databases (GORM connection pooling)
	if err := app.InitDatabases(); err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	// 4. Enable real-time WebSocket room broadcasting
	app.SetupWebSocketRoute("/ws")

	// 5. Static file storage (uploads)
	_ = os.MkdirAll(cfg.StoragePath, os.ModePerm)
	app.Router.Static(cfg.StorageURL, cfg.StoragePath)
	log.Printf("[Storage] Serving files from %s at %s\n", cfg.StoragePath, cfg.StorageURL)

	// 6. Register pluggable domain modules
	userModule := user.NewModule(app)
	app.RegisterModule(userModule)

	// 7. Run server with graceful shutdown listener
	if err := app.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
