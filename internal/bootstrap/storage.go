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

	// Serve shared Laravel CMS media uploads from available paths
	candidateUploads := []string{
		"../probus-implementator-2026/public/uploads",
		"../probus-implementator/public/uploads",
		"../probus-implementator-main/public/uploads",
		"public/uploads",
	}

	for _, dir := range candidateUploads {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			app.Router.Static("/uploads", dir)
			log.Printf("[Storage] Serving CMS uploaded media from %s at /uploads\n", dir)
			return nil
		}
	}

	_ = os.MkdirAll("public/uploads", os.ModePerm)
	app.Router.Static("/uploads", "public/uploads")

	return nil
}
