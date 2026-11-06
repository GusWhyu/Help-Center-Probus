package user

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/database"
	"github.com/rafialfauzan/layrd/pkg/response"
	"github.com/rafialfauzan/layrd/pkg/utils"
)

const maxUserPageSize = 100

type Handler struct {
	service  Service
	uploader *utils.Uploader
}

func NewHandler(service Service, cfg *config.Config) *Handler {
	uploader := utils.NewUploader(cfg.StoragePath, cfg.StorageURL).
		WithAllowedExts(".jpg", ".jpeg", ".png", ".webp").
		WithMaxSize(2 * 1024 * 1024)

	return &Handler{
		service:  service,
		uploader: uploader,
	}
}

func (h *Handler) GetUsers(c *gin.Context) {
	page, err := parseQueryInt(c, "page", 1)
	if err != nil || page < 1 {
		response.BadRequest(c, "Page must be a positive integer", nil)
		return
	}
	limit, err := parseQueryInt(c, "limit", 10)
	if err != nil || limit < 1 || limit > maxUserPageSize {
		response.BadRequest(c, "Limit must be between 1 and 100", nil)
		return
	}
	if page-1 > int(^uint(0)>>1)/limit {
		response.BadRequest(c, "Page is too large", nil)
		return
	}

	users, total, err := h.service.GetUsers(c.Request.Context(), page, limit, c.Query("role"))
	if err != nil {
		log.Printf("[User] Failed to fetch users: %v", err)
		response.InternalError(c, "Failed to fetch users", nil)
		return
	}

	response.Paginated(c, "Users retrieved successfully", users, page, limit, total)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		log.Printf("[User] Failed to fetch user: %v", err)
		response.InternalError(c, "Failed to fetch user", nil)
		return
	}

	response.OK(c, "User retrieved successfully", user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid user data", nil)
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			response.Conflict(c, "Username or email already registered")
			return
		}
		log.Printf("[User] Failed to create user: %v", err)
		response.InternalError(c, "Failed to create user", nil)
		return
	}

	response.Created(c, "User created successfully", user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid user data", nil)
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		log.Printf("[User] Failed to update user: %v", err)
		response.InternalError(c, "Failed to update user", nil)
		return
	}

	response.OK(c, "User updated successfully", user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		log.Printf("[User] Failed to delete user: %v", err)
		response.InternalError(c, "Failed to delete user", nil)
		return
	}

	response.OK(c, "User deleted successfully", nil)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "Avatar image file is required", nil)
		return
	}

	result, err := h.uploader.SaveFile(c, file, "avatars")
	if err != nil {
		if errors.Is(err, utils.ErrFileTooLarge) || errors.Is(err, utils.ErrInvalidFileType) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		log.Printf("[User] Failed to save avatar: %v", err)
		response.InternalError(c, "Failed to save avatar", nil)
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, &UpdateUserRequest{Avatar: &result.URL})
	if err != nil {
		if cleanupErr := h.uploader.DeleteFile(result.RelativePath); cleanupErr != nil {
			log.Printf("[User] Failed to remove orphaned avatar: %v", cleanupErr)
		}
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		log.Printf("[User] Failed to update user avatar: %v", err)
		response.InternalError(c, "Failed to update user avatar", nil)
		return
	}

	response.OK(c, "Avatar uploaded successfully", user)
}

func parseUserID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		response.BadRequest(c, "Invalid user ID", nil)
		return 0, false
	}
	return uint(id), true
}

func parseQueryInt(c *gin.Context, key string, defaultValue int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
