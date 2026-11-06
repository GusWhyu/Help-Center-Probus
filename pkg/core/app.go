package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/database"
	"github.com/rafialfauzan/layrd/pkg/middleware"
	"github.com/rafialfauzan/layrd/pkg/response"
	"github.com/rafialfauzan/layrd/pkg/utils"
	"github.com/rafialfauzan/layrd/pkg/websocket"

	"github.com/gin-gonic/gin"
)

// App is the central framework application container
type App struct {
	Config       *config.Config
	Router       *gin.Engine
	APIGroup     *gin.RouterGroup
	DBManager    *database.DBManager
	WSHub        *websocket.Hub
	TokenService *utils.TokenService
	modules      []Module
	httpServer   *http.Server
}

// NewApp creates and configures a new framework application
func NewApp(cfg *config.Config) *App {
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	isDebug := cfg.GinMode == gin.DebugMode

	// Global Middlewares
	router.Use(gin.Logger())
	router.Use(middleware.Recovery(isDebug))
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS(cfg.CORSAllowedOrigins, isDebug))
	if cfg.RateLimitRPS > 0 {
		router.Use(middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst))
	}

	// Root health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response.OK(c, "Service is healthy", gin.H{
			"app_name": cfg.AppName,
			"env":      cfg.AppEnv,
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	apiGroup := router.Group(cfg.APIPrefix)

	return &App{
		Config:       cfg,
		Router:       router,
		APIGroup:     apiGroup,
		DBManager:    database.NewDBManager(),
		WSHub:        websocket.NewHub(),
		TokenService: utils.NewTokenService(cfg.JWTSecret, cfg.JWTExpirationHours),
		modules:      make([]Module, 0),
	}
}

// InitDatabases connects to all databases configured in Config
func (a *App) InitDatabases() error {
	isDebug := a.Config.GinMode == gin.DebugMode

	for name, dbCfg := range a.Config.Databases {
		_, err := a.DBManager.Connect(dbCfg, isDebug)
		if err != nil {
			return fmt.Errorf("failed to connect to db [%s]: %w", name, err)
		}
	}
	return nil
}

// RegisterModule registers one or more pluggable functional modules
func (a *App) RegisterModule(modules ...Module) {
	for _, mod := range modules {
		log.Printf("[Framework] Registering module: %s\n", mod.Name())

		// Execute module pre-flight initializer if implemented
		if initMod, ok := mod.(Initializer); ok {
			if err := initMod.Init(a); err != nil {
				log.Fatalf("[Framework] Module [%s] initialization failed: %v", mod.Name(), err)
			}
		}

		// Execute auto-migration if implemented
		if migMod, ok := mod.(Migrator); ok {
			if err := migMod.AutoMigrate(a); err != nil {
				log.Fatalf("[Framework] Module [%s] migration failed: %v", mod.Name(), err)
			}
		}

		// Mount module routes under the base API group
		mod.RegisterRoutes(a.APIGroup)
		a.modules = append(a.modules, mod)
	}
}

// SetupWebSocketRoute enables real-time WebSocket connection endpoint
func (a *App) SetupWebSocketRoute(path string) {
	a.Router.GET(path, func(c *gin.Context) {
		room := c.Query("room")
		if room == "" {
			room = "global"
		}
		websocket.ServeWebSocket(a.WSHub, c, room)
	})
	log.Printf("[Framework] WebSocket endpoint enabled at %s?room=:room\n", path)
}

// Run starts the HTTP server with graceful shutdown handling
func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", a.Config.AppPort)

	a.httpServer = &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("==================================================")
		log.Printf(" 🚀 %s is running on port %s (%s mode)", a.Config.AppName, a.Config.AppPort, a.Config.AppEnv)
		log.Printf(" API Base URL: http://localhost:%s%s", a.Config.AppPort, a.Config.APIPrefix)
		log.Printf(" Health Check: http://localhost:%s/health", a.Config.AppPort)
		log.Printf("==================================================")

		serverErrors <- a.httpServer.ListenAndServe()
	}()

	// Listen for OS interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdown:
		log.Printf("[Framework] Received shutdown signal: %v. Initiating graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Stop HTTP server
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		// Close database pools
		a.DBManager.CloseAll()
		log.Println("[Framework] Graceful shutdown completed cleanly.")
	}

	return nil
}
