package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/rafialfauzan/layrd/pkg/core"
)

func setupStorage(app *core.App) error {
	if err := os.MkdirAll(app.Config.StoragePath, os.ModePerm); err != nil {
		return fmt.Errorf("create storage directory %q: %w", app.Config.StoragePath, err)
	}

	app.Router.Static(app.Config.StorageURL, app.Config.StoragePath)
	log.Printf("[Storage] Serving files from %s at %s\n", app.Config.StoragePath, app.Config.StorageURL)
	return nil
}
