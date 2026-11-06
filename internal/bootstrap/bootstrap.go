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
	initialized := false
	defer func() {
		if !initialized {
			app.DBManager.CloseAll()
		}
	}()

	registerRoutes(app)

	if err := app.InitDatabases(); err != nil {
		return nil, fmt.Errorf("initialize databases: %w", err)
	}

	if err := setupStorage(app); err != nil {
		return nil, err
	}

	if err := registerDependencies(app); err != nil {
		return nil, fmt.Errorf("register dependencies: %w", err)
	}

	initialized = true
	return app, nil
}

func Migrate() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	app := core.NewApp(cfg)
	if err := app.InitDatabases(); err != nil {
		app.DBManager.CloseAll()
		return fmt.Errorf("initialize databases: %w", err)
	}
	defer app.DBManager.CloseAll()

	if err := registerDependencies(app); err != nil {
		return fmt.Errorf("register dependencies: %w", err)
	}
	if err := app.MigrateModules(); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
