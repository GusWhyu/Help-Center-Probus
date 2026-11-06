package bootstrap

import (
	"fmt"

	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/core"
)

func Run() error {
	app, err := initializeApp()
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	return app.Run()
}

func initializeApp() (*core.App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	app := core.NewApp(cfg)

	registerRoutes(app)

	if err := app.InitDatabases(); err != nil {
		return nil, fmt.Errorf("initialize databases: %w", err)
	}

	if err := setupStorage(app); err != nil {
		return nil, err
	}

	registerDependencies(app)

	return app, nil
}
