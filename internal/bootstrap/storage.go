package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
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

	// CMS fallback base URL: jika file tidak ada di direktori lokal server Go,
	// arahkan request ke server Laravel CMS tempat file media disimpan.
	cmsBase := os.Getenv("LARAVEL_URL")
	if cmsBase == "" {
		cmsBase = os.Getenv("CMS_STORAGE_URL")
	}
	if cmsBase == "" {
		cmsBase = os.Getenv("LARAVEL_STORAGE_URL")
	}
	if cmsBase == "" {
		if os.Getenv("APP_ENV") == "production" {
			cmsBase = "https://probus.probussystem.com"
		} else {
			cmsBase = "http://127.0.0.1:8000"
		}
	}
	cmsBase = strings.TrimRight(cmsBase, "/")

	uploadsHandler := func(c *gin.Context) {
		reqPath := c.Param("filepath")
		cleanRel := filepath.Clean("/" + reqPath)

		// 1. Periksa direktori lokal (untuk local development)
		for _, dir := range candidateUploads {
			fullPath := filepath.Join(dir, cleanRel)
			if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
				c.File(fullPath)
				return
			}
		}

		// 2. Jika file tidak ada di lokal (misal di live server di mana file diunggah ke Laravel CMS),
		// redirect (HTTP 302 Found) langsung ke domain Laravel CMS
		urlPath := strings.ReplaceAll(cleanRel, "\\", "/")
		targetURL := cmsBase + "/uploads" + urlPath
		c.Redirect(http.StatusFound, targetURL)
	}

	app.Router.GET("/uploads/*filepath", uploadsHandler)
	app.Router.HEAD("/uploads/*filepath", uploadsHandler)

	log.Printf("[Storage] Uploads handler active with CMS fallback: %s\n", cmsBase)
	return nil
}
