package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrFileTooLarge       = errors.New("file size exceeds maximum allowed limit")
	ErrInvalidFileType    = errors.New("file type not allowed")
	ErrFileNotFound       = errors.New("file not found in form request")
	DefaultMaxUploadSize  = int64(5 * 1024 * 1024) // 5 MB
	DefaultAllowedImages  = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	DefaultAllowedDocs    = []string{".pdf", ".csv", ".xlsx", ".docx"}
)

type UploadResult struct {
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	Extension    string `json:"extension"`
	URL          string `json:"url"`
	RelativePath string `json:"relative_path"`
}

type Uploader struct {
	basePath    string
	baseURL     string
	allowedExts []string
	maxSize     int64
}

// NewUploader creates a new file uploader instance
func NewUploader(basePath, baseURL string) *Uploader {
	return &Uploader{
		basePath:    basePath,
		baseURL:     baseURL,
		allowedExts: DefaultAllowedImages,
		maxSize:     DefaultMaxUploadSize,
	}
}

// WithAllowedExts sets custom allowed file extensions (e.g. [".jpg", ".pdf"])
func (u *Uploader) WithAllowedExts(exts ...string) *Uploader {
	normalized := make([]string, 0, len(exts))
	for _, e := range exts {
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		normalized = append(normalized, strings.ToLower(e))
	}
	u.allowedExts = normalized
	return u
}

// WithMaxSize sets maximum file size in bytes
func (u *Uploader) WithMaxSize(bytes int64) *Uploader {
	u.maxSize = bytes
	return u
}

// SaveFile validates and saves an uploaded multipart file to the designated subdirectory
func (u *Uploader) SaveFile(c *gin.Context, file *multipart.FileHeader, subDir string) (*UploadResult, error) {
	if file == nil {
		return nil, ErrFileNotFound
	}

	// 1. Validate file size
	if file.Size > u.maxSize {
		return nil, fmt.Errorf("%w: max %d MB", ErrFileTooLarge, u.maxSize/(1024*1024))
	}

	// 2. Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !u.isAllowedExt(ext) {
		return nil, fmt.Errorf("%w: %s (allowed: %s)", ErrInvalidFileType, ext, strings.Join(u.allowedExts, ", "))
	}

	// 3. Ensure target directory exists
	targetDir := filepath.Join(u.basePath, subDir)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 4. Generate collision-resistant unique filename: <timestamp>_<randomHex><ext>
	randomBytes := make([]byte, 4)
	_, _ = rand.Read(randomBytes)
	uniqueFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(randomBytes), ext)

	// 5. Save the file to disk
	destination := filepath.Join(targetDir, uniqueFilename)
	if err := c.SaveUploadedFile(file, destination); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Construct public web URL and relative path
	// URL example: /storage/products/1718293849_a1b2c3d4.jpg
	subDirClean := filepath.ToSlash(subDir)
	var urlPath string
	if subDirClean != "" && subDirClean != "." {
		urlPath = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(u.baseURL, "/"), subDirClean, uniqueFilename)
	} else {
		urlPath = fmt.Sprintf("%s/%s", strings.TrimSuffix(u.baseURL, "/"), uniqueFilename)
	}

	relPath := filepath.Join(subDir, uniqueFilename)

	return &UploadResult{
		Filename:     uniqueFilename,
		OriginalName: file.Filename,
		Size:         file.Size,
		Extension:    ext,
		URL:          urlPath,
		RelativePath: relPath,
	}, nil
}

// DeleteFile removes a file from the upload storage
func (u *Uploader) DeleteFile(relativePath string) error {
	fullPath := filepath.Join(u.basePath, relativePath)
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		return nil // File already doesn't exist, ignore
	}
	return os.Remove(fullPath)
}

func (u *Uploader) isAllowedExt(ext string) bool {
	if len(u.allowedExts) == 0 {
		return true // No restrictions if list is empty
	}
	for _, allowed := range u.allowedExts {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}
