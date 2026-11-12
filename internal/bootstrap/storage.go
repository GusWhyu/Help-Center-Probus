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

	// Serve shared Laravel CMS media uploads from ../probus-implementator-main/public/uploads
	laravelUploads := "../probus-implementator-main/public/uploads"
	if _, err := os.Stat(laravelUploads); err == nil {
		app.Router.Static("/uploads", laravelUploads)
		log.Printf("[Storage] Serving CMS uploaded media from %s at /uploads\n", laravelUploads)
	} else {
		_ = os.MkdirAll("public/uploads", os.ModePerm)
		app.Router.Static("/uploads", "public/uploads")
	}

	return nil
}
