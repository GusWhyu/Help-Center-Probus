package user

import (
	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/response"
	"github.com/rafialfauzan/layrd/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")

	users, total, err := h.service.GetUsers(c.Request.Context(), page, limit, role)
	if err != nil {
		response.InternalError(c, "Failed to fetch users", err.Error())
		return
	}

	response.Paginated(c, "Users retrieved successfully", users, page, limit, total)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.OK(c, "User retrieved successfully", user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User created successfully", user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "User updated successfully", user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Failed to delete user", err.Error())
		return
	}

	response.OK(c, "User deleted successfully", nil)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "Avatar image file is required", err.Error())
		return
	}

	res, err := h.uploader.SaveFile(c, file, "avatars")
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Update user's avatar in database
	avatarURL := res.URL
	user, err := h.service.UpdateUser(c.Request.Context(), uint(id), &UpdateUserRequest{
		Avatar: &avatarURL,
	})
	if err != nil {
		_ = h.uploader.DeleteFile(res.RelativePath)
		response.InternalError(c, "Failed to update user avatar", err.Error())
		return
	}

	response.OK(c, "Avatar uploaded successfully", user)
}
